package services

import (
	"errors"
	"testing"
	"time"

	"github.com/Nimirandad/bike-rental-service/internal/constants"
	"github.com/Nimirandad/bike-rental-service/internal/models"
	"github.com/stretchr/testify/assert"
)

type MockRentalRepository struct {
	HasActiveRentalFunc        func(userID int) (bool, error)
	CreateFunc                 func(userID, bikeID int, startLat, startLong float64) (*models.Rental, error)
	GetByIDFunc                func(rentalID int) (*models.Rental, error)
	GetActiveRentalsByUserFunc func(userID, page, limit int) ([]*models.Rental, error)
	CountByUserFunc            func(userID int) (int, error)
	GetActiveRentalByUserFunc  func(userID int) (*models.Rental, error)
	EndRentalFunc              func(rentalID int, endLat, endLong float64, durationMinutes int, cost float64) (*models.Rental, error)
}

func (m *MockRentalRepository) HasActiveRental(userID int) (bool, error) {
	return m.HasActiveRentalFunc(userID)
}

func (m *MockRentalRepository) Create(userID, bikeID int, startLat, startLong float64) (*models.Rental, error) {
	return m.CreateFunc(userID, bikeID, startLat, startLong)
}

func (m *MockRentalRepository) GetByID(rentalID int) (*models.Rental, error) {
	return m.GetByIDFunc(rentalID)
}

func (m *MockRentalRepository) GetActiveRentalsByUser(userID, page, limit int) ([]*models.Rental, error) {
	return m.GetActiveRentalsByUserFunc(userID, page, limit)
}

func (m *MockRentalRepository) CountByUser(userID int) (int, error) {
	return m.CountByUserFunc(userID)
}

func (m *MockRentalRepository) GetActiveRentalByUser(userID int) (*models.Rental, error) {
	return m.GetActiveRentalByUserFunc(userID)
}

func (m *MockRentalRepository) EndRental(rentalID int, endLat, endLong float64, durationMinutes int, cost float64) (*models.Rental, error) {
	return m.EndRentalFunc(rentalID, endLat, endLong, durationMinutes, cost)
}

// TestRentalService_StartRental_Success tests successful rental start
func TestRentalService_StartRental_Success(t *testing.T) {
	mockRentalRepo := &MockRentalRepository{
		HasActiveRentalFunc: func(userID int) (bool, error) {
			return false, nil
		},
		CreateFunc: func(userID, bikeID int, startLat, startLong float64) (*models.Rental, error) {
			return &models.Rental{
				ID:             1,
				UserID:         userID,
				BikeID:         bikeID,
				StartLatitude:  startLat,
				StartLongitude: startLong,
				Status:         "running",
				StartTime:      time.Now(),
			}, nil
		},
	}

	mockBikeRepo := &MockBikeRepository{
		GetByIDFunc: func(bikeID int) (*models.Bike, error) {
			return &models.Bike{
				ID:             bikeID,
				Latitude:       40.416775,
				Longitude:      -3.703790,
				IsAvailable:    true,
				PricePerMinute: 0.5,
			}, nil
		},
		UpdateAvailabilityFunc: func(bikeID int, isAvailable bool) error {
			return nil
		},
	}

	service := &RentalService{rentalRepo: mockRentalRepo, bikeRepo: mockBikeRepo}
	rental, err := service.StartRental(1, 1)

	assert.NoError(t, err)
	assert.NotNil(t, rental)
	assert.Equal(t, 1, rental.ID)
	assert.Equal(t, 1, rental.UserID)
	assert.Equal(t, 1, rental.BikeID)
	assert.Equal(t, "running", rental.Status)
}

