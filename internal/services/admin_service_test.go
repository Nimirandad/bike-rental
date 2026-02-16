package services

import (
	"errors"
	"testing"

	"github.com/Nimirandad/bike-rental-service/internal/constants"
	"github.com/Nimirandad/bike-rental-service/internal/models"
	"github.com/stretchr/testify/assert"
)

type MockAdminRepository struct {
	CreateBikeFunc             func(latitude, longitude, pricePerMinute float64) (*models.Bike, error)
	GetAllBikesFunc            func(page, limit int) ([]*models.Bike, error)
	CountAllFunc               func() (int, error)
	UpdateBikeFunc             func(bikeID int, latitude, longitude *float64, isAvailable *bool, pricePerMinute *float64) (*models.Bike, error)
	GetAllUsersFunc            func(page, limit int) ([]*models.User, error)
	CountAllUsersFunc          func() (int, error)
	GetUserByIDFunc            func(userID int) (*models.User, error)
	EmailExistsByOtherUserFunc func(email string, userID int) (bool, error)
	UpdateUserFunc             func(userID int, email, firstName, lastName, hashedPassword *string) (*models.User, error)
	GetAllRentalsFunc          func(page, limit int) ([]*models.Rental, error)
	CountAllRentalsFunc        func() (int, error)
	GetRentalByIDFunc          func(rentalID int) (*models.Rental, error)
	UpdateRentalFunc           func(rentalID int, status *string) (*models.Rental, error)
}

func (m *MockAdminRepository) CreateBike(latitude, longitude, pricePerMinute float64) (*models.Bike, error) {
	return m.CreateBikeFunc(latitude, longitude, pricePerMinute)
}

func (m *MockAdminRepository) GetAllBikes(page, limit int) ([]*models.Bike, error) {
	return m.GetAllBikesFunc(page, limit)
}

func (m *MockAdminRepository) CountAll() (int, error) {
	return m.CountAllFunc()
}

func (m *MockAdminRepository) UpdateBike(bikeID int, latitude, longitude *float64, isAvailable *bool, pricePerMinute *float64) (*models.Bike, error) {
	return m.UpdateBikeFunc(bikeID, latitude, longitude, isAvailable, pricePerMinute)
}

func (m *MockAdminRepository) GetAllUsers(page, limit int) ([]*models.User, error) {
	return m.GetAllUsersFunc(page, limit)
}

func (m *MockAdminRepository) CountAllUsers() (int, error) {
	return m.CountAllUsersFunc()
}

func (m *MockAdminRepository) GetUserByID(userID int) (*models.User, error) {
	return m.GetUserByIDFunc(userID)
}

func (m *MockAdminRepository) EmailExistsByOtherUser(email string, userID int) (bool, error) {
	return m.EmailExistsByOtherUserFunc(email, userID)
}

func (m *MockAdminRepository) UpdateUser(userID int, email, firstName, lastName, hashedPassword *string) (*models.User, error) {
	return m.UpdateUserFunc(userID, email, firstName, lastName, hashedPassword)
}

func (m *MockAdminRepository) GetAllRentals(page, limit int) ([]*models.Rental, error) {
	return m.GetAllRentalsFunc(page, limit)
}

func (m *MockAdminRepository) CountAllRentals() (int, error) {
	return m.CountAllRentalsFunc()
}

func (m *MockAdminRepository) GetRentalByID(rentalID int) (*models.Rental, error) {
	return m.GetRentalByIDFunc(rentalID)
}

func (m *MockAdminRepository) UpdateRental(rentalID int, status *string) (*models.Rental, error) {
	return m.UpdateRentalFunc(rentalID, status)
}

// TestAdminService_CreateBike_Success tests successful bike creation
func TestAdminService_CreateBike_Success(t *testing.T) {
	mockRepo := &MockAdminRepository{
		CreateBikeFunc: func(latitude, longitude, pricePerMinute float64) (*models.Bike, error) {
			return &models.Bike{
				ID:             1,
				Latitude:       latitude,
				Longitude:      longitude,
				IsAvailable:    true,
				PricePerMinute: pricePerMinute,
			}, nil
		},
	}

	service := &AdminService{adminRepo: mockRepo}
	bike, err := service.CreateBike(40.416775, -3.703790, 0.5)

	assert.NoError(t, err)
	assert.NotNil(t, bike)
	assert.Equal(t, 1, bike.ID)
	assert.Equal(t, 40.416775, bike.Latitude)
	assert.Equal(t, -3.703790, bike.Longitude)
	assert.Equal(t, 0.5, bike.PricePerMinute)
	assert.True(t, bike.IsAvailable)
}

