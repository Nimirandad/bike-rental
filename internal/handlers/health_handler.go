package handlers

import (
	"net/http"

	"github.com/Nimirandad/bike-rental-service/internal/logger"
	"github.com/Nimirandad/bike-rental-service/internal/services"
	"github.com/Nimirandad/bike-rental-service/internal/types"
)

type HealthService interface {
	CheckHealth() (*services.HealthStatus, bool)
}

type HealthHandler struct {
	healthService HealthService
}

func NewHealthHandler(healthService *services.HealthService) *HealthHandler {
	return &HealthHandler{
		healthService: healthService,
	}
}

// CheckHealth godoc
// @Summary Health check
// @Description Check API and database health status
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} services.HealthStatus "Service is healthy"
// @Failure 503 {object} services.HealthStatus "Service is unhealthy"
// @Router /status [get]
func (h *HealthHandler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	log := logger.Get()

	health, isHealthy := h.healthService.CheckHealth()

	statusCode := http.StatusOK
	if !isHealthy {
		log.Warn().Str("status", health.Status).Str("database", health.Database).Msg("Health check failed")
		statusCode = http.StatusServiceUnavailable
	} else {
		log.Info().Str("status", health.Status).Msg("Health check passed")
	}

	types.WriteJSON(w, statusCode, health)
}
