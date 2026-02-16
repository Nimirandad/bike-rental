package routes

import (
	"net/http"

	"github.com/Nimirandad/bike-rental-service/internal/middlewares"
	server "github.com/Nimirandad/bike-rental-service/internal/server"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func RegisterRoutes(s *server.Server) {
	s.Chi.Use(middlewares.CORS)
	s.Chi.Use(middleware.Logger)
	s.Chi.Use(middleware.Recoverer)

	// Register routes for the main server
	s.Chi.Route("/api/v1", func(r chi.Router) {

		r.Get("/bikes", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("List of bikes"))
		})
		r.Post("/rent", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Rent a bike"))
		})
	})
}
