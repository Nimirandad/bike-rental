package services

import (
	"math"
	"time"

	"github.com/Nimirandad/bike-rental-service/internal/constants"
	"github.com/Nimirandad/bike-rental-service/internal/models"
	"github.com/Nimirandad/bike-rental-service/internal/repositories"
	"github.com/Nimirandad/bike-rental-service/internal/utils"
)

type RentalRepository interface {
	HasActiveRental(userID int) (bool, error)
	Create(userID, bikeID int, startLat, startLong float64) (*models.Rental, error)
	GetByID(rentalID int) (*models.Rental, error)
	GetActiveRentalsByUser(userID, page, limit int) ([]*models.Rental, error)
	CountByUser(userID int) (int, error)
	GetActiveRentalByUser(userID int) (*models.Rental, error)
	EndRental(rentalID int, endLat, endLong float64, durationMinutes int, cost float64) (*models.Rental, error)
}

type RentalService struct {
	rentalRepo RentalRepository
	bikeRepo   BikeRepository
}

func NewRentalService(rentalRepo *repositories.RentalRepository, bikeRepo *repositories.BikeRepository) *RentalService {
	return &RentalService{
		rentalRepo: rentalRepo,
		bikeRepo:   bikeRepo,
	}
}

func (s *RentalService) StartRental(userID, bikeID int) (*models.Rental, error) {
	hasActive, err := s.rentalRepo.HasActiveRental(userID)
	if err != nil {
		return nil, err
	}
	if hasActive {
		return nil, constants.ErrUserHasActiveRental
	}

	bike, err := s.bikeRepo.GetByID(bikeID)
	if err != nil {
		return nil, constants.ErrBikeNotFound
	}

	if !bike.IsAvailable {
		return nil, constants.ErrBikeNotAvailable
	}

	rental, err := s.rentalRepo.Create(userID, bikeID, bike.Latitude, bike.Longitude)
	if err != nil {
		return nil, err
	}

	err = s.bikeRepo.UpdateAvailability(bikeID, false)
	if err != nil {
		return nil, err
	}

	return rental, nil
}

func (s *RentalService) GetRentalHistory(userID int, page, limit int) ([]*models.Rental, int, error) {
	total, err := s.rentalRepo.CountByUser(userID)
	if err != nil {
		return nil, 0, err
	}

	rentals, err := s.rentalRepo.GetActiveRentalsByUser(userID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return rentals, total, nil
}

func (s *RentalService) EndRental(userID int, endLat, endLong float64) (*models.Rental, error) {
	activeRental, err := s.rentalRepo.GetActiveRentalByUser(userID)
	if err != nil {
		return nil, err
	}

	if activeRental == nil {
		return nil, constants.ErrNoActiveRental
	}

	distance := utils.HaversineDistance(
		activeRental.StartLatitude,
		activeRental.StartLongitude,
		endLat,
		endLong,
	)

	if distance > 5.0 {
		return nil, constants.ErrEndLocationTooFar
	}

	bike, err := s.bikeRepo.GetByID(activeRental.BikeID)
	if err != nil {
		return nil, err
	}

	duration := time.Since(activeRental.StartTime)
	durationMinutes := int(math.Ceil(duration.Minutes()))

	cost := float64(durationMinutes) * bike.PricePerMinute

	rental, err := s.rentalRepo.EndRental(activeRental.ID, endLat, endLong, durationMinutes, cost)
	if err != nil {
		return nil, err
	}

	err = s.bikeRepo.UpdateAvailability(activeRental.BikeID, true)
	if err != nil {
		return nil, err
	}

	return rental, nil
}