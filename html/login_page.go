package html

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func LoginPage(errorMessage string) Node {
	return Layout("/login", false,
		Div(Class("flex items-center justify-center min-h-[80vh] px-4"),
			Div(Class("w-full max-w-sm bg-white rounded-2xl shadow-xl p-8 space-y-6"),
				Div(Class("text-center"),
					H2(Class("text-2xl font-bold text-indigo-700"), Text("Welcome Back")),
					P(Class("mt-1 text-sm text-gray-500"), Text("Sign in to continue")),
				),

				If(errorMessage != "", Div(
					Class("bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded text-sm"),
					Text(errorMessage),
				)),

				Form(
					Action("/login"),
					Method("POST"),
					Attr("hx-post", "/login"),
					Attr("hx-target", "#main"),
					Attr("hx-swap", "innerHTML"),
					Class("space-y-4"),

					Div(
						Label(Class("block text-sm font-medium text-gray-700"), Text("Email or Username")),
						Input(Type("text"), Name("identity"), Required(),
							Class("mt-1 w-full px-4 py-2 border rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-500")),
					),
					Div(
						Label(Class("block text-sm font-medium text-gray-700"), Text("Password")),
						Input(Type("password"), Name("password"), Required(),
							Class("mt-1 w-full px-4 py-2 border rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-500")),
					),

					Button(
						Type("submit"),
						Class("w-full bg-indigo-600 hover:bg-indigo-700 text-white font-semibold py-2 px-4 rounded-lg transition"),
						Text("Sign In"),
					),
				),
			),
		),
	)
}
