package html

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// NewItemPage renders the "Add New Menu Item" form page.
// This form allows the admin or staff to create a new menu item
// with a name, price, optional description, and stock quantity.
//
// Parameters:
//   - errorMessage: optional error message to display at the top of the form.

func NewItemPage(errorMessage string) Node {
	return Layout("/items/new", true,
		Div(
			ID("main"),
			Class("max-w-md mx-auto mt-12"),

			H2(Class("text-3xl font-bold mb-6 text-gray-800"), Text("Add New Menu Item")),

			If(errorMessage != "",
				Div(Class("bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-6"),
					Text(errorMessage),
				),
			),

			// Form to create a new item
			Form(
				Method("POST"),
				Attr("hx-post", "/items/new"),
				Attr("hx-target", "#main"),
				Attr("hx-swap", "innerHTML"),
				Class("space-y-6"),

				Div(
					Label(For("name"), Class("block font-medium text-gray-700 mb-1"), Text("Item Name")),
					Input(Type("text"), ID("name"), Name("name"),
						Class("w-full border border-gray-300 rounded p-2 focus:ring focus:border-indigo-500"),
						Required(),
					),
				),

				Div(
					Label(For("price"), Class("block font-medium text-gray-700 mb-1"), Text("Price (TZS)")),
					Input(Type("number"), ID("price"), Name("price"),
						Class("w-full border border-gray-300 rounded p-2 focus:ring focus:border-indigo-500"),
						Step("0.01"),
						Required(),
					),
				),

				Div(
					Label(For("description"), Class("block font-medium text-gray-700 mb-1"), Text("Description (optional)")),
					Textarea(Name("description"), ID("description"),
						Class("w-full border border-gray-300 rounded p-2 focus:ring focus:border-indigo-500"),
						Text(""),
					),
				),

				Div(
					Label(For("quantity"), Class("block font-medium text-gray-700 mb-1"), Text("Available Quantity")),
					Input(Type("number"), ID("quantity"), Name("quantity"),
						Class("w-full border border-gray-300 rounded p-2 focus:ring focus:border-indigo-500"),
						Min("0"),
						Value("1"),
					),
				),

				Button(Type("submit"),
					Class("w-full bg-indigo-600 text-white font-semibold py-2 px-4 rounded hover:bg-indigo-700 transition"),
					Text("Add Item"),
				),
			),
		),
	)
}