// TestAdminService_CreateBike_Error tests error when creating bike
func TestAdminService_CreateBike_Error(t *testing.T) {
	mockRepo := &MockAdminRepository{
		CreateBikeFunc: func(latitude, longitude, pricePerMinute float64) (*models.Bike, error) {
			return nil, errors.New("database error")
		},
	}

	service := &AdminService{adminRepo: mockRepo}
	bike, err := service.CreateBike(40.416775, -3.703790, 0.5)

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
	assert.Nil(t, bike)
}

// TestAdminService_GetAllBikes_Success tests successful retrieval of all bikes
func TestAdminService_GetAllBikes_Success(t *testing.T) {
	mockRepo := &MockAdminRepository{
		CountAllFunc: func() (int, error) {
			return 20, nil
		},
		GetAllBikesFunc: func(page, limit int) ([]*models.Bike, error) {
			bikes := []*models.Bike{
				{ID: 1, Latitude: 40.416775, Longitude: -3.703790, IsAvailable: true, PricePerMinute: 0.5},
				{ID: 2, Latitude: 40.417832, Longitude: -3.705064, IsAvailable: false, PricePerMinute: 0.5},
			}
			return bikes, nil
		},
	}

	service := &AdminService{adminRepo: mockRepo}
	bikes, total, err := service.GetAllBikes(1, 10)

	assert.NoError(t, err)
	assert.Equal(t, 20, total)
	assert.Len(t, bikes, 2)
	assert.Equal(t, 1, bikes[0].ID)
	assert.Equal(t, 2, bikes[1].ID)
}

// TestAdminService_GetAllBikes_CountError tests error when counting bikes
func TestAdminService_GetAllBikes_CountError(t *testing.T) {
	mockRepo := &MockAdminRepository{
		CountAllFunc: func() (int, error) {
			return 0, errors.New("count error")
		},
	}

	service := &AdminService{adminRepo: mockRepo}
	bikes, total, err := service.GetAllBikes(1, 10)

	assert.Error(t, err)
	assert.Equal(t, "count error", err.Error())
	assert.Equal(t, 0, total)
	assert.Nil(t, bikes)
}

// TestAdminService_GetAllBikes_GetBikesError tests error when getting bikes
func TestAdminService_GetAllBikes_GetBikesError(t *testing.T) {
	mockRepo := &MockAdminRepository{
		CountAllFunc: func() (int, error) {
			return 20, nil
		},
		GetAllBikesFunc: func(page, limit int) ([]*models.Bike, error) {
			return nil, errors.New("query error")
		},
	}

	service := &AdminService{adminRepo: mockRepo}
	bikes, total, err := service.GetAllBikes(1, 10)

	assert.Error(t, err)
	assert.Equal(t, "query error", err.Error())
	assert.Equal(t, 0, total)
	assert.Nil(t, bikes)
}

// TestAdminService_UpdateBike_Success tests successful bike update
func TestAdminService_UpdateBike_Success(t *testing.T) {
	newLat := 41.0
	newLong := -4.0
	newAvailable := false
	newPrice := 0.75

	mockRepo := &MockAdminRepository{
		UpdateBikeFunc: func(bikeID int, latitude, longitude *float64, isAvailable *bool, pricePerMinute *float64) (*models.Bike, error) {
			return &models.Bike{
				ID:             bikeID,
				Latitude:       *latitude,
				Longitude:      *longitude,
				IsAvailable:    *isAvailable,
				PricePerMinute: *pricePerMinute,
			}, nil
		},
	}

	service := &AdminService{adminRepo: mockRepo}
	bike, err := service.UpdateBike(1, &newLat, &newLong, &newAvailable, &newPrice)

	assert.NoError(t, err)
	assert.NotNil(t, bike)
	assert.Equal(t, 1, bike.ID)
	assert.Equal(t, 41.0, bike.Latitude)
	assert.Equal(t, -4.0, bike.Longitude)
	assert.False(t, bike.IsAvailable)
	assert.Equal(t, 0.75, bike.PricePerMinute)
}

