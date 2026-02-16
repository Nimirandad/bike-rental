package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nimirandad/bike-rental-service/internal/services"
	"github.com/stretchr/testify/assert"
)

type MockHealthService struct {
	CheckHealthFunc func() (*services.HealthStatus, bool)
}

func (m *MockHealthService) CheckHealth() (*services.HealthStatus, bool) {
	return m.CheckHealthFunc()
}

func TestHealthHandler_CheckHealth_Healthy(t *testing.T) {
	mockService := &MockHealthService{
		CheckHealthFunc: func() (*services.HealthStatus, bool) {
			return &services.HealthStatus{
				Status:    "healthy",
				Timestamp: "2024-01-01T00:00:00Z",
				Service:   "bike-rental-service",
				Version:   "1.0.0",
				Database:  "connected",
			}, true
		},
	}

	handler := &HealthHandler{healthService: mockService}

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	handler.CheckHealth(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHealthHandler_CheckHealth_Unhealthy(t *testing.T) {
	mockService := &MockHealthService{
		CheckHealthFunc: func() (*services.HealthStatus, bool) {
			return &services.HealthStatus{
				Status:    "unhealthy",
				Timestamp: "2024-01-01T00:00:00Z",
				Service:   "bike-rental-service",
				Version:   "1.0.0",
				Database:  "disconnected",
			}, false
		},
	}

	handler := &HealthHandler{healthService: mockService}

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	handler.CheckHealth(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}