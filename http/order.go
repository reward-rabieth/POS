package http

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"maragu.dev/gomponents"
	. "maragu.dev/gomponents"
	ghttp "maragu.dev/gomponents/http"

	"github.com/go-chi/chi/v5"
	"github.com/rustacean-dev/possystem/html"
	"github.com/rustacean-dev/possystem/internal/compute"
	"github.com/rustacean-dev/possystem/model"
	"github.com/rustacean-dev/possystem/repository"
	. "maragu.dev/gomponents/html"
)

func OrderRoutes(r chi.Router) {
	r.Get("/orders", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return nil, nil
		}

		orders, err := repository.GetAllOrders(cookie.Value)
		if err != nil {
			return html.LoginPage("Failed to fetch orders"), nil
		}

		return html.OrderHistoryPage(orders), nil
	}))

	// Show the order form
	r.Get("/orders/new", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return nil, nil
		}

		items, err := repository.GetAllItems(cookie.Value)
		if err != nil {
			return html.CreateOrderForm("Failed to fetch items", nil), nil
		}

		return html.CreateOrderForm("", items), nil

	}))

	// Create the order
	r.Post("/orders", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		/* ---------- 1. Auth ---------- */
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return nil, nil
		}

		/* ---------- 2. Parse form ---------- */
		if err := r.ParseForm(); err != nil {
			return html.CreateOrderForm("Invalid form", nil), nil
		}

		itemID := r.FormValue("item_id")
		quantityStr := r.FormValue("quantity")

		qty, err := strconv.Atoi(quantityStr)
		if err != nil || qty <= 0 {
			return html.CreateOrderForm("Quantity must be ≥ 1", nil), nil
		}

		/* ------- 3. Fetch item ------- */
		item, err := repository.GetItemByID(itemID, cookie.Value)
		if err != nil {
			return html.CreateOrderForm("Item not found", nil), nil
		}
		if qty > item.Quantity {
			if item.Quantity == 0 {
				return html.CreateOrderForm(fmt.Sprintf("'%s' is out of stock", item.Name), nil), nil
			}
			return html.CreateOrderForm(fmt.Sprintf("Only %d '%s' left in stock", item.Quantity, item.Name), nil), nil
		}
		/* ---------- 4. Calculate & log total ---------- */
		total := compute.OrderTotal(item.Price, qty)
		fmt.Printf("Create order – item:%s qty:%d total:%s\n",
			item.Name, qty, FormatTZS(total))

		/* ------- 5. Persist order ------- */
		order := model.Order{
			Items:     []model.Item{{ID: item.ID, Name: item.Name, Price: item.Price, Quantity: qty}},
			TotalCost: total,
			Status:    "pending",
		}

		if err := repository.CreateOrder(order, cookie.Value); err != nil {
			return html.CreateOrderForm("Failed to create order", nil), nil
		}

		err = repository.UpdateItemStock(item.ID, item.Quantity-qty, cookie.Value)
		if err != nil {
			fmt.Println(" Failed to update item stock:", err)
			// Optionally: rollback order creation, or just log the error
			log.Println(" Failed to update item stock:", err)
		}

		/* ---------- 6. Redirect to history ---------- */
		http.Redirect(w, r, "/orders", http.StatusSeeOther)
		return nil, nil
	}))

	// Live total calculator for HTMX
	r.Post("/orders/total", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (gomponents.Node, error) {
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return nil, nil
		}

		if err := r.ParseForm(); err != nil {
			return Div(Text("Invalid input")), nil
		}

		itemID := r.FormValue("item_id")
		qtyStr := r.FormValue("quantity")

		qty, err := strconv.Atoi(qtyStr)
		if err != nil || qty <= 0 {
			qty = 1 // default fallback
		}

		item, err := repository.GetItemByID(itemID, cookie.Value)
		if err != nil {
			return Div(Text("Item not found")), nil
		}

		total := item.Price * float64(qty)

		return Div(
			Class("font-semibold text-gray-800"),
			Text(fmt.Sprintf("Total: %.0f TZS", total)),
		), nil
	}))

}
