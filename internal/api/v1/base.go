package v1

import (
	"context"
	"net/http"

	"github.com/DillonEnge/jolt/internal/api"
	"github.com/DillonEnge/jolt/internal/auth"
	"github.com/DillonEnge/jolt/templates"
	"github.com/alexedwards/scs/v2"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
)

func HandleBase(sm *scs.SessionManager, authClient *auth.Client, config *api.Config) api.HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) *api.ApiError {
		var claims *casdoorsdk.Claims
		var err error

		token := sm.GetString(r.Context(), "authToken")
		if token != "" {
			claims, err = authClient.ParseJwtToken(token)
			if err != nil {
				return &api.ApiError{
					Status: http.StatusInternalServerError,
					Err:    err,
				}
			}
		}

		templates.Base(claims, config).Render(context.Background(), w)

		return nil
	}
}
