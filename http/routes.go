package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"maragu.dev/httph"
)

func (s *Server) setupRoutes() {
	s.mux.Group(func(r chi.Router) {
		r.Use(middleware.Compress(5))
		r.Group(func(r chi.Router) {
			r.Use(httph.VersionedAssets)

			Static(r)
		})

		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
		))

		// r.Get("/docs/*", httpSwagger.Handler(
		// 	httpSwagger.URL("/docs/swagger.json"),
		// ))

		Home(r)
		Auth(r)
		OrderRoutes(r)
		ItemRoutes(r)

	})
}
