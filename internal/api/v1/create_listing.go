package v1

import (
	"net/http"

	"github.com/DillonEnge/jolt/internal/api"
	"github.com/DillonEnge/jolt/internal/auth"
	"github.com/DillonEnge/jolt/templates"
	"github.com/alexedwards/scs/v2"
)

func HandleCreateListing(sm *scs.SessionManager, authClient *auth.Client) api.HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) *api.ApiError {
		token := sm.GetString(r.Context(), "authToken")

		claims, err := authClient.ParseJwtToken(token)
		if err != nil {
			return &api.ApiError{
				Status: http.StatusInternalServerError,
				Err:    err,
			}
		}
		templates.CreateListing(claims).Render(r.Context(), w)

		return nil
	}
}
