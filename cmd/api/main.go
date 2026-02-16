package main

import (
	"github.com/Nimirandad/bike-rental-service/internal/app"
	"github.com/Nimirandad/bike-rental-service/internal/config"
)

// @title Bike Rental Service API
// @version 1.0
// @description API para gestión de alquiler de bicicletas
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@bikerental.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @securityDefinitions.basic BasicAuth

func main() {
	cfg := config.Load()

	app.Run(&cfg)
}
