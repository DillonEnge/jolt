package v1

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/DillonEnge/jolt/database"
	"github.com/DillonEnge/jolt/internal/api"
	"github.com/DillonEnge/jolt/internal/auth"
	"github.com/DillonEnge/jolt/templates"
	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func HandleChat(db *pgxpool.Pool, sm *scs.SessionManager, authClient *auth.Client, config *api.Config) api.HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) *api.ApiError {
		negotiationID := r.URL.Query().Get("negotiation_id")
		if negotiationID == "" {
			return &api.ApiError{
				Status: http.StatusBadRequest,
				Err:    fmt.Errorf("failed to provide negotiation_id query param"),
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

		messages, err := queries.MessagesByNegotiationID(r.Context(), negotiationID)
		if err != nil {
			return &api.ApiError{
				Status: http.StatusInternalServerError,
				Err:    err,
			}
		}

		slog.Info("messages", "messages", messages)

		templates.Chat(messages, negotiationID, claims).Render(r.Context(), w)

		return nil
	}
}
