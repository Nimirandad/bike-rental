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
	"github.com/stretchr/testify/assert"
)

type MockAdminService2 struct {
	CreateBikeFunc    func(latitude, longitude, pricePerMinute float64) (*models.Bike, error)
	GetAllBikesFunc   func(page, limit int) ([]*models.Bike, int, error)
	UpdateBikeFunc    func(bikeID int, latitude, longitude *float64, isAvailable *bool, pricePerMinute *float64) (*models.Bike, error)
	GetAllUsersFunc   func(page, limit int) ([]*models.User, int, error)
	GetUserByIDFunc   func(userID int) (*models.User, error)
	UpdateUserFunc    func(userID int, email, firstName, lastName, hashedPassword *string) (*models.User, error)
	GetAllRentalsFunc func(page, limit int) ([]*models.Rental, int, error)
	GetRentalByIDFunc func(rentalID int) (*models.Rental, error)
	UpdateRentalFunc  func(rentalID int, status *string) (*models.Rental, error)
}

func (m *MockAdminService2) CreateBike(latitude, longitude, pricePerMinute float64) (*models.Bike, error) {
	return m.CreateBikeFunc(latitude, longitude, pricePerMinute)
}

func (m *MockAdminService2) GetAllBikes(page, limit int) ([]*models.Bike, int, error) {
	return m.GetAllBikesFunc(page, limit)
}

func (m *MockAdminService2) UpdateBike(bikeID int, latitude, longitude *float64, isAvailable *bool, pricePerMinute *float64) (*models.Bike, error) {
	return m.UpdateBikeFunc(bikeID, latitude, longitude, isAvailable, pricePerMinute)
}

func (m *MockAdminService2) GetAllUsers(page, limit int) ([]*models.User, int, error) {
	return m.GetAllUsersFunc(page, limit)
}

func (m *MockAdminService2) GetUserByID(userID int) (*models.User, error) {
	return m.GetUserByIDFunc(userID)
}

func (m *MockAdminService2) UpdateUser(userID int, email, firstName, lastName, hashedPassword *string) (*models.User, error) {
	return m.UpdateUserFunc(userID, email, firstName, lastName, hashedPassword)
}

func (m *MockAdminService2) GetAllRentals(page, limit int) ([]*models.Rental, int, error) {
	return m.GetAllRentalsFunc(page, limit)
}

func (m *MockAdminService2) GetRentalByID(rentalID int) (*models.Rental, error) {
	return m.GetRentalByIDFunc(rentalID)
}

func (m *MockAdminService2) UpdateRental(rentalID int, status *string) (*models.Rental, error) {
	return m.UpdateRentalFunc(rentalID, status)
}

