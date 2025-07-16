package html

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// HomePage renders the main landing page for the POS system.
// Parameters:
//   - authenticated: bool indicating whether the user is logged in (used by Layout).
func HomePage(authenticated bool) Node {
	return Layout("/", authenticated,
		Section(
			Class("flex items-center justify-center min-h-[calc(100vh-4rem)] px-6"),
			Div( // tinted panel
				Class("w-full max-w-7xl mx-auto bg-gradient-to-br from-indigo-50 to-purple-100 rounded-3xl shadow-xl"),
				Div(
					Class("grid md:grid-cols-2 gap-10 p-10 md:p-16"),

					Div(Class("flex flex-col justify-center space-y-8"),
						H1(Class("text-4xl md:text-5xl font-extrabold text-indigo-800 leading-tight"),
							Text("Modern POS for Small Restaurants & Cafés"),
						),
						P(Class("text-lg md:text-xl text-gray-700"),
							Text("Track sales, manage orders, and simplify inventory – all with a clean, minimal interface powered by PocketBase and Go."),
						),
						Div(Class("flex flex-wrap gap-4 pt-2"),
							cta("/orders", "View Orders", "blue"),
							cta("/orders/new", "Create Order", "green"),
							cta("/items/new", "Add Item", "purple"),
						),
					),

					Div(Class("flex items-center justify-center")),
				),
			),
		),
	)
}

// cta generates a styled call-to-action <a> button.
func cta(href, label, color string) Node {
	return A(Href(href),
		Class("px-6 py-3 rounded-lg font-medium text-white shadow transition "+
			"bg-"+color+"-600 hover:bg-"+color+"-700 focus:ring-4 focus:ring-"+color+"-300 focus:outline-none"),
		Text(label),
	)
}
