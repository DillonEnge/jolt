package v1

import (
	"context"
	"fmt"
	"net/http"

	"github.com/DillonEnge/jolt/database"
	"github.com/DillonEnge/jolt/internal/api"
	"github.com/DillonEnge/jolt/internal/auth"
	"github.com/DillonEnge/jolt/templates"
	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NegotiationQuerier interface {
	RecordNegotiation(ctx context.Context, arg database.RecordNegotiationParams) (database.Negotiation, error)
	NegotiationByListingIDAndBuyerEmail(ctx context.Context, arg database.NegotiationByListingIDAndBuyerEmailParams) (database.Negotiation, error)
}

func HandleNegotiations(db *pgxpool.Pool, authClient *auth.Client, sm *scs.SessionManager) api.HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) *api.ApiError {
		token := sm.GetString(r.Context(), "authToken")

		if token == "" {
			return &api.ApiError{
				Status: http.StatusNotFound,
				Err:    fmt.Errorf("failed to find session auth token"),
			}
		}

		claims, err := authClient.ParseJwtToken(token)
		if err != nil {
			return &api.ApiError{
				Status: http.StatusInternalServerError,
				Err:    err,
			}
		}

		tx, err := db.Begin(r.Context())
		if err != nil {
			return &api.ApiError{
				Status: http.StatusInternalServerError,
				Err:    err,
			}
		}
		defer tx.Rollback(r.Context())

		queries := database.New(tx)

		negotiations, err := queries.NegotiationsByEmail(r.Context(), claims.Email)
		templates.Negotiations(negotiations, claims).Render(r.Context(), w)

		return nil
	}
}

func HandlePostNegotiation(db NegotiationQuerier, authClient *auth.Client, sm *scs.SessionManager) api.HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) *api.ApiError {
		listingID := r.URL.Query().Get("listing_id")
		if listingID == "" {
			return &api.ApiError{
				Status: http.StatusBadRequest,
				Err:    fmt.Errorf("failed to provide listing_id query param"),
			}
		}

		claims, err := authClient.GetClaims(r.Context(), sm)
		if err != nil {
			return &api.ApiError{
				Status: http.StatusInternalServerError,
				Err:    err,
			}
		}

		negotiation, err := db.RecordNegotiation(r.Context(), database.RecordNegotiationParams{
			ListingID:  listingID,
			BuyerEmail: claims.Email,
		})

		if negotiation.ID == "" {
			negotiation, err = db.NegotiationByListingIDAndBuyerEmail(r.Context(), database.NegotiationByListingIDAndBuyerEmailParams{
				ListingID:  listingID,
				BuyerEmail: claims.Email,
			})
		}

		templates.Loader(fmt.Sprintf("/chat?negotiation_id=%s", negotiation.ID)).Render(r.Context(), w)

		return nil
	}
}