// TestRentalService_StartRental_UserHasActiveRental tests error when user already has active rental
func TestRentalService_StartRental_UserHasActiveRental(t *testing.T) {
	mockRentalRepo := &MockRentalRepository{
		HasActiveRentalFunc: func(userID int) (bool, error) {
			return true, nil
		},
	}

	service := &RentalService{rentalRepo: mockRentalRepo, bikeRepo: nil}
	rental, err := service.StartRental(1, 1)

	assert.Error(t, err)
	assert.Equal(t, constants.ErrUserHasActiveRental, err)
	assert.Nil(t, rental)
}

// TestRentalService_StartRental_HasActiveRentalCheckError tests error when checking for active rental
func TestRentalService_StartRental_HasActiveRentalCheckError(t *testing.T) {
	mockRentalRepo := &MockRentalRepository{
		HasActiveRentalFunc: func(userID int) (bool, error) {
			return false, errors.New("database error")
		},
	}

	service := &RentalService{rentalRepo: mockRentalRepo, bikeRepo: nil}
	rental, err := service.StartRental(1, 1)

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
	assert.Nil(t, rental)
}

// TestRentalService_StartRental_BikeNotFound tests error when bike doesn't exist
func TestRentalService_StartRental_BikeNotFound(t *testing.T) {
	mockRentalRepo := &MockRentalRepository{
		HasActiveRentalFunc: func(userID int) (bool, error) {
			return false, nil
		},
	}

	mockBikeRepo := &MockBikeRepository{
		GetByIDFunc: func(bikeID int) (*models.Bike, error) {
			return nil, errors.New("bike not found")
		},
	}

	service := &RentalService{rentalRepo: mockRentalRepo, bikeRepo: mockBikeRepo}
	rental, err := service.StartRental(1, 1)

	assert.Error(t, err)
	assert.Equal(t, constants.ErrBikeNotFound, err)
	assert.Nil(t, rental)
}

// TestRentalService_StartRental_BikeNotAvailable tests error when bike is not available
func TestRentalService_StartRental_BikeNotAvailable(t *testing.T) {
	mockRentalRepo := &MockRentalRepository{
		HasActiveRentalFunc: func(userID int) (bool, error) {
			return false, nil
		},
	}

	mockBikeRepo := &MockBikeRepository{
		GetByIDFunc: func(bikeID int) (*models.Bike, error) {
			return &models.Bike{
				ID:             bikeID,
				Latitude:       40.416775,
				Longitude:      -3.703790,
				IsAvailable:    false,
				PricePerMinute: 0.5,
			}, nil
		},
	}

	service := &RentalService{rentalRepo: mockRentalRepo, bikeRepo: mockBikeRepo}
	rental, err := service.StartRental(1, 1)

	assert.Error(t, err)
	assert.Equal(t, constants.ErrBikeNotAvailable, err)
	assert.Nil(t, rental)
}

// TestRentalService_StartRental_CreateRentalError tests error when creating rental
func TestRentalService_StartRental_CreateRentalError(t *testing.T) {
	mockRentalRepo := &MockRentalRepository{
		HasActiveRentalFunc: func(userID int) (bool, error) {
			return false, nil
		},
		CreateFunc: func(userID, bikeID int, startLat, startLong float64) (*models.Rental, error) {
			return nil, errors.New("create error")
		},
	}

	mockBikeRepo := &MockBikeRepository{
		GetByIDFunc: func(bikeID int) (*models.Bike, error) {
			return &models.Bike{
				ID:             bikeID,
				Latitude:       40.416775,
				Longitude:      -3.703790,
				IsAvailable:    true,
				PricePerMinute: 0.5,
			}, nil
		},
	}

	service := &RentalService{rentalRepo: mockRentalRepo, bikeRepo: mockBikeRepo}
	rental, err := service.StartRental(1, 1)

	assert.Error(t, err)
	assert.Equal(t, "create error", err.Error())
	assert.Nil(t, rental)
}

