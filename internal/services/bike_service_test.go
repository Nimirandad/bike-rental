package services

import (
	"errors"
	"testing"

	"github.com/Nimirandad/bike-rental-service/internal/models"
	"github.com/stretchr/testify/assert"
)

type MockBikeRepository struct {
	CountAvailableFunc     func() (int, error)
	GetAvailableFunc       func(page, limit int) ([]*models.Bike, error)
	GetByIDFunc            func(bikeID int) (*models.Bike, error)
	UpdateAvailabilityFunc func(bikeID int, isAvailable bool) error
}

func (m *MockBikeRepository) CountAvailable() (int, error) {
	return m.CountAvailableFunc()
}

func (m *MockBikeRepository) GetAvailable(page, limit int) ([]*models.Bike, error) {
	return m.GetAvailableFunc(page, limit)
}

func (m *MockBikeRepository) GetByID(bikeID int) (*models.Bike, error) {
	return m.GetByIDFunc(bikeID)
}

func (m *MockBikeRepository) UpdateAvailability(bikeID int, isAvailable bool) error {
	return m.UpdateAvailabilityFunc(bikeID, isAvailable)
}

func TestBikeService_GetAvailableBikes_Success(t *testing.T) {
	mockRepo := &MockBikeRepository{
		CountAvailableFunc: func() (int, error) {
			return 5, nil
		},
		GetAvailableFunc: func(page, limit int) ([]*models.Bike, error) {
			bikes := []*models.Bike{
				{ID: 1, Latitude: 40.416775, Longitude: -3.703790, IsAvailable: true, PricePerMinute: 0.5},
				{ID: 2, Latitude: 40.417832, Longitude: -3.705064, IsAvailable: true, PricePerMinute: 0.5},
			}
			return bikes, nil
		},
	}

	service := &BikeService{bikeRepo: mockRepo}
	bikes, total, err := service.GetAvailableBikes(1, 10)

	assert.NoError(t, err)
	assert.Equal(t, 5, total)
	assert.Len(t, bikes, 2)
	assert.Equal(t, 1, bikes[0].ID)
	assert.Equal(t, 2, bikes[1].ID)
}

func TestBikeService_GetAvailableBikes_CountError(t *testing.T) {
	mockRepo := &MockBikeRepository{
		CountAvailableFunc: func() (int, error) {
			return 0, errors.New("database error")
		},
	}

	service := &BikeService{bikeRepo: mockRepo}
	bikes, total, err := service.GetAvailableBikes(1, 10)

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
	assert.Equal(t, 0, total)
	assert.Nil(t, bikes)
}

func TestBikeService_GetAvailableBikes_GetAvailableError(t *testing.T) {
	mockRepo := &MockBikeRepository{
		CountAvailableFunc: func() (int, error) {
			return 5, nil
		},
		GetAvailableFunc: func(page, limit int) ([]*models.Bike, error) {
			return nil, errors.New("query error")
		},
	}

	service := &BikeService{bikeRepo: mockRepo}
	bikes, total, err := service.GetAvailableBikes(1, 10)

	assert.Error(t, err)
	assert.Equal(t, "query error", err.Error())
	assert.Equal(t, 0, total)
	assert.Nil(t, bikes)
}

func TestBikeService_GetAvailableBikes_EmptyResult(t *testing.T) {
	mockRepo := &MockBikeRepository{
		CountAvailableFunc: func() (int, error) {
			return 0, nil
		},
		GetAvailableFunc: func(page, limit int) ([]*models.Bike, error) {
			return []*models.Bike{}, nil
		},
	}

	service := &BikeService{bikeRepo: mockRepo}
	bikes, total, err := service.GetAvailableBikes(1, 10)

	assert.NoError(t, err)
	assert.Equal(t, 0, total)
	assert.Empty(t, bikes)
}