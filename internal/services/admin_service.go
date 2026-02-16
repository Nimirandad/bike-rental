package services

import (
	"github.com/Nimirandad/bike-rental-service/internal/constants"
	"github.com/Nimirandad/bike-rental-service/internal/models"
	"github.com/Nimirandad/bike-rental-service/internal/repositories"
)

type AdminRepository interface {
	CreateBike(latitude, longitude, pricePerMinute float64) (*models.Bike, error)
	GetAllBikes(page, limit int) ([]*models.Bike, error)
	CountAll() (int, error)
	UpdateBike(bikeID int, latitude, longitude *float64, isAvailable *bool, pricePerMinute *float64) (*models.Bike, error)
	GetAllUsers(page, limit int) ([]*models.User, error)
	CountAllUsers() (int, error)
	GetUserByID(userID int) (*models.User, error)
	EmailExistsByOtherUser(email string, userID int) (bool, error)
	UpdateUser(userID int, email, firstName, lastName, hashedPassword *string) (*models.User, error)
	GetAllRentals(page, limit int) ([]*models.Rental, error)
	CountAllRentals() (int, error)
	GetRentalByID(rentalID int) (*models.Rental, error)
	UpdateRental(rentalID int, status *string) (*models.Rental, error)
}

type AdminService struct {
	adminRepo AdminRepository
}

func NewAdminService(adminRepo *repositories.AdminRepository) *AdminService {
	return &AdminService{
		adminRepo: adminRepo,
	}
}

func (s *AdminService) CreateBike(latitude, longitude, pricePerMinute float64) (*models.Bike, error) {
	return s.adminRepo.CreateBike(latitude, longitude, pricePerMinute)
}

func (s *AdminService) GetAllBikes(page, limit int) ([]*models.Bike, int, error) {
	total, err := s.adminRepo.CountAll()
	if err != nil {
		return nil, 0, err
	}

	bikes, err := s.adminRepo.GetAllBikes(page, limit)
	if err != nil {
		return nil, 0, err
	}

	return bikes, total, nil
}

func (s *AdminService) UpdateBike(bikeID int, latitude, longitude *float64, isAvailable *bool, pricePerMinute *float64) (*models.Bike, error) {
	return s.adminRepo.UpdateBike(bikeID, latitude, longitude, isAvailable, pricePerMinute)
}

func (s *AdminService) GetAllUsers(page, limit int) ([]*models.User, int, error) {
	total, err := s.adminRepo.CountAllUsers()
	if err != nil {
		return nil, 0, err
	}

	users, err := s.adminRepo.GetAllUsers(page, limit)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (s *AdminService) GetUserByID(userID int) (*models.User, error) {
	return s.adminRepo.GetUserByID(userID)
}

func (s *AdminService) UpdateUser(userID int, email, firstName, lastName, hashedPassword *string) (*models.User, error) {
	if email != nil {
		exists, err := s.adminRepo.EmailExistsByOtherUser(*email, userID)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, constants.ErrEmailAlreadyExists
		}
	}

	return s.adminRepo.UpdateUser(userID, email, firstName, lastName, hashedPassword)
}

func (s *AdminService) GetAllRentals(page, limit int) ([]*models.Rental, int, error) {
	total, err := s.adminRepo.CountAllRentals()
	if err != nil {
		return nil, 0, err
	}

	rentals, err := s.adminRepo.GetAllRentals(page, limit)
	if err != nil {
		return nil, 0, err
	}

	return rentals, total, nil
}

func (s *AdminService) GetRentalByID(rentalID int) (*models.Rental, error) {
	return s.adminRepo.GetRentalByID(rentalID)
}

func (s *AdminService) UpdateRental(rentalID int, status *string) (*models.Rental, error) {
	return s.adminRepo.UpdateRental(rentalID, status)
}