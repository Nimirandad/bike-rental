package server

import (
	"database/sql"
	"net/http"

	"github.com/Nimirandad/bike-rental-service/internal/config"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	Chi      *chi.Mux
	AdminChi *chi.Mux
	Config   *config.Config
	DB *sql.DB
}

func NewServer(cfg *config.Config, db *sql.DB) *Server {
	return &Server{
		Chi:      chi.NewRouter(),
		AdminChi: chi.NewRouter(),
		Config:   cfg,
		DB:       db,
	}
}

func (server *Server) Start(port string) error {
	return http.ListenAndServe(":" + port, server.Chi)
}

func (server *Server) StartAdmin(port string) error {
	return http.ListenAndServe(":" + port, server.AdminChi)
}