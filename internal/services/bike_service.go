package services

import (
	"github.com/Nimirandad/bike-rental-service/internal/models"
	"github.com/Nimirandad/bike-rental-service/internal/repositories"
)

type BikeRepository interface {
	CountAvailable() (int, error)
	GetAvailable(page, limit int) ([]*models.Bike, error)
	GetByID(bikeID int) (*models.Bike, error)
	UpdateAvailability(bikeID int, isAvailable bool) error
}

type BikeService struct {
	bikeRepo BikeRepository
}

func NewBikeService(bikeRepo *repositories.BikeRepository) *BikeService {
	return &BikeService{bikeRepo: bikeRepo}
}

func (s *BikeService) GetAvailableBikes(page, limit int) ([]*models.Bike, int, error) {
	total, err := s.bikeRepo.CountAvailable()
	if err != nil {
		return nil, 0, err
	}

	bikes, err := s.bikeRepo.GetAvailable(page, limit)
	if err != nil {
		return nil, 0, err
	}

	return bikes, total, nil
}