// TestAdminService_UpdateBike_Error tests error when updating bike
func TestAdminService_UpdateBike_Error(t *testing.T) {
	newLat := 41.0

	mockRepo := &MockAdminRepository{
		UpdateBikeFunc: func(bikeID int, latitude, longitude *float64, isAvailable *bool, pricePerMinute *float64) (*models.Bike, error) {
			return nil, errors.New("update error")
		},
	}

	service := &AdminService{adminRepo: mockRepo}
	bike, err := service.UpdateBike(1, &newLat, nil, nil, nil)

	assert.Error(t, err)
	assert.Equal(t, "update error", err.Error())
	assert.Nil(t, bike)
}

// TestAdminService_GetAllUsers_Success tests successful retrieval of all users
func TestAdminService_GetAllUsers_Success(t *testing.T) {
	mockRepo := &MockAdminRepository{
		CountAllUsersFunc: func() (int, error) {
			return 50, nil
		},
		GetAllUsersFunc: func(page, limit int) ([]*models.User, error) {
			users := []*models.User{
				{ID: 1, FirstName: "John", LastName: "Doe", Email: "john@example.com"},
				{ID: 2, FirstName: "Jane", LastName: "Smith", Email: "jane@example.com"},
			}
			return users, nil
		},
	}

	service := &AdminService{adminRepo: mockRepo}
	users, total, err := service.GetAllUsers(1, 10)

	assert.NoError(t, err)
	assert.Equal(t, 50, total)
	assert.Len(t, users, 2)
	assert.Equal(t, "John", users[0].FirstName)
	assert.Equal(t, "Jane", users[1].FirstName)
}

// TestAdminService_GetAllUsers_CountError tests error when counting users
func TestAdminService_GetAllUsers_CountError(t *testing.T) {
	mockRepo := &MockAdminRepository{
		CountAllUsersFunc: func() (int, error) {
			return 0, errors.New("count error")
		},
	}

	service := &AdminService{adminRepo: mockRepo}
	users, total, err := service.GetAllUsers(1, 10)

	assert.Error(t, err)
	assert.Equal(t, "count error", err.Error())
	assert.Equal(t, 0, total)
	assert.Nil(t, users)
}

// TestAdminService_GetAllUsers_GetUsersError tests error when getting users
func TestAdminService_GetAllUsers_GetUsersError(t *testing.T) {
	mockRepo := &MockAdminRepository{
		CountAllUsersFunc: func() (int, error) {
			return 50, nil
		},
		GetAllUsersFunc: func(page, limit int) ([]*models.User, error) {
			return nil, errors.New("query error")
		},
	}

	service := &AdminService{adminRepo: mockRepo}
	users, total, err := service.GetAllUsers(1, 10)

	assert.Error(t, err)
	assert.Equal(t, "query error", err.Error())
	assert.Equal(t, 0, total)
	assert.Nil(t, users)
}

// TestAdminService_GetUserByID_Success tests successful user retrieval by ID
func TestAdminService_GetUserByID_Success(t *testing.T) {
	mockRepo := &MockAdminRepository{
		GetUserByIDFunc: func(userID int) (*models.User, error) {
			return &models.User{
				ID:        userID,
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@example.com",
			}, nil
		},
	}

	service := &AdminService{adminRepo: mockRepo}
	user, err := service.GetUserByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, 1, user.ID)
	assert.Equal(t, "John", user.FirstName)
}

// TestAdminService_GetUserByID_Error tests error when getting user by ID
func TestAdminService_GetUserByID_Error(t *testing.T) {
	mockRepo := &MockAdminRepository{
		GetUserByIDFunc: func(userID int) (*models.User, error) {
			return nil, errors.New("user not found")
		},
	}

	service := &AdminService{adminRepo: mockRepo}
	user, err := service.GetUserByID(1)

	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	assert.Nil(t, user)
}

