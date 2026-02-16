package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Nimirandad/bike-rental-service/internal/constants"
	"github.com/Nimirandad/bike-rental-service/internal/models"
	"github.com/Nimirandad/bike-rental-service/internal/utils"
	"github.com/stretchr/testify/assert"
)

type MockRentalService struct {
	StartRentalFunc      func(userID, bikeID int) (*models.Rental, error)
	EndRentalFunc        func(userID int, endLat, endLong float64) (*models.Rental, error)
	GetRentalHistoryFunc func(userID, page, limit int) ([]*models.Rental, int, error)
}

func (m *MockRentalService) StartRental(userID, bikeID int) (*models.Rental, error) {
	return m.StartRentalFunc(userID, bikeID)
}

func (m *MockRentalService) EndRental(userID int, endLat, endLong float64) (*models.Rental, error) {
	return m.EndRentalFunc(userID, endLat, endLong)
}

func (m *MockRentalService) GetRentalHistory(userID, page, limit int) ([]*models.Rental, int, error) {
	return m.GetRentalHistoryFunc(userID, page, limit)
}

func TestRentalHandler_StartRental_Success(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	testUser := &models.User{ID: 1, Email: "test@example.com", FirstName: "John", LastName: "Doe"}
	token, _ := utils.GenerateJWT(testUser)

	mockService := &MockRentalService{
		StartRentalFunc: func(userID, bikeID int) (*models.Rental, error) {
			return &models.Rental{ID: 1, UserID: userID, BikeID: bikeID, Status: "running"}, nil
		},
	}

	handler := &RentalHandler{rentalService: mockService}

	reqBody := map[string]int{"bike_id": 1}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/rentals/start", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler.StartRental(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRentalHandler_StartRental_NoAuthHeader(t *testing.T) {
	handler := &RentalHandler{}

	req := httptest.NewRequest(http.MethodPost, "/api/rentals/start", nil)
	w := httptest.NewRecorder()

	handler.StartRental(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRentalHandler_StartRental_InvalidBikeID(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	testUser := &models.User{ID: 1, Email: "test@example.com", FirstName: "John", LastName: "Doe"}
	token, _ := utils.GenerateJWT(testUser)

	handler := &RentalHandler{}

	reqBody := map[string]int{"bike_id": 0}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/rentals/start", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler.StartRental(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRentalHandler_StartRental_UserHasActiveRental(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	testUser := &models.User{ID: 1, Email: "test@example.com", FirstName: "John", LastName: "Doe"}
	token, _ := utils.GenerateJWT(testUser)

	mockService := &MockRentalService{
		StartRentalFunc: func(userID, bikeID int) (*models.Rental, error) {
			return nil, constants.ErrUserHasActiveRental
		},
	}

	handler := &RentalHandler{rentalService: mockService}

	reqBody := map[string]int{"bike_id": 1}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/rentals/start", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler.StartRental(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestRentalHandler_StartRental_BikeNotAvailable(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	testUser := &models.User{ID: 1, Email: "test@example.com", FirstName: "John", LastName: "Doe"}
	token, _ := utils.GenerateJWT(testUser)

	mockService := &MockRentalService{
		StartRentalFunc: func(userID, bikeID int) (*models.Rental, error) {
			return nil, constants.ErrBikeNotAvailable
		},
	}

	handler := &RentalHandler{rentalService: mockService}

	reqBody := map[string]int{"bike_id": 1}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/rentals/start", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler.StartRental(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestRentalHandler_StartRental_BikeNotFound(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	testUser := &models.User{ID: 1, Email: "test@example.com", FirstName: "John", LastName: "Doe"}
	token, _ := utils.GenerateJWT(testUser)

	mockService := &MockRentalService{
		StartRentalFunc: func(userID, bikeID int) (*models.Rental, error) {
			return nil, constants.ErrBikeNotFound
		},
	}

	handler := &RentalHandler{rentalService: mockService}

	reqBody := map[string]int{"bike_id": 999}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/rentals/start", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler.StartRental(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRentalHandler_EndRental_Success(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	testUser := &models.User{ID: 1, Email: "test@example.com", FirstName: "John", LastName: "Doe"}
	token, _ := utils.GenerateJWT(testUser)

	mockService := &MockRentalService{
		EndRentalFunc: func(userID int, endLat, endLong float64) (*models.Rental, error) {
			return &models.Rental{ID: 1, UserID: userID, Status: "ended"}, nil
		},
	}

	handler := &RentalHandler{rentalService: mockService}

	reqBody := map[string]float64{"latitude": 40.416775, "longitude": -3.703790}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/rentals/end", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler.EndRental(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRentalHandler_EndRental_InvalidLatitude(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	testUser := &models.User{ID: 1, Email: "test@example.com", FirstName: "John", LastName: "Doe"}
	token, _ := utils.GenerateJWT(testUser)

	handler := &RentalHandler{}

	reqBody := map[string]float64{"latitude": 91.0, "longitude": -3.703790}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/rentals/end", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler.EndRental(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRentalHandler_EndRental_InvalidLongitude(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	testUser := &models.User{ID: 1, Email: "test@example.com", FirstName: "John", LastName: "Doe"}
	token, _ := utils.GenerateJWT(testUser)

	handler := &RentalHandler{}

	reqBody := map[string]float64{"latitude": 40.416775, "longitude": 181.0}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/rentals/end", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler.EndRental(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRentalHandler_EndRental_NoActiveRental(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	testUser := &models.User{ID: 1, Email: "test@example.com", FirstName: "John", LastName: "Doe"}
	token, _ := utils.GenerateJWT(testUser)

	mockService := &MockRentalService{
		EndRentalFunc: func(userID int, endLat, endLong float64) (*models.Rental, error) {
			return nil, constants.ErrNoActiveRental
		},
	}

	handler := &RentalHandler{rentalService: mockService}

	reqBody := map[string]float64{"latitude": 40.416775, "longitude": -3.703790}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/rentals/end", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler.EndRental(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestRentalHandler_EndRental_LocationTooFar(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	testUser := &models.User{ID: 1, Email: "test@example.com", FirstName: "John", LastName: "Doe"}
	token, _ := utils.GenerateJWT(testUser)

	mockService := &MockRentalService{
		EndRentalFunc: func(userID int, endLat, endLong float64) (*models.Rental, error) {
			return nil, constants.ErrEndLocationTooFar
		},
	}

	handler := &RentalHandler{rentalService: mockService}

	reqBody := map[string]float64{"latitude": 48.856614, "longitude": 2.352222}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/rentals/end", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler.EndRental(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRentalHandler_EndRental_InternalError(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	testUser := &models.User{ID: 1, Email: "test@example.com", FirstName: "John", LastName: "Doe"}
	token, _ := utils.GenerateJWT(testUser)

	mockService := &MockRentalService{
		EndRentalFunc: func(userID int, endLat, endLong float64) (*models.Rental, error) {
			return nil, errors.New("database error")
		},
	}

	handler := &RentalHandler{rentalService: mockService}

	reqBody := map[string]float64{"latitude": 40.416775, "longitude": -3.703790}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/rentals/end", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler.EndRental(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestRentalHandler_GetRentalHistory_Success(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	testUser := &models.User{ID: 1, Email: "test@example.com", FirstName: "John", LastName: "Doe"}
	token, _ := utils.GenerateJWT(testUser)

	mockService := &MockRentalService{
		GetRentalHistoryFunc: func(userID, page, limit int) ([]*models.Rental, int, error) {
			rentals := []*models.Rental{
				{ID: 1, UserID: userID, BikeID: 1, Status: "ended"},
				{ID: 2, UserID: userID, BikeID: 2, Status: "running"},
			}
			return rentals, 2, nil
		},
	}

	handler := &RentalHandler{rentalService: mockService}

	req := httptest.NewRequest(http.MethodGet, "/api/rentals?page=1&limit=10", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler.GetRentalHistory(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRentalHandler_StartRental_InternalError(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	testUser := &models.User{ID: 1, Email: "test@example.com", FirstName: "John", LastName: "Doe"}
	token, _ := utils.GenerateJWT(testUser)

	mockService := &MockRentalService{
		StartRentalFunc: func(userID, bikeID int) (*models.Rental, error) {
			return nil, errors.New("database error")
		},
	}

	handler := &RentalHandler{rentalService: mockService}

	reqBody := map[string]int{"bike_id": 1}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/rentals/start", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler.StartRental(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestRentalHandler_StartRental_InvalidJSON(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	testUser := &models.User{ID: 1, Email: "test@example.com"}
	token, _ := utils.GenerateJWT(testUser)

	handler := &RentalHandler{}
	req := httptest.NewRequest(http.MethodPost, "/api/rentals/start", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler.StartRental(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRentalHandler_EndRental_InvalidJSON(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	testUser := &models.User{ID: 1, Email: "test@example.com"}
	token, _ := utils.GenerateJWT(testUser)

	handler := &RentalHandler{}
	req := httptest.NewRequest(http.MethodPost, "/api/rentals/end", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler.EndRental(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRentalHandler_GetRentalHistory_ServiceError(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	testUser := &models.User{ID: 1, Email: "test@example.com"}
	token, _ := utils.GenerateJWT(testUser)

	mockService := &MockRentalService{
		GetRentalHistoryFunc: func(userID, page, limit int) ([]*models.Rental, int, error) {
			return nil, 0, errors.New("database error")
		},
	}

	handler := &RentalHandler{rentalService: mockService}
	req := httptest.NewRequest(http.MethodGet, "/api/rentals/history", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler.GetRentalHistory(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestRentalHandler_StartRental_InvalidAuthFormat(t *testing.T) {
	handler := &RentalHandler{}
	reqBody := map[string]int{"bike_id": 1}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/rentals/start", bytes.NewReader(body))
	req.Header.Set("Authorization", "InvalidFormat")
	w := httptest.NewRecorder()

	handler.StartRental(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRentalHandler_EndRental_InvalidAuthFormat(t *testing.T) {
	handler := &RentalHandler{}
	reqBody := map[string]float64{"latitude": 40.416775, "longitude": -3.703790}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/rentals/end", bytes.NewReader(body))
	req.Header.Set("Authorization", "InvalidFormat")
	w := httptest.NewRecorder()

	handler.EndRental(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRentalHandler_GetRentalHistory_InvalidAuthFormat(t *testing.T) {
	handler := &RentalHandler{}
	req := httptest.NewRequest(http.MethodGet, "/api/rentals/history", nil)
	req.Header.Set("Authorization", "InvalidFormat")
	w := httptest.NewRecorder()

	handler.GetRentalHistory(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRentalHandler_GetRentalHistory_WithPagination(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	testUser := &models.User{ID: 1}
	token, _ := utils.GenerateJWT(testUser)

	mockService := &MockRentalService{
		GetRentalHistoryFunc: func(userID, page, limit int) ([]*models.Rental, int, error) {
			return []*models.Rental{
				{ID: 1, UserID: userID, BikeID: 1},
			}, 1, nil
		},
	}

	handler := &RentalHandler{rentalService: mockService}
	req := httptest.NewRequest(http.MethodGet, "/rentals/history?page=2&limit=5", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler.GetRentalHistory(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
