package v1

import (
	"fmt"
	"net/http"

	"github.com/DillonEnge/jolt/internal/api"
	"github.com/DillonEnge/jolt/templates"
)

func HandleNavbar() api.HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) *api.ApiError {
		active := r.URL.Query().Get("active")

		if active == "" {
			return &api.ApiError{
				Status: http.StatusInternalServerError,
				Err:    fmt.Errorf("failed to find active query param"),
			}
		}

		items := []templates.NavbarItemData{
			{
				Route: "/listings/popular",
				Name:  "trending",
				Icon:  "trending-up",
			},
			{
				Route: "/search",
				Name:  "search",
				Icon:  "search",
			},
		}

		templates.Navbar(items, active).Render(r.Context(), w)
		return nil
	}
}