// TestAdminService_UpdateUser_Success tests successful user update
func TestAdminService_UpdateUser_Success(t *testing.T) {
	newEmail := "newemail@example.com"
	newFirstName := "Johnny"

	mockRepo := &MockAdminRepository{
		EmailExistsByOtherUserFunc: func(email string, userID int) (bool, error) {
			return false, nil
		},
		UpdateUserFunc: func(userID int, email, firstName, lastName, hashedPassword *string) (*models.User, error) {
			return &models.User{
				ID:        userID,
				Email:     *email,
				FirstName: *firstName,
			}, nil
		},
	}

	service := &AdminService{adminRepo: mockRepo}
	user, err := service.UpdateUser(1, &newEmail, &newFirstName, nil, nil)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, 1, user.ID)
	assert.Equal(t, "newemail@example.com", user.Email)
	assert.Equal(t, "Johnny", user.FirstName)
}

// TestAdminService_UpdateUser_EmailAlreadyExists tests error when email exists for another user
func TestAdminService_UpdateUser_EmailAlreadyExists(t *testing.T) {
	newEmail := "existing@example.com"

	mockRepo := &MockAdminRepository{
		EmailExistsByOtherUserFunc: func(email string, userID int) (bool, error) {
			return true, nil
		},
	}

	service := &AdminService{adminRepo: mockRepo}
	user, err := service.UpdateUser(1, &newEmail, nil, nil, nil)

	assert.Error(t, err)
	assert.Equal(t, constants.ErrEmailAlreadyExists, err)
	assert.Nil(t, user)
}

// TestAdminService_UpdateUser_EmailCheckError tests error when checking email existence
func TestAdminService_UpdateUser_EmailCheckError(t *testing.T) {
	newEmail := "test@example.com"

	mockRepo := &MockAdminRepository{
		EmailExistsByOtherUserFunc: func(email string, userID int) (bool, error) {
			return false, errors.New("database error")
		},
	}

	service := &AdminService{adminRepo: mockRepo}
	user, err := service.UpdateUser(1, &newEmail, nil, nil, nil)

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
	assert.Nil(t, user)
}

// TestAdminService_UpdateUser_NoEmailChange tests update without email change
func TestAdminService_UpdateUser_NoEmailChange(t *testing.T) {
	newFirstName := "Johnny"

	mockRepo := &MockAdminRepository{
		UpdateUserFunc: func(userID int, email, firstName, lastName, hashedPassword *string) (*models.User, error) {
			return &models.User{
				ID:        userID,
				FirstName: *firstName,
				Email:     "original@example.com",
			}, nil
		},
	}

	service := &AdminService{adminRepo: mockRepo}
	user, err := service.UpdateUser(1, nil, &newFirstName, nil, nil)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "Johnny", user.FirstName)
}

// TestAdminService_UpdateUser_UpdateError tests error when updating user
func TestAdminService_UpdateUser_UpdateError(t *testing.T) {
	newEmail := "newemail@example.com"

	mockRepo := &MockAdminRepository{
		EmailExistsByOtherUserFunc: func(email string, userID int) (bool, error) {
			return false, nil
		},
		UpdateUserFunc: func(userID int, email, firstName, lastName, hashedPassword *string) (*models.User, error) {
			return nil, errors.New("update error")
		},
	}

	service := &AdminService{adminRepo: mockRepo}
	user, err := service.UpdateUser(1, &newEmail, nil, nil, nil)

	assert.Error(t, err)
	assert.Equal(t, "update error", err.Error())
	assert.Nil(t, user)
}

// TestAdminService_GetAllRentals_Success tests successful retrieval of all rentals
func TestAdminService_GetAllRentals_Success(t *testing.T) {
	mockRepo := &MockAdminRepository{
		CountAllRentalsFunc: func() (int, error) {
			return 100, nil
		},
		GetAllRentalsFunc: func(page, limit int) ([]*models.Rental, error) {
			rentals := []*models.Rental{
				{ID: 1, UserID: 1, BikeID: 1, Status: "running"},
				{ID: 2, UserID: 2, BikeID: 2, Status: "ended"},
			}
			return rentals, nil
		},
	}

	service := &AdminService{adminRepo: mockRepo}
	rentals, total, err := service.GetAllRentals(1, 10)

	assert.NoError(t, err)
	assert.Equal(t, 100, total)
	assert.Len(t, rentals, 2)
	assert.Equal(t, 1, rentals[0].ID)
	assert.Equal(t, "running", rentals[0].Status)
}

