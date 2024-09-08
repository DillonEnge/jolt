package v1

import (
	"fmt"
	"net/http"

	"github.com/DillonEnge/jolt/internal/api"
	"github.com/alexedwards/scs/v2"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
)

func HandleSignin(sm *scs.SessionManager, authClient *casdoorsdk.Client) api.HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) *api.ApiError {
		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")

		if code == "" || state == "" {
			return &api.ApiError{
				Status: http.StatusInternalServerError,
				Err:    fmt.Errorf("failed to find code or state query params"),
			}
		}

		token, err := authClient.GetOAuthToken(code, state)
		if err != nil {
			return &api.ApiError{
				Status: http.StatusInternalServerError,
				Err:    err,
			}
		}

		sm.Put(r.Context(), "authToken", token.AccessToken)

		http.Redirect(w, r, "/", http.StatusFound)

		return nil
	}
}