func TestAdminHandler_AddBike_Success(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	mockService := &MockAdminService2{
		CreateBikeFunc: func(latitude, longitude, pricePerMinute float64) (*models.Bike, error) {
			return &models.Bike{ID: 1, Latitude: latitude, Longitude: longitude, IsAvailable: true, PricePerMinute: pricePerMinute}, nil
		},
	}

	handler := &AdminHandler{adminService: mockService}
	body, _ := json.Marshal(map[string]interface{}{"latitude": 40.416775, "longitude": -3.703790, "pricePerMinute": 0.5})
	req := httptest.NewRequest(http.MethodPost, "/admin/bikes", bytes.NewReader(body))
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=") // base64("admin:admin123")
	w := httptest.NewRecorder()

	handler.AddBike(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAdminHandler_AddBike_Unauthorized(t *testing.T) {
	handler := &AdminHandler{}
	req := httptest.NewRequest(http.MethodPost, "/admin/bikes", nil)
	w := httptest.NewRecorder()

	handler.AddBike(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAdminHandler_AddBike_InvalidLatitude(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	handler := &AdminHandler{}
	body, _ := json.Marshal(map[string]interface{}{"latitude": 91.0, "longitude": -3.703790})
	req := httptest.NewRequest(http.MethodPost, "/admin/bikes", bytes.NewReader(body))
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.AddBike(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAdminHandler_AddBike_InvalidLongitude(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	handler := &AdminHandler{}
	body, _ := json.Marshal(map[string]interface{}{"latitude": 40.416775, "longitude": 181.0})
	req := httptest.NewRequest(http.MethodPost, "/admin/bikes", bytes.NewReader(body))
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.AddBike(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAdminHandler_AddBike_InvalidPrice(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	handler := &AdminHandler{}
	price := -0.5
	body, _ := json.Marshal(map[string]interface{}{"latitude": 40.416775, "longitude": -3.703790, "price_per_minute": price})
	req := httptest.NewRequest(http.MethodPost, "/admin/bikes", bytes.NewReader(body))
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.AddBike(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAdminHandler_AddBike_ServiceError(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	mockService := &MockAdminService2{
		CreateBikeFunc: func(latitude, longitude, pricePerMinute float64) (*models.Bike, error) {
			return nil, errors.New("database error")
		},
	}

	handler := &AdminHandler{adminService: mockService}
	body, _ := json.Marshal(map[string]interface{}{"latitude": 40.416775, "longitude": -3.703790})
	req := httptest.NewRequest(http.MethodPost, "/admin/bikes", bytes.NewReader(body))
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.AddBike(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAdminHandler_GetAllBikes_Success(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	mockService := &MockAdminService2{
		GetAllBikesFunc: func(page, limit int) ([]*models.Bike, int, error) {
			bikes := []*models.Bike{
				{ID: 1, Latitude: 40.416775, Longitude: -3.703790, IsAvailable: true, PricePerMinute: 0.5},
				{ID: 2, Latitude: 40.417832, Longitude: -3.705064, IsAvailable: false, PricePerMinute: 0.5},
			}
			return bikes, 10, nil
		},
	}

	handler := &AdminHandler{adminService: mockService}
	req := httptest.NewRequest(http.MethodGet, "/admin/bikes?page=1&limit=10", nil)
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.GetAllBikes(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAdminHandler_GetAllBikes_Unauthorized(t *testing.T) {
	handler := &AdminHandler{}
	req := httptest.NewRequest(http.MethodGet, "/admin/bikes", nil)
	w := httptest.NewRecorder()

	handler.GetAllBikes(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAdminHandler_GetAllUsers_Success(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	mockService := &MockAdminService2{
		GetAllUsersFunc: func(page, limit int) ([]*models.User, int, error) {
			users := []*models.User{
				{ID: 1, Email: "user1@example.com", FirstName: "John", LastName: "Doe"},
				{ID: 2, Email: "user2@example.com", FirstName: "Jane", LastName: "Smith"},
			}
			return users, 10, nil
		},
	}

	handler := &AdminHandler{adminService: mockService}
	req := httptest.NewRequest(http.MethodGet, "/admin/users?page=1&limit=10", nil)
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.GetAllUsers(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAdminHandler_GetAllUsers_Unauthorized(t *testing.T) {
	handler := &AdminHandler{}
	req := httptest.NewRequest(http.MethodGet, "/admin/users", nil)
	w := httptest.NewRecorder()

	handler.GetAllUsers(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAdminHandler_GetAllRentals_Success(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	mockService := &MockAdminService2{
		GetAllRentalsFunc: func(page, limit int) ([]*models.Rental, int, error) {
			rentals := []*models.Rental{
				{ID: 1, UserID: 1, BikeID: 1, Status: "running"},
				{ID: 2, UserID: 2, BikeID: 2, Status: "ended"},
			}
			return rentals, 10, nil
		},
	}

	handler := &AdminHandler{adminService: mockService}
	req := httptest.NewRequest(http.MethodGet, "/admin/rentals?page=1&limit=10", nil)
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.GetAllRentals(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAdminHandler_GetAllRentals_Unauthorized(t *testing.T) {
	handler := &AdminHandler{}
	req := httptest.NewRequest(http.MethodGet, "/admin/rentals", nil)
	w := httptest.NewRecorder()

	handler.GetAllRentals(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAdminHandler_UpdateBike_Success(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	latitude := 40.416775
	price := 0.75
	mockService := &MockAdminService2{
		UpdateBikeFunc: func(bikeID int, lat, long *float64, isAvailable *bool, pricePerMinute *float64) (*models.Bike, error) {
			bike := &models.Bike{ID: bikeID, Latitude: 40.0, PricePerMinute: 0.5}
			if lat != nil {
				bike.Latitude = *lat
			}
			if pricePerMinute != nil {
				bike.PricePerMinute = *pricePerMinute
			}
			return bike, nil
		},
	}

	handler := &AdminHandler{adminService: mockService}
	body, _ := json.Marshal(map[string]interface{}{"latitude": latitude, "pricePerMinute": price})
	req := httptest.NewRequest(http.MethodPut, "/admin/bikes/1", bytes.NewReader(body))
	req.SetPathValue("bike-id", "1")
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.UpdateBike(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAdminHandler_UpdateBike_InvalidBikeID(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	handler := &AdminHandler{}
	body, _ := json.Marshal(map[string]interface{}{"latitude": 40.416775})
	req := httptest.NewRequest(http.MethodPut, "/admin/bikes/abc", bytes.NewReader(body))
	req.SetPathValue("bike-id", "abc")
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.UpdateBike(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAdminHandler_UpdateBike_NoFieldsProvided(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	handler := &AdminHandler{}
	body, _ := json.Marshal(map[string]interface{}{})
	req := httptest.NewRequest(http.MethodPut, "/admin/bikes/1", bytes.NewReader(body))
	req.SetPathValue("bike-id", "1")
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.UpdateBike(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAdminHandler_UpdateBike_InvalidLatitude(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	handler := &AdminHandler{}
	latitude := 100.0
	body, _ := json.Marshal(map[string]interface{}{"latitude": latitude})
	req := httptest.NewRequest(http.MethodPut, "/admin/bikes/1", bytes.NewReader(body))
	req.SetPathValue("bike-id", "1")
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.UpdateBike(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAdminHandler_UpdateBike_NotFound(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	mockService := &MockAdminService2{
		UpdateBikeFunc: func(bikeID int, lat, long *float64, isAvailable *bool, pricePerMinute *float64) (*models.Bike, error) {
			return nil, errors.New("bike with id 99 not found")
		},
	}

	handler := &AdminHandler{adminService: mockService}
	latitude := 40.416775
	body, _ := json.Marshal(map[string]interface{}{"latitude": latitude})
	req := httptest.NewRequest(http.MethodPut, "/admin/bikes/99", bytes.NewReader(body))
	req.SetPathValue("bike-id", "99")
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.UpdateBike(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAdminHandler_UpdateUser_Success(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	email := "newemail@example.com"
	mockService := &MockAdminService2{
		UpdateUserFunc: func(userID int, em, firstName, lastName, hashedPassword *string) (*models.User, error) {
			user := &models.User{ID: userID, Email: "old@example.com"}
			if em != nil {
				user.Email = *em
			}
			return user, nil
		},
	}

	handler := &AdminHandler{adminService: mockService}
	body, _ := json.Marshal(map[string]interface{}{"email": email})
	req := httptest.NewRequest(http.MethodPut, "/admin/users/1", bytes.NewReader(body))
	req.SetPathValue("user-id", "1")
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.UpdateUser(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAdminHandler_UpdateUser_NoFieldsProvided(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	handler := &AdminHandler{}
	body, _ := json.Marshal(map[string]interface{}{})
	req := httptest.NewRequest(http.MethodPut, "/admin/users/1", bytes.NewReader(body))
	req.SetPathValue("user-id", "1")
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.UpdateUser(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAdminHandler_UpdateRental_Success(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	status := "ended"
	mockService := &MockAdminService2{
		UpdateRentalFunc: func(rentalID int, st *string) (*models.Rental, error) {
			rental := &models.Rental{ID: rentalID, Status: "running"}
			if st != nil {
				rental.Status = *st
			}
			return rental, nil
		},
	}

	handler := &AdminHandler{adminService: mockService}
	body, _ := json.Marshal(map[string]interface{}{"status": status})
	req := httptest.NewRequest(http.MethodPut, "/admin/rentals/1", bytes.NewReader(body))
	req.SetPathValue("rental-id", "1")
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.UpdateRental(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAdminHandler_UpdateRental_InvalidStatus(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	status := "invalid"
	mockService := &MockAdminService2{
		UpdateRentalFunc: func(rentalID int, st *string) (*models.Rental, error) {
			return nil, errors.New("invalid status")
		},
	}

	handler := &AdminHandler{adminService: mockService}
	body, _ := json.Marshal(map[string]interface{}{"status": status})
	req := httptest.NewRequest(http.MethodPut, "/admin/rentals/1", bytes.NewReader(body))
	req.SetPathValue("rental-id", "1")
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.UpdateRental(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAdminHandler_GetUserDetails_Success(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	mockService := &MockAdminService2{
		GetUserByIDFunc: func(userID int) (*models.User, error) {
			return &models.User{ID: userID, Email: "user@example.com"}, nil
		},
	}

	handler := &AdminHandler{adminService: mockService}
	req := httptest.NewRequest(http.MethodGet, "/admin/users/1", nil)
	req.SetPathValue("user-id", "1")
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.GetUserDetails(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAdminHandler_GetRentalDetails_Success(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	mockService := &MockAdminService2{
		GetRentalByIDFunc: func(rentalID int) (*models.Rental, error) {
			return &models.Rental{ID: rentalID, Status: "running"}, nil
		},
	}

	handler := &AdminHandler{adminService: mockService}
	req := httptest.NewRequest(http.MethodGet, "/admin/rentals/1", nil)
	req.SetPathValue("rental-id", "1")
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.GetRentalDetails(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAdminHandler_GetUserDetails_InvalidUserID(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	handler := &AdminHandler{}
	req := httptest.NewRequest(http.MethodGet, "/admin/users/abc", nil)
	req.SetPathValue("user-id", "abc")
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.GetUserDetails(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAdminHandler_GetUserDetails_NotFound(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	mockService := &MockAdminService2{
		GetUserByIDFunc: func(userID int) (*models.User, error) {
			return nil, errors.New("user not found")
		},
	}

	handler := &AdminHandler{adminService: mockService}
	req := httptest.NewRequest(http.MethodGet, "/admin/users/999", nil)
	req.SetPathValue("user-id", "999")
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.GetUserDetails(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAdminHandler_GetRentalDetails_InvalidRentalID(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	handler := &AdminHandler{}
	req := httptest.NewRequest(http.MethodGet, "/admin/rentals/abc", nil)
	req.SetPathValue("rental-id", "abc")
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.GetRentalDetails(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAdminHandler_UpdateBike_InvalidLongitude(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	handler := &AdminHandler{}
	longitude := 200.0
	body, _ := json.Marshal(map[string]interface{}{"longitude": longitude})
	req := httptest.NewRequest(http.MethodPut, "/admin/bikes/1", bytes.NewReader(body))
	req.SetPathValue("bike-id", "1")
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.UpdateBike(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAdminHandler_UpdateBike_InvalidPrice(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	handler := &AdminHandler{}
	price := -1.0
	body, _ := json.Marshal(map[string]interface{}{"pricePerMinute": price})
	req := httptest.NewRequest(http.MethodPut, "/admin/bikes/1", bytes.NewReader(body))
	req.SetPathValue("bike-id", "1")
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.UpdateBike(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAdminHandler_UpdateBike_ServiceError(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	mockService := &MockAdminService2{
		UpdateBikeFunc: func(bikeID int, lat, long *float64, isAvailable *bool, pricePerMinute *float64) (*models.Bike, error) {
			return nil, errors.New("database error")
		},
	}

	handler := &AdminHandler{adminService: mockService}
	latitude := 40.416775
	body, _ := json.Marshal(map[string]interface{}{"latitude": latitude})
	req := httptest.NewRequest(http.MethodPut, "/admin/bikes/1", bytes.NewReader(body))
	req.SetPathValue("bike-id", "1")
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.UpdateBike(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAdminHandler_UpdateUser_InvalidUserID(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	handler := &AdminHandler{}
	body, _ := json.Marshal(map[string]interface{}{"email": "new@example.com"})
	req := httptest.NewRequest(http.MethodPut, "/admin/users/abc", bytes.NewReader(body))
	req.SetPathValue("user-id", "abc")
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.UpdateUser(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAdminHandler_UpdateUser_EmailConflict(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	mockService := &MockAdminService2{
		UpdateUserFunc: func(userID int, em, firstName, lastName, hashedPassword *string) (*models.User, error) {
			return nil, constants.ErrEmailAlreadyExists
		},
	}

	handler := &AdminHandler{adminService: mockService}
	body, _ := json.Marshal(map[string]interface{}{"email": "existing@example.com"})
	req := httptest.NewRequest(http.MethodPut, "/admin/users/1", bytes.NewReader(body))
	req.SetPathValue("user-id", "1")
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.UpdateUser(w, req)
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestAdminHandler_UpdateUser_ServiceError(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	mockService := &MockAdminService2{
		UpdateUserFunc: func(userID int, em, firstName, lastName, hashedPassword *string) (*models.User, error) {
			return nil, errors.New("database error")
		},
	}

	handler := &AdminHandler{adminService: mockService}
	body, _ := json.Marshal(map[string]interface{}{"email": "new@example.com"})
	req := httptest.NewRequest(http.MethodPut, "/admin/users/1", bytes.NewReader(body))
	req.SetPathValue("user-id", "1")
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.UpdateUser(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAdminHandler_UpdateRental_InvalidRentalID(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	handler := &AdminHandler{}
	body, _ := json.Marshal(map[string]interface{}{"status": "ended"})
	req := httptest.NewRequest(http.MethodPut, "/admin/rentals/abc", bytes.NewReader(body))
	req.SetPathValue("rental-id", "abc")
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.UpdateRental(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAdminHandler_UpdateRental_NoStatus(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	handler := &AdminHandler{}
	body, _ := json.Marshal(map[string]interface{}{})
	req := httptest.NewRequest(http.MethodPut, "/admin/rentals/1", bytes.NewReader(body))
	req.SetPathValue("rental-id", "1")
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.UpdateRental(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAdminHandler_GetRentalDetails_NotFound(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	mockService := &MockAdminService2{
		GetRentalByIDFunc: func(rentalID int) (*models.Rental, error) {
			return nil, errors.New("rental not found")
		},
	}

	handler := &AdminHandler{adminService: mockService}
	req := httptest.NewRequest(http.MethodGet, "/admin/rentals/999", nil)
	req.SetPathValue("rental-id", "999")
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.GetRentalDetails(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAdminHandler_GetRentalDetails_Unauthorized(t *testing.T) {
	handler := &AdminHandler{}
	req := httptest.NewRequest(http.MethodGet, "/admin/rentals/1", nil)
	req.SetPathValue("rental-id", "1")
	w := httptest.NewRecorder()

	handler.GetRentalDetails(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAdminHandler_UpdateUser_WithPassword(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	mockService := &MockAdminService2{
		UpdateUserFunc: func(userID int, em, firstName, lastName, hashedPassword *string) (*models.User, error) {
			user := &models.User{ID: userID, Email: "test@example.com"}
			if firstName != nil {
				user.FirstName = *firstName
			}
			return user, nil
		},
	}

	handler := &AdminHandler{adminService: mockService}
	body, _ := json.Marshal(map[string]interface{}{"first_name": "John", "password": "newpassword123"})
	req := httptest.NewRequest(http.MethodPut, "/admin/users/1", bytes.NewReader(body))
	req.SetPathValue("user-id", "1")
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.UpdateUser(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAdminHandler_UpdateRental_InvalidJSON(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	handler := &AdminHandler{}
	req := httptest.NewRequest(http.MethodPut, "/admin/rentals/1", bytes.NewReader([]byte("invalid")))
	req.SetPathValue("rental-id", "1")
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.UpdateRental(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAdminHandler_UpdateBike_InvalidJSON(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	handler := &AdminHandler{}
	req := httptest.NewRequest(http.MethodPut, "/admin/bikes/1", bytes.NewReader([]byte("invalid")))
	req.SetPathValue("bike-id", "1")
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.UpdateBike(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAdminHandler_UpdateUser_InvalidJSON(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	handler := &AdminHandler{}
	req := httptest.NewRequest(http.MethodPut, "/admin/users/1", bytes.NewReader([]byte("invalid")))
	req.SetPathValue("user-id", "1")
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.UpdateUser(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAdminHandler_AddBike_InvalidJSON(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	handler := &AdminHandler{}
	req := httptest.NewRequest(http.MethodPost, "/admin/bikes", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.AddBike(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAdminHandler_GetAllBikes_ServiceError(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	mockService := &MockAdminService2{
		GetAllBikesFunc: func(page, limit int) ([]*models.Bike, int, error) {
			return nil, 0, errors.New("database error")
		},
	}

	handler := &AdminHandler{adminService: mockService}
	req := httptest.NewRequest(http.MethodGet, "/admin/bikes", nil)
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.GetAllBikes(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAdminHandler_GetAllUsers_ServiceError(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	mockService := &MockAdminService2{
		GetAllUsersFunc: func(page, limit int) ([]*models.User, int, error) {
			return nil, 0, errors.New("database error")
		},
	}

	handler := &AdminHandler{adminService: mockService}
	req := httptest.NewRequest(http.MethodGet, "/admin/users", nil)
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.GetAllUsers(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAdminHandler_GetAllRentals_ServiceError(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	mockService := &MockAdminService2{
		GetAllRentalsFunc: func(page, limit int) ([]*models.Rental, int, error) {
			return nil, 0, errors.New("database error")
		},
	}

	handler := &AdminHandler{adminService: mockService}
	req := httptest.NewRequest(http.MethodGet, "/admin/rentals", nil)
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.GetAllRentals(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAdminHandler_UpdateBike_InvalidJSON2(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	handler := &AdminHandler{}
	req := httptest.NewRequest(http.MethodPut, "/admin/bikes/1", bytes.NewReader([]byte("invalid-json")))
	req.SetPathValue("bike-id", "1")
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.UpdateBike(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAdminHandler_UpdateRental_ServiceError(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	mockService := &MockAdminService2{
		UpdateRentalFunc: func(rentalID int, status *string) (*models.Rental, error) {
			return nil, errors.New("database error")
		},
	}

	handler := &AdminHandler{adminService: mockService}
	body, _ := json.Marshal(map[string]interface{}{"status": "ended"})
	req := httptest.NewRequest(http.MethodPut, "/admin/rentals/1", bytes.NewReader(body))
	req.SetPathValue("rental-id", "1")
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.UpdateRental(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAdminHandler_GetAllBikes_InvalidPage(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	mockService := &MockAdminService2{
		GetAllBikesFunc: func(page, limit int) ([]*models.Bike, int, error) {
			return []*models.Bike{}, 0, nil
		},
	}

	handler := &AdminHandler{adminService: mockService}
	req := httptest.NewRequest(http.MethodGet, "/admin/bikes?page=0", nil)
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.GetAllBikes(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAdminHandler_GetAllUsers_WithPagination(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	mockService := &MockAdminService2{
		GetAllUsersFunc: func(page, limit int) ([]*models.User, int, error) {
			return []*models.User{
				{ID: 1, Email: "user1@test.com"},
			}, 1, nil
		},
	}

	handler := &AdminHandler{adminService: mockService}
	req := httptest.NewRequest(http.MethodGet, "/admin/users?page=2&limit=5", nil)
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.GetAllUsers(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAdminHandler_GetAllRentals_WithPagination(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	mockService := &MockAdminService2{
		GetAllRentalsFunc: func(page, limit int) ([]*models.Rental, int, error) {
			return []*models.Rental{
				{ID: 1, UserID: 1, BikeID: 1},
			}, 1, nil
		},
	}

	handler := &AdminHandler{adminService: mockService}
	req := httptest.NewRequest(http.MethodGet, "/admin/rentals?page=1&limit=10", nil)
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.GetAllRentals(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAdminHandler_UpdateUser_WithPasswordHash(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	hashedPass := "hashedPassword123"
	mockService := &MockAdminService2{
		UpdateUserFunc: func(userID int, email, firstName, lastName, hashedPassword *string) (*models.User, error) {
			return &models.User{ID: userID, Email: "updated@test.com"}, nil
		},
	}

	handler := &AdminHandler{adminService: mockService}
	body, _ := json.Marshal(map[string]interface{}{
		"email":    "updated@test.com",
		"password": hashedPass,
	})
	req := httptest.NewRequest(http.MethodPut, "/admin/users/1", bytes.NewReader(body))
	req.SetPathValue("user-id", "1")
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.UpdateUser(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAdminHandler_AddBike_MinimumPrice(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	mockService := &MockAdminService2{
		CreateBikeFunc: func(latitude, longitude, pricePerMinute float64) (*models.Bike, error) {
			return &models.Bike{ID: 1, Latitude: latitude, Longitude: longitude, IsAvailable: true, PricePerMinute: pricePerMinute}, nil
		},
	}

	handler := &AdminHandler{adminService: mockService}
	body, _ := json.Marshal(map[string]interface{}{"latitude": 40.416775, "longitude": -3.703790, "pricePerMinute": 0.01})
	req := httptest.NewRequest(http.MethodPost, "/admin/bikes", bytes.NewReader(body))
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.AddBike(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAdminHandler_UpdateBike_AllFields(t *testing.T) {
	os.Setenv("ADMIN_CREDENTIALS", "YWRtaW46YWRtaW4xMjM=")
	defer os.Unsetenv("ADMIN_CREDENTIALS")

	latitude := 41.0
	longitude := -4.0
	price := 1.0
	available := false

	mockService := &MockAdminService2{
		UpdateBikeFunc: func(bikeID int, lat, long *float64, isAvailable *bool, pricePerMinute *float64) (*models.Bike, error) {
			bike := &models.Bike{ID: bikeID}
			if lat != nil {
				bike.Latitude = *lat
			}
			if long != nil {
				bike.Longitude = *long
			}
			if isAvailable != nil {
				bike.IsAvailable = *isAvailable
			}
			if pricePerMinute != nil {
				bike.PricePerMinute = *pricePerMinute
			}
			return bike, nil
		},
	}

	handler := &AdminHandler{adminService: mockService}
	body, _ := json.Marshal(map[string]interface{}{
		"latitude":       latitude,
		"longitude":      longitude,
		"pricePerMinute": price,
		"isAvailable":    available,
	})
	req := httptest.NewRequest(http.MethodPut, "/admin/bikes/1", bytes.NewReader(body))
	req.SetPathValue("bike-id", "1")
	req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4xMjM=")
	w := httptest.NewRecorder()

	handler.UpdateBike(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