// TestAdminService_GetAllRentals_CountError tests error when counting rentals
func TestAdminService_GetAllRentals_CountError(t *testing.T) {
	mockRepo := &MockAdminRepository{
		CountAllRentalsFunc: func() (int, error) {
			return 0, errors.New("count error")
		},
	}

	service := &AdminService{adminRepo: mockRepo}
	rentals, total, err := service.GetAllRentals(1, 10)

	assert.Error(t, err)
	assert.Equal(t, "count error", err.Error())
	assert.Equal(t, 0, total)
	assert.Nil(t, rentals)
}

// TestAdminService_GetAllRentals_GetRentalsError tests error when getting rentals
func TestAdminService_GetAllRentals_GetRentalsError(t *testing.T) {
	mockRepo := &MockAdminRepository{
		CountAllRentalsFunc: func() (int, error) {
			return 100, nil
		},
		GetAllRentalsFunc: func(page, limit int) ([]*models.Rental, error) {
			return nil, errors.New("query error")
		},
	}

	service := &AdminService{adminRepo: mockRepo}
	rentals, total, err := service.GetAllRentals(1, 10)

	assert.Error(t, err)
	assert.Equal(t, "query error", err.Error())
	assert.Equal(t, 0, total)
	assert.Nil(t, rentals)
}

// TestAdminService_GetRentalByID_Success tests successful rental retrieval by ID
func TestAdminService_GetRentalByID_Success(t *testing.T) {
	mockRepo := &MockAdminRepository{
		GetRentalByIDFunc: func(rentalID int) (*models.Rental, error) {
			return &models.Rental{
				ID:     rentalID,
				UserID: 1,
				BikeID: 1,
				Status: "running",
			}, nil
		},
	}

	service := &AdminService{adminRepo: mockRepo}
	rental, err := service.GetRentalByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, rental)
	assert.Equal(t, 1, rental.ID)
	assert.Equal(t, "running", rental.Status)
}

// TestAdminService_GetRentalByID_Error tests error when getting rental by ID
func TestAdminService_GetRentalByID_Error(t *testing.T) {
	mockRepo := &MockAdminRepository{
		GetRentalByIDFunc: func(rentalID int) (*models.Rental, error) {
			return nil, errors.New("rental not found")
		},
	}

	service := &AdminService{adminRepo: mockRepo}
	rental, err := service.GetRentalByID(1)

	assert.Error(t, err)
	assert.Equal(t, "rental not found", err.Error())
	assert.Nil(t, rental)
}

// TestAdminService_UpdateRental_Success tests successful rental update
func TestAdminService_UpdateRental_Success(t *testing.T) {
	newStatus := "ended"

	mockRepo := &MockAdminRepository{
		UpdateRentalFunc: func(rentalID int, status *string) (*models.Rental, error) {
			return &models.Rental{
				ID:     rentalID,
				Status: *status,
			}, nil
		},
	}

	service := &AdminService{adminRepo: mockRepo}
	rental, err := service.UpdateRental(1, &newStatus)

	assert.NoError(t, err)
	assert.NotNil(t, rental)
	assert.Equal(t, 1, rental.ID)
	assert.Equal(t, "ended", rental.Status)
}

// TestAdminService_UpdateRental_Error tests error when updating rental
func TestAdminService_UpdateRental_Error(t *testing.T) {
	newStatus := "ended"

	mockRepo := &MockAdminRepository{
		UpdateRentalFunc: func(rentalID int, status *string) (*models.Rental, error) {
			return nil, errors.New("update error")
		},
	}

	service := &AdminService{adminRepo: mockRepo}
	rental, err := service.UpdateRental(1, &newStatus)

	assert.Error(t, err)
	assert.Equal(t, "update error", err.Error())
	assert.Nil(t, rental)
}