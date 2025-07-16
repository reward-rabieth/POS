package http

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/rustacean-dev/possystem/model"

	"github.com/go-chi/chi/v5"
	"github.com/rustacean-dev/possystem/html"
	"github.com/rustacean-dev/possystem/repository"
	. "maragu.dev/gomponents"
	ghttp "maragu.dev/gomponents/http"
)

func ItemRoutes(r chi.Router) {
	r.Get("/items/new", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		//  Redirect if not logged in
		cookie, err := r.Cookie("token")
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return nil, nil
		}

		return html.NewItemPage(""), nil
	}))

	r.Post("/items/new", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		//  Redirect if not logged in
		cookie, err := r.Cookie("token")
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return nil, nil
		}

		// Parse form
		if err := r.ParseForm(); err != nil {
			return html.NewItemPage("Invalid form submission"), nil
		}

		// Parse price
		price, err := strconv.ParseFloat(r.FormValue("price"), 64)
		if err != nil || price <= 0 {
			return html.NewItemPage("Please enter a valid price"), nil
		}

		// Parse quantity (optional: default to 0)
		quantity := 0
		if q := r.FormValue("quantity"); q != "" {
			quantity, _ = strconv.Atoi(q)
		}

		// Normalize item name (case-insensitive matching)
		name := strings.ToLower(strings.TrimSpace(r.FormValue("name")))
		description := r.FormValue("description")

		// Check if item already exists
		existingItem, err := repository.GetItemByName(name, cookie.Value)
		if err == nil {
			//  If item exists, increase quantity only
			newQty := existingItem.Quantity + quantity
			err := repository.UpdateItemStock(existingItem.ID, newQty, cookie.Value)
			if err != nil {
				return html.NewItemPage("Failed to update stock of existing item"), nil
			}

			//  Redirect (stock updated)
			w.Header().Set("HX-Redirect", "/orders/new")
			return nil, nil
		}

		//  If not found, create new item
		item := model.Item{
			Name:        name,
			Price:       price,
			Description: description,
			Quantity:    quantity,
		}

		err = repository.CreateItem(item, cookie.Value)
		if err != nil {
			return html.NewItemPage(fmt.Sprintf("Failed to create item: %s", err.Error())), nil
		}

		w.Header().Set("HX-Redirect", "/orders/new")
		return nil, nil
	}))

}
