package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Nimirandad/bike-rental-service/internal/models"
	"github.com/Nimirandad/bike-rental-service/internal/utils"
	"github.com/stretchr/testify/assert"
)

type MockBikeService struct {
	GetAvailableBikesFunc func(page, limit int) ([]*models.Bike, int, error)
}

func (m *MockBikeService) GetAvailableBikes(page, limit int) ([]*models.Bike, int, error) {
	return m.GetAvailableBikesFunc(page, limit)
}

func TestBikeHandler_GetAvailableBikes_Success(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	testUser := &models.User{ID: 1, Email: "test@example.com", FirstName: "John", LastName: "Doe"}
	token, _ := utils.GenerateJWT(testUser)

	mockService := &MockBikeService{
		GetAvailableBikesFunc: func(page, limit int) ([]*models.Bike, int, error) {
			bikes := []*models.Bike{
				{ID: 1, Latitude: 40.416775, Longitude: -3.703790, IsAvailable: true, PricePerMinute: 0.5},
				{ID: 2, Latitude: 40.417832, Longitude: -3.705064, IsAvailable: true, PricePerMinute: 0.5},
			}
			return bikes, 10, nil
		},
	}

	handler := &BikeHandler{bikeService: mockService}

	req := httptest.NewRequest(http.MethodGet, "/api/bikes?page=1&limit=10", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler.GetAvailableBikes(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBikeHandler_GetAvailableBikes_NoAuthHeader(t *testing.T) {
	handler := &BikeHandler{}

	req := httptest.NewRequest(http.MethodGet, "/api/bikes", nil)
	w := httptest.NewRecorder()

	handler.GetAvailableBikes(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestBikeHandler_GetAvailableBikes_InvalidToken(t *testing.T) {
	handler := &BikeHandler{}

	req := httptest.NewRequest(http.MethodGet, "/api/bikes", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()

	handler.GetAvailableBikes(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestBikeHandler_GetAvailableBikes_ServiceError(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	testUser := &models.User{ID: 1, Email: "test@example.com", FirstName: "John", LastName: "Doe"}
	token, _ := utils.GenerateJWT(testUser)

	mockService := &MockBikeService{
		GetAvailableBikesFunc: func(page, limit int) ([]*models.Bike, int, error) {
			return nil, 0, errors.New("database error")
		},
	}

	handler := &BikeHandler{bikeService: mockService}

	req := httptest.NewRequest(http.MethodGet, "/api/bikes", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler.GetAvailableBikes(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestBikeHandler_GetAvailableBikes_InvalidAuthFormat(t *testing.T) {
	handler := &BikeHandler{}

	req := httptest.NewRequest(http.MethodGet, "/api/bikes", nil)
	req.Header.Set("Authorization", "InvalidFormat")
	w := httptest.NewRecorder()

	handler.GetAvailableBikes(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestBikeHandler_GetAvailableBikes_WithPagination(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	testUser := &models.User{ID: 1, Email: "test@example.com"}
	token, _ := utils.GenerateJWT(testUser)

	mockService := &MockBikeService{
		GetAvailableBikesFunc: func(page, limit int) ([]*models.Bike, int, error) {
			assert.Equal(t, 2, page)
			assert.Equal(t, 15, limit)
			return []*models.Bike{}, 0, nil
		},
	}

	handler := &BikeHandler{bikeService: mockService}

	req := httptest.NewRequest(http.MethodGet, "/api/bikes?page=2&limit=15", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler.GetAvailableBikes(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}