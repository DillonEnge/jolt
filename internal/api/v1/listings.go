package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/DillonEnge/jolt/database"
	"github.com/DillonEnge/jolt/internal/api"
	"github.com/DillonEnge/jolt/internal/auth"
	"github.com/DillonEnge/jolt/templates"
	"github.com/alexedwards/scs/v2"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ListingFetcher interface {
	ListingByID(ctx context.Context, listingID string) ([]database.Listing, error)
	ListingsByLikeName(ctx context.Context, listingName string) ([]database.Listing, error)
}

type ListingViewsUpserter interface {
	UpsertListingViews(ctx context.Context, listingID string) (database.ListingView, error)
}

type ListingRecorder interface {
	RecordListing(ctx context.Context, arg database.RecordListingParams) (database.Listing, error)
}

type ListingsByViewsFetcher interface {
	ListingsByViews(ctx context.Context, arg database.ListingsByViewsParams) ([]database.Listing, error)
}

type RecordListingParams struct {
	SellerEmail string  `json:"seller_email"`
	ListingName string  `json:"listing_name"`
	Description string  `json:"description"`
	Price       float32 `json:"price,string"`
}

func HandleListings(db ListingFetcher, sm *scs.SessionManager) api.HandlerFuncWithError {
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

		listings, err := db.ListingsByLikeName(r.Context(), name)
		if err != nil {
			return &api.ApiError{
				Status: http.StatusInternalServerError,
				Err:    err,
			}
		}

		if len(listings) == 0 {
			templates.NoResults().Render(r.Context(), w)
			return nil
		}

		token := sm.GetString(r.Context(), "authToken")

		w.WriteHeader(http.StatusOK)
		templates.Listings(title, listings, token != "").Render(r.Context(), w)

		return nil
	}
}

func HandlePopularListings(db ListingsByViewsFetcher, sm *scs.SessionManager) api.HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) *api.ApiError {
		pageSizeParam := r.URL.Query().Get("page_size")

		if pageSizeParam == "" {
			pageSizeParam = "10"
		}

		pageSize, err := strconv.Atoi(pageSizeParam)
		if err != nil {
			return &api.ApiError{
				Status: http.StatusInternalServerError,
				Err:    err,
			}
		}

		pageNumberParam := r.URL.Query().Get("page_number")

		if pageNumberParam == "" {
			pageNumberParam = "1"
		}

		pageNumber, err := strconv.Atoi(pageNumberParam)
		if err != nil {
			return &api.ApiError{
				Status: http.StatusInternalServerError,
				Err:    err,
			}
		}

		rows, err := db.ListingsByViews(r.Context(), database.ListingsByViewsParams{
			Limit:  int32(pageSize),
			Offset: int32((pageNumber - 1) * pageSize),
		})

		if err != nil {
			return &api.ApiError{
				Status: http.StatusInternalServerError,
				Err:    err,
			}
		}

		token := sm.GetString(r.Context(), "authToken")

		w.WriteHeader(http.StatusOK)
		templates.Listings("Popular Listings", rows, token != "").Render(r.Context(), w)

		return nil
	}
}

func HandleMyListings(db *pgxpool.Pool, authClient *auth.Client, sm *scs.SessionManager) api.HandlerFuncWithError {
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

		queries, tx, err := database.NewQueries(r.Context(), db)
		if err != nil {
			return &api.ApiError{
				Status: http.StatusInternalServerError,
				Err:    err,
			}
		}
		defer tx.Rollback(r.Context())

		listings, err := queries.ListingsBySellerEmail(r.Context(), claims.Email)
		if err != nil {
			return &api.ApiError{
				Status: http.StatusInternalServerError,
				Err:    err,
			}
		}

		if len(listings) == 0 {
			templates.NoResults().Render(r.Context(), w)
			return nil
		}

		w.WriteHeader(http.StatusOK)
		templates.Listings("My Listings", listings, true).Render(r.Context(), w)

		return nil
	}
}

func HandlePostListings(db ListingRecorder) api.HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) *api.ApiError {
		var params RecordListingParams
		err := json.NewDecoder(r.Body).Decode(&params)
		if err != nil {
			return &api.ApiError{
				Status: http.StatusInternalServerError,
				Err:    err,
			}
		}
		row, err := db.RecordListing(r.Context(), database.RecordListingParams{
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

		w.WriteHeader(http.StatusOK)
		templates.IndividualListing(row, false).Render(r.Context(), w)

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

		queries, tx, err := database.NewQueries(r.Context(), db)
		if err != nil {
			return &api.ApiError{
				Status: http.StatusInternalServerError,
				Err:    err,
			}
		}
		defer tx.Rollback(r.Context())

		listing, err := queries.DeleteListing(r.Context(), id)
		if err != nil {
			return &api.ApiError{
				Status: http.StatusInternalServerError,
				Err:    err,
			}
		}

		tx.Commit(r.Context())

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(listing)

		return nil
	}
}

func HandlePatchListing(db ListingViewsUpserter) api.HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) *api.ApiError {
		id := r.URL.Query().Get("id")

		if _, err := uuid.FromString(id); err != nil {
			return &api.ApiError{
				Status: http.StatusInternalServerError,
				Err:    err,
			}
		}

		_, err := db.UpsertListingViews(r.Context(), id)
		if err != nil {
			return &api.ApiError{
				Status: http.StatusInternalServerError,
				Err:    err,
			}
		}

		return nil
	}
}
