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
	"github.com/DillonEnge/jolt/internal/messagequeue"
	"github.com/DillonEnge/jolt/templates"
	"github.com/alexedwards/scs/v2"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MessageRecorder interface {
	RecordMessage(context.Context, database.RecordMessageParams) (database.Message, error)
}

type PostMessageParams struct {
	NegotiationID string `json:"negotiation_id"`
	Message       string `json:"message"`
}

func HandleMessageWS(db MessageRecorder, authClient *auth.Client, mq *messagequeue.Store, sm *scs.SessionManager) api.HandlerFuncWithError {
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

		mq.AddTopic(negotiationID)

		topicChan, err := mq.Subscribe(ctx, negotiationID)
		if err != nil {
			return &api.ApiError{
				Status: http.StatusInternalServerError,
				Err:    err,
			}
		}
		defer mq.Unsubscribe(negotiationID, topicChan)

		slog.Info("Subscribed", "topic", initialParams.NegotiationID, "user", claims.Email)

		go func() {
			for {
				message, ok := <-topicChan
				if !ok {
					slog.Info("topic channel closed")
					break
				}

				if message.MessageText == "" {
					continue
				}

				var buf bytes.Buffer
				templates.MessageOOB(*message, claims).Render(ctx, &buf)

				d, err := io.ReadAll(&buf)
				if err != nil {
					slog.Error("failed to read from message data buffer")
					continue
				}

				c.Write(ctx, websocket.MessageText, d)
			}
		}()

		for {
			var params PostMessageParams
			err = wsjson.Read(ctx, c, &params)
			if err != nil {
				break
			}

			if params.Message == "" {
				slog.Warn("encountered blank message, skipping...", "params", params)
				continue
			}

			p := database.RecordMessageParams{
				NegotiationID: negotiationID,
				SenderEmail:   claims.Email,
				SenderName:    claims.DisplayName,
				MessageText:   params.Message,
			}

			newMessage, err := db.RecordMessage(ctx, p)
			if err != nil {
				break
			}

			mq.Publish(negotiationID, &database.Message{
				SenderEmail: newMessage.SenderEmail,
				SenderName:  newMessage.SenderName,
				MessageText: newMessage.MessageText,
				TimeSent:    newMessage.TimeSent,
			})
		}

		// c.Close(websocket.StatusNormalClosure, "")

		slog.Error("encountered err", "error", err)
		return &api.ApiError{
			Status: http.StatusInternalServerError,
			Err:    err,
		}

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