// TestRentalService_StartRental_UpdateAvailabilityError tests error when updating bike availability
func TestRentalService_StartRental_UpdateAvailabilityError(t *testing.T) {
	mockRentalRepo := &MockRentalRepository{
		HasActiveRentalFunc: func(userID int) (bool, error) {
			return false, nil
		},
		CreateFunc: func(userID, bikeID int, startLat, startLong float64) (*models.Rental, error) {
			return &models.Rental{
				ID:             1,
				UserID:         userID,
				BikeID:         bikeID,
				StartLatitude:  startLat,
				StartLongitude: startLong,
				Status:         "running",
				StartTime:      time.Now(),
			}, nil
		},
	}

	mockBikeRepo := &MockBikeRepository{
		GetByIDFunc: func(bikeID int) (*models.Bike, error) {
			return &models.Bike{
				ID:             bikeID,
				Latitude:       40.416775,
				Longitude:      -3.703790,
				IsAvailable:    true,
				PricePerMinute: 0.5,
			}, nil
		},
		UpdateAvailabilityFunc: func(bikeID int, isAvailable bool) error {
			return errors.New("update error")
		},
	}

	service := &RentalService{rentalRepo: mockRentalRepo, bikeRepo: mockBikeRepo}
	rental, err := service.StartRental(1, 1)

	assert.Error(t, err)
	assert.Equal(t, "update error", err.Error())
	assert.Nil(t, rental)
}

// TestRentalService_GetRentalHistory_Success tests successful rental history retrieval
func TestRentalService_GetRentalHistory_Success(t *testing.T) {
	mockRentalRepo := &MockRentalRepository{
		CountByUserFunc: func(userID int) (int, error) {
			return 10, nil
		},
		GetActiveRentalsByUserFunc: func(userID, page, limit int) ([]*models.Rental, error) {
			rentals := []*models.Rental{
				{ID: 1, UserID: userID, BikeID: 1, Status: "running"},
				{ID: 2, UserID: userID, BikeID: 2, Status: "ended"},
			}
			return rentals, nil
		},
	}

	service := &RentalService{rentalRepo: mockRentalRepo, bikeRepo: nil}
	rentals, total, err := service.GetRentalHistory(1, 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, 10, total)
	assert.Len(t, rentals, 2)
	assert.Equal(t, 1, rentals[0].ID)
	assert.Equal(t, 2, rentals[1].ID)
}

// TestRentalService_GetRentalHistory_CountError tests error when counting rentals
func TestRentalService_GetRentalHistory_CountError(t *testing.T) {
	mockRentalRepo := &MockRentalRepository{
		CountByUserFunc: func(userID int) (int, error) {
			return 0, errors.New("count error")
		},
	}

	service := &RentalService{rentalRepo: mockRentalRepo, bikeRepo: nil}
	rentals, total, err := service.GetRentalHistory(1, 1, 10)

	assert.Error(t, err)
	assert.Equal(t, "count error", err.Error())
	assert.Equal(t, 0, total)
	assert.Nil(t, rentals)
}

// TestRentalService_GetRentalHistory_GetRentalsError tests error when getting rentals
func TestRentalService_GetRentalHistory_GetRentalsError(t *testing.T) {
	mockRentalRepo := &MockRentalRepository{
		CountByUserFunc: func(userID int) (int, error) {
			return 10, nil
		},
		GetActiveRentalsByUserFunc: func(userID, page, limit int) ([]*models.Rental, error) {
			return nil, errors.New("query error")
		},
	}

	service := &RentalService{rentalRepo: mockRentalRepo, bikeRepo: nil}
	rentals, total, err := service.GetRentalHistory(1, 1, 10)

	assert.Error(t, err)
	assert.Equal(t, "query error", err.Error())
	assert.Equal(t, 0, total)
	assert.Nil(t, rentals)
}

