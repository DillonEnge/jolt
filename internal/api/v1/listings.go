package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/DillonEnge/jolt/database"
	"github.com/DillonEnge/jolt/internal/api"
	"github.com/DillonEnge/jolt/internal/auth"
	"github.com/DillonEnge/jolt/templates"
	"github.com/DillonEnge/seaweedfs-go-client"
	"github.com/alexedwards/scs/v2"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ListingFetcher interface {
	ListingByID(ctx context.Context, listingID string) (database.ListingWithImageUrl, error)
	ListingsByLikeName(ctx context.Context, listingName string) ([]database.ListingWithImageUrl, error)
}

type ListingViewsUpserter interface {
	UpsertListingViews(ctx context.Context, listingID string) (database.ListingView, error)
}

type ListingRecorderFetcher interface {
	RecordListing(ctx context.Context, arg database.RecordListingParams) (database.Listing, error)
	RecordListingImages(ctx context.Context, arg database.RecordListingImagesParams) ([]database.ListingImage, error)
	ListingByID(ctx context.Context, listingID string) (database.ListingWithImageUrl, error)
}

type ListingsByViewsFetcher interface {
	ListingsByViews(ctx context.Context, arg database.ListingsByViewsParams) ([]database.ListingWithImageUrl, error)
}

type RecordListingParams struct {
	SellerEmail string  `json:"seller_email"`
	ListingName string  `json:"listing_name"`
	Description string  `json:"description"`
	Price       float32 `json:"price,string"`
}

func HandleListings(db ListingFetcher, authClient *auth.Client, sm *scs.SessionManager) api.HandlerFuncWithError {
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
		claims, err := authClient.ParseJwtToken(token)
		if err != nil {
			slog.Error("failed to decode token", "err", err)
			claims = nil
		}

		w.WriteHeader(http.StatusOK)
		templates.Listings(title, listings, claims, token != "").Render(r.Context(), w)

		return nil
	}
}

func HandlePopularListings(db ListingsByViewsFetcher, authClient *auth.Client, sm *scs.SessionManager) api.HandlerFuncWithError {
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
		claims, err := authClient.ParseJwtToken(token)
		if err != nil {
			slog.Error("failed to decode token", "err", err)
			claims = nil
		}

		w.WriteHeader(http.StatusOK)
		templates.Listings("Popular Listings", rows, claims, token != "").Render(r.Context(), w)

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
		templates.Listings("My Listings", listings, nil, false).Render(r.Context(), w)

		return nil
	}
}

func HandlePostListings(db ListingRecorderFetcher, fsClient *seaweedfs.Client, config *api.Config) api.HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) *api.ApiError {
		// Parse multipart form with 10MB max memory
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			return &api.ApiError{
				Status: http.StatusBadRequest,
				Err:    err,
			}
		}

		// Get form values
		sellerEmail := r.FormValue("seller_email")
		listingName := r.FormValue("listing_name")
		description := r.FormValue("description")
		priceStr := r.FormValue("price")

		// Convert price to float
		price, err := strconv.ParseFloat(priceStr, 32)
		if err != nil {
			return &api.ApiError{
				Status: http.StatusBadRequest,
				Err:    fmt.Errorf("invalid price format: %v", err),
			}
		}

		listingID, err := uuid.NewV4()
		if err != nil {
			return &api.ApiError{
				Status: http.StatusInternalServerError,
				Err:    fmt.Errorf("unable to create listing id: %v", err),
			}
		}
		imageURLs := []string{}
		// Get image files
		files := r.MultipartForm.File["images"]
		if len(files) > 0 {
			slog.Info("Processing images", "count", len(files))
			// TODO: Save images to storage service, attach to listing
			for i, fileHeader := range files {
				slog.Info("Image file", "index", i, "filename", fileHeader.Filename, "size", fileHeader.Size)
				f, err := fileHeader.Open()
				if err != nil {
					return &api.ApiError{
						Status: http.StatusInternalServerError,
						Err:    fmt.Errorf("unable to open image fileHeader: %v", err),
					}
				}
				defer f.Close()

				daResp, err := fsClient.DirAssign()
				if err != nil {
					return &api.ApiError{
						Status: http.StatusInternalServerError,
						Err:    fmt.Errorf("failed to assign a dir via fsClient: %v", err),
					}
				}

				ufResp, err := fsClient.UploadFile(f, fileHeader.Filename, daResp.FID)
				if err != nil {
					return &api.ApiError{
						Status: http.StatusInternalServerError,
						Err:    fmt.Errorf("failed to upload file via fsClient: %v", err),
					}
				}
				if ufResp.Size == 0 {
					return &api.ApiError{
						Status: http.StatusInternalServerError,
						Err:    fmt.Errorf("failed to upload file via fsClient: %s", "image size is 0"),
					}
				}

				imageURLs = append(imageURLs, fmt.Sprintf("%s/%s", config.SeaweedFS.VolumesURL, daResp.FID))
			}
		}

		// Record the listing in the database
		_, err = db.RecordListing(r.Context(), database.RecordListingParams{
			ID:          listingID.String(),
			SellerEmail: sellerEmail,
			ListingName: listingName,
			Description: description,
			Price:       int32(float32(price) * 100),
		})
		if err != nil {
			return &api.ApiError{
				Status: http.StatusInternalServerError,
				Err:    err,
			}
		}

		_, err = db.RecordListingImages(r.Context(), database.RecordListingImagesParams{
			ListingID:     listingID.String(),
			ImageUrlArray: imageURLs,
		})
		if err != nil {
			return &api.ApiError{
				Status: http.StatusInternalServerError,
				Err:    err,
			}
		}

		listing, err := db.ListingByID(r.Context(), listingID.String())

		w.WriteHeader(http.StatusOK)
		templates.IndividualListing(listing, nil, false).Render(r.Context(), w)

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
