package routes

import (
	"github.com/Nimirandad/bike-rental-service/internal/handlers"
	"github.com/Nimirandad/bike-rental-service/internal/repositories"
	"github.com/Nimirandad/bike-rental-service/internal/server"
	"github.com/Nimirandad/bike-rental-service/internal/server/middlewares"
	"github.com/Nimirandad/bike-rental-service/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	_ "github.com/Nimirandad/bike-rental-service/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

func RegisterRoutes(s *server.Server) {
	s.Chi.Use(middlewares.CORS)
	s.Chi.Use(middlewares.LoggingMiddleware)
	s.Chi.Use(middleware.Recoverer)

	userRepo := repositories.NewUserRepository(s.DB)
	bikeRepo := repositories.NewBikeRepository(s.DB)
	adminRepo := repositories.NewAdminRepository(s.DB)
	rentalRepo := repositories.NewRentalRepository(s.DB)

	userService := services.NewUserService(userRepo)
	bikeService := services.NewBikeService(bikeRepo)
	rentalService := services.NewRentalService(rentalRepo, bikeRepo)
	adminService := services.NewAdminService(adminRepo)
	healthService := services.NewHealthService(s.DB)

	userHandler := handlers.NewUserHandler(userService)
	bikeHandler := handlers.NewBikeHandler(bikeService)
	rentalHandler := handlers.NewRentalHandler(rentalService)
	adminHandler := handlers.NewAdminHandler(adminService)
	healthHandler := handlers.NewHealthHandler(healthService)

	s.Chi.Get("/status", healthHandler.CheckHealth)
	s.Chi.Get("/swagger/*", httpSwagger.WrapHandler)

	s.Chi.Route("/api/v1", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/register", userHandler.RegisterUser)
			r.Post("/login", userHandler.LoginUser)
			r.Get("/profile", userHandler.GetUserProfile)
			r.Patch("/profile", userHandler.UpdateUserProfile)
		})

		r.Route("/bikes", func(r chi.Router) {
			r.Get("/available", bikeHandler.GetAvailableBikes)
		})

		r.Route("/rentals", func(r chi.Router) {
			r.Post("/start", rentalHandler.StartRental)
			r.Post("/end", rentalHandler.EndRental)
			r.Get("/history", rentalHandler.GetRentalHistory)
		})

		r.Route("/admin", func(r chi.Router) {
			r.Route("/bikes", func(r chi.Router) {
				r.Get("/", adminHandler.GetAllBikes)
				r.Post("/", adminHandler.AddBike)
				r.Patch("/{bike-id}", adminHandler.UpdateBike)
			})

			r.Route("/users", func(r chi.Router) {
				r.Get("/", adminHandler.GetAllUsers)
				r.Get("/{user-id}", adminHandler.GetUserDetails)
				r.Patch("/{user-id}", adminHandler.UpdateUser)
			})

			r.Route("/rentals", func(r chi.Router) {
				r.Get("/", adminHandler.GetAllRentals)
				r.Get("/{rental-id}", adminHandler.GetRentalDetails)
				r.Patch("/{rental-id}", adminHandler.UpdateRental)
			})

		})
	})
}
