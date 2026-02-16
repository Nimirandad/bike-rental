package services

import (
	"database/sql"
	"time"
)

type HealthService struct {
	db *sql.DB
}

func NewHealthService(db *sql.DB) *HealthService {
	return &HealthService{
		db: db,
	}
}

type HealthStatus struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Service   string `json:"service"`
	Version   string `json:"version"`
	Database  string `json:"database"`
}

func (s *HealthService) CheckHealth() (*HealthStatus, bool) {
	health := &HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
		Service:   "bike-rental-service",
		Version:   "1.0.0",
		Database:  "connected",
	}

	if err := s.db.Ping(); err != nil {
		health.Status = "unhealthy"
		health.Database = "disconnected"
		return health, false
	}

	return health, true
}