// TestRentalService_EndRental_Success tests successful rental ending within 5km
func TestRentalService_EndRental_Success(t *testing.T) {
	startTime := time.Now().Add(-30 * time.Minute)

	mockRentalRepo := &MockRentalRepository{
		GetActiveRentalByUserFunc: func(userID int) (*models.Rental, error) {
			return &models.Rental{
				ID:             1,
				UserID:         userID,
				BikeID:         1,
				StartLatitude:  40.416775,
				StartLongitude: -3.703790,
				Status:         "running",
				StartTime:      startTime,
			}, nil
		},
		EndRentalFunc: func(rentalID int, endLat, endLong float64, durationMinutes int, cost float64) (*models.Rental, error) {
			return &models.Rental{
				ID:     rentalID,
				Status: "ended",
			}, nil
		},
	}

	mockBikeRepo := &MockBikeRepository{
		GetByIDFunc: func(bikeID int) (*models.Bike, error) {
			return &models.Bike{
				ID:             bikeID,
				Latitude:       40.416775,
				Longitude:      -3.703790,
				IsAvailable:    false,
				PricePerMinute: 0.5,
			}, nil
		},
		UpdateAvailabilityFunc: func(bikeID int, isAvailable bool) error {
			return nil
		},
	}

	service := &RentalService{rentalRepo: mockRentalRepo, bikeRepo: mockBikeRepo}
	// End location within 5km (approximately same location)
	rental, err := service.EndRental(1, 40.420000, -3.700000)

	assert.NoError(t, err)
	assert.NotNil(t, rental)
	assert.Equal(t, 1, rental.ID)
	assert.Equal(t, "ended", rental.Status)
}

// TestRentalService_EndRental_NoActiveRental tests error when user has no active rental
func TestRentalService_EndRental_NoActiveRental(t *testing.T) {
	mockRentalRepo := &MockRentalRepository{
		GetActiveRentalByUserFunc: func(userID int) (*models.Rental, error) {
			return nil, nil
		},
	}

	service := &RentalService{rentalRepo: mockRentalRepo, bikeRepo: nil}
	rental, err := service.EndRental(1, 40.420000, -3.700000)

	assert.Error(t, err)
	assert.Equal(t, constants.ErrNoActiveRental, err)
	assert.Nil(t, rental)
}

// TestRentalService_EndRental_GetActiveRentalError tests error when getting active rental
func TestRentalService_EndRental_GetActiveRentalError(t *testing.T) {
	mockRentalRepo := &MockRentalRepository{
		GetActiveRentalByUserFunc: func(userID int) (*models.Rental, error) {
			return nil, errors.New("database error")
		},
	}

	service := &RentalService{rentalRepo: mockRentalRepo, bikeRepo: nil}
	rental, err := service.EndRental(1, 40.420000, -3.700000)

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
	assert.Nil(t, rental)
}

// TestRentalService_EndRental_EndLocationTooFar tests error when end location is more than 5km away
func TestRentalService_EndRental_EndLocationTooFar(t *testing.T) {
	startTime := time.Now().Add(-30 * time.Minute)

	mockRentalRepo := &MockRentalRepository{
		GetActiveRentalByUserFunc: func(userID int) (*models.Rental, error) {
			return &models.Rental{
				ID:             1,
				UserID:         userID,
				BikeID:         1,
				StartLatitude:  40.416775,
				StartLongitude: -3.703790,
				Status:         "running",
				StartTime:      startTime,
			}, nil
		},
	}

	service := &RentalService{rentalRepo: mockRentalRepo, bikeRepo: nil}
	// End location more than 5km away (Paris coordinates - ~1050km from Madrid)
	rental, err := service.EndRental(1, 48.856614, 2.352222)

	assert.Error(t, err)
	assert.Equal(t, constants.ErrEndLocationTooFar, err)
	assert.Nil(t, rental)
}

