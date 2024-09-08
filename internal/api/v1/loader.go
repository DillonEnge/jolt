package v1

import (
	"errors"
	"net/http"

	"github.com/DillonEnge/jolt/internal/api"
	"github.com/DillonEnge/jolt/templates"
	"github.com/a-h/templ"
)

func HandleLoader() api.HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) *api.ApiError {
		route := r.URL.Query().Get("route")
		if route == "" {
			return &api.ApiError{
				Status: http.StatusInternalServerError,
				Err:    errors.New("failed to provide 'route' query param"),
			}
		}

		templ.Handler(templates.Loader(route)).ServeHTTP(w, r)

		return nil
	}
}
