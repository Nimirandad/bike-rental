package services

import (
	"fmt"

	"github.com/Nimirandad/bike-rental-service/internal/constants"
	"github.com/Nimirandad/bike-rental-service/internal/models"
	"github.com/Nimirandad/bike-rental-service/internal/repositories"
	"github.com/Nimirandad/bike-rental-service/internal/utils"
)

type UserRepository interface {
	EmailExists(email string) (bool, error)
	Create(email, hashedPassword, firstName, lastName string) (*models.User, error)
	GetByID(userID int) (*models.User, error)
	GetPasswordHashByEmail(email string) (string, *models.User, error)
	EmailExistsByOtherUser(email string, userID int) (bool, error)
	Update(userID int, email, firstName, lastName *string) (*models.User, error)
}

type UserService struct {
	userRepo UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) RegisterUser(email, password, firstName, lastName string) (*models.User, error) {
	exists, err := s.userRepo.EmailExists(email)
	if err != nil {
		return nil, fmt.Errorf("error checking email: %w", err)
	}
	if exists {
		return nil, constants.ErrEmailAlreadyExists
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	user, err := s.userRepo.Create(email, hashedPassword, firstName, lastName)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	return user, nil
}

func (s *UserService) GetByID(userID int) (*models.User, error) {
	return s.userRepo.GetByID(userID)
}

func (s *UserService) Login(email, password string) (*models.User, error) {
	hashedPassword, user, err := s.userRepo.GetPasswordHashByEmail(email)
	if err != nil {
		return nil, constants.ErrInvalidCredentials
	}

	if !utils.VerifyPassword(password, hashedPassword) {
		return nil, constants.ErrInvalidCredentials
	}

	return user, nil
}

func (s *UserService) UpdateUser(userID int, email, firstName, lastName *string) (*models.User, error) {
	if email != nil && *email != "" {
		exists, err := s.userRepo.EmailExistsByOtherUser(*email, userID)
		if err != nil {
			return nil, fmt.Errorf("error checking email: %w", err)
		}
		if exists {
			return nil, constants.ErrEmailAlreadyExists
		}
	}

	user, err := s.userRepo.Update(userID, email, firstName, lastName)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	return user, nil
}