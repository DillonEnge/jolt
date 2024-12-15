package v1

import (
	"fmt"
	"net/http"

	"github.com/DillonEnge/jolt/internal/api"
	"github.com/DillonEnge/jolt/internal/auth"
	"github.com/DillonEnge/jolt/templates"
	"github.com/alexedwards/scs/v2"
)

func HandleNavbar(sm *scs.SessionManager, authClient *auth.Client) api.HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) *api.ApiError {
		claims, _ := authClient.GetClaims(r.Context(), sm)

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

		if claims != nil {
			items = append(
				items,
				templates.NavbarItemData{
					Route: "/my-listings",
					Name:  "mylistings",
					Icon:  "folder",
				},
				templates.NavbarItemData{
					Route: "/create-listing",
					Name:  "createlisting",
					Icon:  "plus-square",
				},
				templates.NavbarItemData{
					Route: "/negotiations",
					Name:  "negotiations",
					Icon:  "dollar-sign",
				})
		}

		templates.Navbar(items, active).Render(r.Context(), w)
		return nil
	}
}
