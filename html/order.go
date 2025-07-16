package html

import (
	"fmt"
	"math"

	"github.com/dustin/go-humanize"
	"github.com/rustacean-dev/possystem/model"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// OrderHistoryPage renders the /orders page.
// It displays a table of all orders, showing key information:
// - Order ID
// - Customer name (if available, otherwise "Unknown")
// - Order creation date and time
// - Total cost (formatted in TSh)
// - Order status (e.g., pending, completed)

func OrderHistoryPage(orders []model.Order) Node {
	return Layout("/orders", true,
		Div(
			ID("main"),
			Class("max-w-6xl mx-auto mt-12"),

			H2(Class("text-2xl font-bold mb-4 text-gray-800"), Text("Order History")),

			Table(Class("min-w-full bg-white border border-gray-200 rounded-md overflow-hidden"),
				THead(Class("bg-indigo-700 text-white"),
					Tr(
						Th(Class("px-4 py-2 text-left"), Text("Order ID")),
						Th(Class("px-4 py-2 text-left"), Text("Customer")),
						Th(Class("px-4 py-2 text-left"), Text("Date")),
						Th(Class("px-4 py-2 text-left"), Text("Time")),
						Th(Class("px-4 py-2 text-left"), Text("Total")),
						Th(Class("px-4 py-2 text-left"), Text("Status")),
					),
				),

				TBody(
					Group(func() []Node {
						var rows []Node
						for _, o := range orders {
							// Fallback for missing customer name
							customer := o.Expand.User.Username
							if customer == "" {
								customer = "Unknown"
							}

							// Extract date & time from ISO timestamp
							// Expected format: "YYYY-MM-DD HH:MM:SS.ZZZ" or ISO8601.
							// This safely slices the string to get:
							// - date = first 10 characters (YYYY-MM-DD)
							// - time = characters 11 to 18 (HH:MM:SS)
							//
							// Only runs if length is at least 19 to avoid runtime panics.
							date := "-"
							time := "-"
							if len(o.CreatedAt) >= 19 {
								date = o.CreatedAt[:10]   // YYYY-MM-DD
								time = o.CreatedAt[11:19] // HH:MM:SS
							}

							rows = append(rows, Tr(
								Td(Class("px-4 py-2 border-t"), Text(o.ID)),
								Td(Class("px-4 py-2 border-t"), Text(customer)),
								Td(Class("px-4 py-2 border-t"), Text(date)),
								Td(Class("px-4 py-2 border-t"), Text(time)),
								Td(Class("px-4 py-2 border-t"), Text(FormatTZS(o.TotalCost))),
								Td(Class("px-4 py-2 border-t capitalize"), Text(o.Status)),
							))
						}
						return rows
					}()),
				),
			),
		),
	)
}

func CreateOrderForm(errorMsg string, items []model.Item) Node {
	// Build <option> nodes with item IDs and display prices
	opts := []Node{}
	for _, item := range items {
		opts = append(opts,
			Option(
				Value(item.ID),
				Text(fmt.Sprintf("%s (%.0f TZS)", item.Name, item.Price)),
			),
		)
	}

	return Layout("/orders/new", true,
		Div(
			ID("main"),
			Class("max-w-md mx-auto"),

			H2(Class("text-2xl font-semibold mb-6 text-gray-800"), Text("Create New Order")),

			If(errorMsg != "",
				Div(Class("mb-4 text-red-600 font-medium"), Text(errorMsg)),
			),

			Form(
				Action("/orders"),
				Method("POST"),

				// HTMX attributes to update live total
				Attr("hx-post", "/orders/total"),
				Attr("hx-target", "#total-display"),
				Attr("hx-trigger", "change from:#item_id, change from:#quantity"),
				Attr("hx-swap", "innerHTML"),

				Class("space-y-6"),

				// Item dropdown
				Div(
					Label(For("item_id"), Class("block mb-1 font-medium text-gray-700"), Text("Select Item")),
					Select(
						append([]Node{
							Name("item_id"),
							ID("item_id"),
							Class("w-full border border-gray-300 p-2 rounded"),
						}, opts...)..., // item options
					),
				),

				// Quantity input
				Div(
					Label(For("quantity"), Class("block mb-1 font-medium text-gray-700"), Text("Quantity")),
					Input(
						Type("number"),
						Name("quantity"),
						ID("quantity"),
						Min("1"),
						Value("1"),
						Class("w-full border border-gray-300 p-2 rounded"),
					),
				),

				//  Live total cost display (HTMX target)
				Div(
					ID("total-display"),
					Class("text-lg font-semibold text-gray-700"),
					Text("Total: 0 TZS"),
				),

				// Submit
				Button(
					Type("submit"),
					Class("bg-green-600 text-white px-4 py-2 rounded hover:bg-green-700 transition"),
					Text("Place Order"),
				),
			),
		),
	)
}

func FormatTZS(v float64) string {
	// We don’t expect decimals; round to nearest shilling
	n := int64(math.Round(v))
	return fmt.Sprintf("%s TZS", humanize.Comma(n)) // → “35,000 TZS”
}
