// Package HTML holds all the common HTML components and utilities.
package html

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

var hashOnce sync.Once
var appCSSPath, appJSPath, htmxJSPath string

// PageProps are properties for the [page] component.
type PageProps struct {
	Title       string
	Description string
}

func Layout(path string, authenticated bool, children ...Node) Node {
	// Run only once to compute file paths
	hashOnce.Do(func() {
		appCSSPath = getHashedPath("public/styles/app.css")
		htmxJSPath = getHashedPath("public/scripts/htmx.js")
		appJSPath = getHashedPath("public/scripts/app.js")
	})

	isLoginPage := strings.HasPrefix(path, "/login")

	return HTML5(HTML5Props{
		Language:    "en",
		Title:       "POS System",
		Description: "Modern point‑of‑sale system built with PocketBase and Gomponents",
		Head: []Node{
			Meta(Name("viewport"), Content("width=device-width, initial-scale=1")), // 📱 responsive
			Script(Src("https://cdn.tailwindcss.com?plugins=typography,forms")),
			Script(Src(htmxJSPath), Defer()),
			Script(Src(appJSPath), Defer()),
		},
		Body: []Node{
			Div(Class("min-h-screen bg-gray-50 font-sans flex flex-col"),
				If(!isLoginPage,
					Header(
						Class("h-16 flex items-center justify-between px-6 bg-gradient-to-r from-indigo-700 via-purple-700 to-pink-700 text-white shadow-md sticky top-0 z-50"),
						H1(Class("text-lg md:text-xl font-bold tracking-tight select-none"), Text("POS System")),
						Nav(Class("flex space-x-4 text-sm font-medium"),
							navLink("/", "Home"),
							navLink("/orders", "Orders"),
							navLink("/orders/new", "New Order"),
							navLink("/items/new", "Add Item"),
							If(authenticated,
								navLink("/logout", "Logout"),
							),
						),
					),
				),
				Main(ID("main"), Class("flex-grow"), Group(children)),
			),
		},
	})
}

func navLink(href, label string) Node {
	return A(Href(href),
		Class("transition hover:text-yellow-300 whitespace-nowrap"), Text(label),
	)
}

// header bar with logo and navigation.
func header() Node {
	return Div(Class("bg-indigo-600 text-white shadow-sm"),
		container(true, false,
			Div(Class("h-16 flex items-center justify-between"),
				A(Href("/"), Class("inline-flex items-center text-xl font-semibold"),
					Img(Src("/images/logo.png"), Alt("Logo"), Class("h-12 w-auto bg-white rounded-full mr-4")),
					Text("Home"),
				),
			),
		),
	)
}

// container restricts the width and sets padding.
func container(padX, padY bool, children ...Node) Node {
	return Div(
		Classes{
			"max-w-7xl mx-auto px":  true,
			"px-4 md:px-8 lg:px-16": padX,
			"py-4 md:py-8":          padY,
		},
		Group(children),
	)
}

func getHashedPath(path string) string {
	externalPath := strings.TrimPrefix(path, "public")
	ext := filepath.Ext(path)
	if ext == "" {
		panic("no extension found")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Sprintf("%v.x%v", strings.TrimSuffix(externalPath, ext), ext)
	}

	return fmt.Sprintf("%v.%x%v", strings.TrimSuffix(externalPath, ext), sha256.Sum256(data), ext)
}
