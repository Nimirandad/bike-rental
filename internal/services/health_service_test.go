package services

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// TestHealthService_CheckHealth_Healthy tests successful health check with database connected
func TestHealthService_CheckHealth_Healthy(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPing()

	service := NewHealthService(db)
	health, isHealthy := service.CheckHealth()

	assert.True(t, isHealthy)
	assert.NotNil(t, health)
	assert.Equal(t, "healthy", health.Status)
	assert.Equal(t, "bike-rental-service", health.Service)
	assert.Equal(t, "1.0.0", health.Version)
	assert.Equal(t, "connected", health.Database)
	assert.NotEmpty(t, health.Timestamp)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestHealthService_CheckHealth_Unhealthy tests health check when database is disconnected
func TestHealthService_CheckHealth_Unhealthy(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPing().WillReturnError(errors.New("connection refused"))

	service := NewHealthService(db)
	health, isHealthy := service.CheckHealth()

	assert.False(t, isHealthy)
	assert.NotNil(t, health)
	assert.Equal(t, "unhealthy", health.Status)
	assert.Equal(t, "bike-rental-service", health.Service)
	assert.Equal(t, "1.0.0", health.Version)
	assert.Equal(t, "disconnected", health.Database)
	assert.NotEmpty(t, health.Timestamp)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestHealthService_CheckHealth_DatabaseNil tests health check with nil database (edge case)
func TestHealthService_CheckHealth_DatabaseNil(t *testing.T) {
	var db *sql.DB = nil

	service := NewHealthService(db)

	// This should panic or return unhealthy - we expect a panic for nil pointer
	assert.Panics(t, func() {
		service.CheckHealth()
	})
}