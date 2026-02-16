package app

import (
	"github.com/Nimirandad/bike-rental-service/internal/config"
	"github.com/Nimirandad/bike-rental-service/internal/database"
	"github.com/Nimirandad/bike-rental-service/internal/logger"
	"github.com/Nimirandad/bike-rental-service/internal/routes"
	"github.com/Nimirandad/bike-rental-service/internal/server"
)

func Run(cfg *config.Config) {
	logger.Setup(cfg.LogLevel)
	log := logger.Get()

	log.Info().Str("log_level", cfg.LogLevel).Msg("Logger initialized")

	db, err := database.Connect(cfg.SQLitePath)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()

	log.Info().Str("db_path", cfg.SQLitePath).Msg("Database connected")

	srv := server.NewServer(cfg, db.DB)
	routes.RegisterRoutes(srv)

	log.Info().Str("port", cfg.Port).Msg("Starting server")
	if err := srv.Start(cfg.Port); err != nil {
		log.Fatal().Err(err).Msg("Server failed to start")
	}
}
