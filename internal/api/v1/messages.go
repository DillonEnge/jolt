package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/DillonEnge/jolt/database"
	"github.com/DillonEnge/jolt/internal/api"
	"github.com/DillonEnge/jolt/internal/auth"
	"github.com/DillonEnge/jolt/templates"
	"github.com/alexedwards/scs/v2"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
)

type MessageRecorder interface {
	RecordMessage(context.Context, database.RecordMessageParams) (database.Message, error)
}

type PostMessageParams struct {
	NegotiationID string `json:"negotiation_id"`
	Message       string `json:"message"`
}

func HandleMessageWS(db MessageRecorder, authClient *auth.Client, nc *nats.Conn, sm *scs.SessionManager) api.HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) *api.ApiError {
		claims, err := authClient.GetClaims(r.Context(), sm)
		if err != nil {
			return &api.ApiError{
				Status: http.StatusUnauthorized,
				Err:    err,
			}
		}

		c, err := websocket.Accept(w, r, nil)
		if err != nil {
			return &api.ApiError{
				Status: http.StatusInternalServerError,
				Err:    err,
			}
		}
		defer c.CloseNow()

		// Set the context as needed. Use of r.Context() is not recommended
		// to avoid surprising behavior (see http.Hijacker).
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute*30)
		defer cancel()

		// Get initial negotiation ID
		var initialParams PostMessageParams
		err = wsjson.Read(ctx, c, &initialParams)
		if err != nil {
			return &api.ApiError{
				Status: http.StatusInternalServerError,
				Err:    err,
			}
		}

		negotiationID := initialParams.NegotiationID

		nc.Subscribe(negotiationID, func(msg *nats.Msg) {
			var m database.Message
			var msgErr error
			msgErr = json.Unmarshal(msg.Data, &m)
			if msgErr != nil {
				slog.Error("failed to unmarshal msg data", "msg", msg.Data)
				return
			}

			var buf bytes.Buffer
			templates.MessageOOB(m, claims).Render(ctx, &buf)

			d, msgErr := io.ReadAll(&buf)
			if msgErr != nil {
				slog.Error("failed to read from message oob buf", "err", err)
				return
			}

			c.Write(ctx, websocket.MessageText, d)
		})

		slog.Info("Subscribed", "topic", initialParams.NegotiationID, "user", claims.Email)

		for {
			var params PostMessageParams
			err = wsjson.Read(ctx, c, &params)
			if err != nil {
				slog.Error("error reading json from ws", "err", err)
				break
			}

			if params.Message == "" {
				slog.Warn("encountered blank message, skipping...", "params", params)
				continue
			}

			slog.Info("recieved message from ws", "msg", params)

			p := database.RecordMessageParams{
				NegotiationID: negotiationID,
				SenderEmail:   claims.Email,
				SenderName:    claims.DisplayName,
				MessageText:   params.Message,
			}

			newMessage, err := db.RecordMessage(ctx, p)
			if err != nil {
				slog.Error("failed to persist message in db", "err", err)
				continue
			}

			payload, err := json.Marshal(newMessage)
			if err != nil {
				slog.Error("failed to marshal message to json", "err", err)
				continue
			}

			nc.Publish(negotiationID, payload)
		}

		slog.Error("encountered err", "error", err)

		return nil
	}
}

func HandleMessages(db *pgxpool.Pool) api.HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) *api.ApiError {

		return nil
	}
}

func HandlePostMessage(db *pgxpool.Pool, sm *scs.SessionManager, authClient *auth.Client) api.HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) *api.ApiError {
		var params PostMessageParams

		err := json.NewDecoder(r.Body).Decode(&params)
		if err != nil {
			return &api.ApiError{
				Status: http.StatusInternalServerError,
				Err:    err,
			}
		}

		claims, err := authClient.GetClaims(r.Context(), sm)
		if err != nil {
			return &api.ApiError{
				Status: http.StatusUnauthorized,
				Err:    err,
			}
		}

		queries, tx, err := database.NewQueries(r.Context(), db)
		if err != nil {
			return &api.ApiError{
				Status: http.StatusInternalServerError,
				Err:    err,
			}
		}
		defer tx.Rollback(r.Context())

		_, err = queries.RecordMessage(r.Context(), database.RecordMessageParams{
			NegotiationID: params.NegotiationID,
			SenderEmail:   claims.Email,
			SenderName:    claims.Name,
			MessageText:   params.Message,
		})
		if err != nil {
			return &api.ApiError{
				Status: http.StatusInternalServerError,
				Err:    err,
			}
		}

		tx.Commit(r.Context())

		return nil
	}
}
