package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rustacean-dev/possystem/html"
	. "maragu.dev/gomponents"

	ghttp "maragu.dev/gomponents/http"
)

// Home sets up the root ("/") route.
//
// This route renders the homepage and checks for a valid JWT token
// in the `token` cookie to determine if the user is authenticated.
//
// If authenticated, the page might render a dashboard or show user-specific data.
// Otherwise, it renders a guest-friendly homepage.

func Home(r chi.Router) {
	r.Get("/", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		// Check if user is authenticated
		cookie, err := r.Cookie("token")
		authenticated := err == nil && cookie.Value != ""

		return html.HomePage(authenticated), nil
	}))
}
