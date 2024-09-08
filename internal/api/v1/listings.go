package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DillonEnge/jolt/database"
	"github.com/DillonEnge/jolt/internal/api"
	"github.com/DillonEnge/jolt/templates"
	"github.com/alexedwards/scs/v2"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RecordListingParams struct {
	SellerEmail string  `json:"seller_email"`
	ListingName string  `json:"listing_name"`
	Description string  `json:"description"`
	Price       float32 `json:"price,string"`
}

func HandleListings(db *pgxpool.Pool) api.HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) *api.ApiError {
		title := r.URL.Query().Get("title")

		if title == "" {
			return &api.ApiError{
				Status: http.StatusInternalServerError,
				Err:    fmt.Errorf("failed to provide title query param"),
			}
		}

		name := r.URL.Query().Get("name")

		if name == "" {
			templates.NoResults().Render(r.Context(), w)
			return nil
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

		rows, err := queries.ListingsByLikeName(r.Context(), name)
		if err != nil {
			return &api.ApiError{
				Status: http.StatusInternalServerError,
				Err:    err,
			}
		}

		if len(rows) == 0 {
			templates.NoResults().Render(r.Context(), w)
			return nil
		}

		w.WriteHeader(http.StatusOK)
		templates.Listings(title, rows).Render(r.Context(), w)

		return nil
	}
}

func HandleMyListings(db *pgxpool.Pool, authClient *casdoorsdk.Client, sm *scs.SessionManager) api.HandlerFuncWithError {
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

		rows, err := queries.ListingsBySellerEmail(r.Context(), claims.Email)
		if err != nil {
			return &api.ApiError{
				Status: http.StatusInternalServerError,
				Err:    err,
			}
		}

		if len(rows) == 0 {
			templates.NoResults().Render(r.Context(), w)
			return nil
		}

		w.WriteHeader(http.StatusOK)
		templates.Listings("My Listings", rows).Render(r.Context(), w)

		return nil
	}
}

func HandlePostListings(db *pgxpool.Pool) api.HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) *api.ApiError {
		var params RecordListingParams
		err := json.NewDecoder(r.Body).Decode(&params)
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

		row, err := queries.RecordListing(r.Context(), database.RecordListingParams{
			SellerEmail: params.SellerEmail,
			ListingName: params.ListingName,
			Description: params.Description,
			Price:       int32(params.Price * 100),
		})
		if err != nil {
			return &api.ApiError{
				Status: http.StatusInternalServerError,
				Err:    err,
			}
		}

		tx.Commit(r.Context())

		w.WriteHeader(http.StatusOK)
		templates.IndividualListing(row).Render(r.Context(), w)

		return nil
	}
}

func HandleDeleteListings(db *pgxpool.Pool) api.HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) *api.ApiError {
		id := r.URL.Query().Get("id")

		if _, err := uuid.FromString(id); err != nil {
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

		row, err := queries.DeleteListing(r.Context(), id)
		if err != nil {
			return &api.ApiError{
				Status: http.StatusInternalServerError,
				Err:    err,
			}
		}

		tx.Commit(r.Context())

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(row)

		return nil
	}
}