// TestRentalService_EndRental_GetBikeError tests error when getting bike
func TestRentalService_EndRental_GetBikeError(t *testing.T) {
	startTime := time.Now().Add(-30 * time.Minute)

	mockRentalRepo := &MockRentalRepository{
		GetActiveRentalByUserFunc: func(userID int) (*models.Rental, error) {
			return &models.Rental{
				ID:             1,
				UserID:         userID,
				BikeID:         1,
				StartLatitude:  40.416775,
				StartLongitude: -3.703790,
				Status:         "running",
				StartTime:      startTime,
			}, nil
		},
	}

	mockBikeRepo := &MockBikeRepository{
		GetByIDFunc: func(bikeID int) (*models.Bike, error) {
			return nil, errors.New("bike not found")
		},
	}

	service := &RentalService{rentalRepo: mockRentalRepo, bikeRepo: mockBikeRepo}
	rental, err := service.EndRental(1, 40.420000, -3.700000)

	assert.Error(t, err)
	assert.Equal(t, "bike not found", err.Error())
	assert.Nil(t, rental)
}

// TestRentalService_EndRental_EndRentalError tests error when ending rental in repository
func TestRentalService_EndRental_EndRentalError(t *testing.T) {
	startTime := time.Now().Add(-30 * time.Minute)

	mockRentalRepo := &MockRentalRepository{
		GetActiveRentalByUserFunc: func(userID int) (*models.Rental, error) {
			return &models.Rental{
				ID:             1,
				UserID:         userID,
				BikeID:         1,
				StartLatitude:  40.416775,
				StartLongitude: -3.703790,
				Status:         "running",
				StartTime:      startTime,
			}, nil
		},
		EndRentalFunc: func(rentalID int, endLat, endLong float64, durationMinutes int, cost float64) (*models.Rental, error) {
			return nil, errors.New("update error")
		},
	}

	mockBikeRepo := &MockBikeRepository{
		GetByIDFunc: func(bikeID int) (*models.Bike, error) {
			return &models.Bike{
				ID:             bikeID,
				Latitude:       40.416775,
				Longitude:      -3.703790,
				IsAvailable:    false,
				PricePerMinute: 0.5,
			}, nil
		},
	}

	service := &RentalService{rentalRepo: mockRentalRepo, bikeRepo: mockBikeRepo}
	rental, err := service.EndRental(1, 40.420000, -3.700000)

	assert.Error(t, err)
	assert.Equal(t, "update error", err.Error())
	assert.Nil(t, rental)
}

// TestRentalService_EndRental_UpdateAvailabilityError tests error when updating bike availability after ending rental
func TestRentalService_EndRental_UpdateAvailabilityError(t *testing.T) {
	startTime := time.Now().Add(-30 * time.Minute)

	mockRentalRepo := &MockRentalRepository{
		GetActiveRentalByUserFunc: func(userID int) (*models.Rental, error) {
			return &models.Rental{
				ID:             1,
				UserID:         userID,
				BikeID:         1,
				StartLatitude:  40.416775,
				StartLongitude: -3.703790,
				Status:         "running",
				StartTime:      startTime,
			}, nil
		},
		EndRentalFunc: func(rentalID int, endLat, endLong float64, durationMinutes int, cost float64) (*models.Rental, error) {
			return &models.Rental{
				ID:     rentalID,
				Status: "ended",
			}, nil
		},
	}

	mockBikeRepo := &MockBikeRepository{
		GetByIDFunc: func(bikeID int) (*models.Bike, error) {
			return &models.Bike{
				ID:             bikeID,
				Latitude:       40.416775,
				Longitude:      -3.703790,
				IsAvailable:    false,
				PricePerMinute: 0.5,
			}, nil
		},
		UpdateAvailabilityFunc: func(bikeID int, isAvailable bool) error {
			return errors.New("availability update error")
		},
	}

	service := &RentalService{rentalRepo: mockRentalRepo, bikeRepo: mockBikeRepo}
	rental, err := service.EndRental(1, 40.420000, -3.700000)

	assert.Error(t, err)
	assert.Equal(t, "availability update error", err.Error())
	assert.Nil(t, rental)
}