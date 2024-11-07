package v1

import (
	"net/http"

	"github.com/DillonEnge/jolt/internal/api"
	"github.com/alexedwards/scs/v2"
)

func HandleLogout(sm *scs.SessionManager) api.HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) *api.ApiError {
		sm.Remove(r.Context(), "authToken")

		http.Redirect(w, r, "/", http.StatusFound)

		return nil
	}
}
