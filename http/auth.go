package http

import (
	"fmt"
	"net/http"

	"github.com/rustacean-dev/possystem/html"
	"github.com/rustacean-dev/possystem/model"
	"github.com/rustacean-dev/possystem/repository"

	"github.com/go-chi/chi/v5"
	. "maragu.dev/gomponents"
	ghttp "maragu.dev/gomponents/http"
)

// AuthRoutes registers login and logout endpoints for authentication.
func Auth(r chi.Router) {

	// GET /login – Serve the login page
	r.Get("/login", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		// Render the login page
		return html.LoginPage(""), nil

	}))

	// POST /login – Handle login form submission
	r.Post("/login", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		// Parse the form values
		err := r.ParseForm()
		if err != nil {
			return html.LoginPage("Something went wrong. Please try again."), nil
		}

		// Create a login request from form inputs
		login := model.LoginRequest{
			Identity: r.FormValue("identity"),
			Password: r.FormValue("password"),
		}

		// Attempt to authenticate user
		res, err := repository.LoginUser(login)
		fmt.Println("Login response:", res, "Error:", err)
		if err != nil {
			// Show friendly error if credentials are invalid
			return html.LoginPage("Invalid email or password. Please try again."), nil
		}

		// Set token in a secure, HTTP-only cookie on successful login
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    res.Token,
			Path:     "/",
			HttpOnly: true,
		})

		// Tell HTMX to redirect to /orders after login
		w.Header().Set("HX-Redirect", "/orders")
		return nil, nil
	}))

	// GET /logout – Clear the authentication cookie and redirect to home
	r.Get("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			MaxAge:   -1,
		})

		// Redirect to home page after logout
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

}
