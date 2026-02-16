package services

import (
	"errors"
	"testing"
	"time"

	"github.com/Nimirandad/bike-rental-service/internal/constants"
	"github.com/Nimirandad/bike-rental-service/internal/models"
	"github.com/Nimirandad/bike-rental-service/internal/utils"
	"github.com/stretchr/testify/assert"
)

type MockUserRepository struct {
	EmailExistsFunc            func(email string) (bool, error)
	CreateFunc                 func(email, hashedPassword, firstName, lastName string) (*models.User, error)
	GetByIDFunc                func(userID int) (*models.User, error)
	GetPasswordHashByEmailFunc func(email string) (string, *models.User, error)
	EmailExistsByOtherUserFunc func(email string, userID int) (bool, error)
	UpdateFunc                 func(userID int, email, firstName, lastName *string) (*models.User, error)
}

func (m *MockUserRepository) EmailExists(email string) (bool, error) {
	return m.EmailExistsFunc(email)
}

func (m *MockUserRepository) Create(email, hashedPassword, firstName, lastName string) (*models.User, error) {
	return m.CreateFunc(email, hashedPassword, firstName, lastName)
}

func (m *MockUserRepository) GetByID(userID int) (*models.User, error) {
	return m.GetByIDFunc(userID)
}

func (m *MockUserRepository) GetPasswordHashByEmail(email string) (string, *models.User, error) {
	return m.GetPasswordHashByEmailFunc(email)
}

func (m *MockUserRepository) EmailExistsByOtherUser(email string, userID int) (bool, error) {
	return m.EmailExistsByOtherUserFunc(email, userID)
}

func (m *MockUserRepository) Update(userID int, email, firstName, lastName *string) (*models.User, error) {
	return m.UpdateFunc(userID, email, firstName, lastName)
}

func TestUserService_RegisterUser(t *testing.T) {
	t.Run("Successfully register user", func(t *testing.T) {
		mockRepo := &MockUserRepository{
			EmailExistsFunc: func(email string) (bool, error) {
				return false, nil
			},
			CreateFunc: func(email, hashedPassword, firstName, lastName string) (*models.User, error) {
				return &models.User{
					ID:        1,
					Email:     email,
					FirstName: firstName,
					LastName:  lastName,
					CreatedAt: time.Now(),
				}, nil
			},
		}

		service := &UserService{userRepo: mockRepo}

		user, err := service.RegisterUser("test@example.com", "Password123", "John", "Doe")

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "test@example.com", user.Email)
		assert.Equal(t, "John", user.FirstName)
	})

	t.Run("Email already exists", func(t *testing.T) {
		mockRepo := &MockUserRepository{
			EmailExistsFunc: func(email string) (bool, error) {
				return true, nil
			},
		}

		service := &UserService{userRepo: mockRepo}

		user, err := service.RegisterUser("test@example.com", "Password123", "John", "Doe")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, constants.ErrEmailAlreadyExists, err)
	})

	t.Run("Error checking email", func(t *testing.T) {
		mockRepo := &MockUserRepository{
			EmailExistsFunc: func(email string) (bool, error) {
				return false, errors.New("database error")
			},
		}

		service := &UserService{userRepo: mockRepo}

		user, err := service.RegisterUser("test@example.com", "Password123", "John", "Doe")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "error checking email")
	})
}

func TestUserService_GetByID(t *testing.T) {
	t.Run("Successfully get user", func(t *testing.T) {
		mockRepo := &MockUserRepository{
			GetByIDFunc: func(userID int) (*models.User, error) {
				return &models.User{
					ID:        userID,
					Email:     "test@example.com",
					FirstName: "John",
					LastName:  "Doe",
					CreatedAt: time.Now(),
				}, nil
			},
		}

		service := &UserService{userRepo: mockRepo}

		user, err := service.GetByID(1)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, 1, user.ID)
	})

	t.Run("User not found", func(t *testing.T) {
		mockRepo := &MockUserRepository{
			GetByIDFunc: func(userID int) (*models.User, error) {
				return nil, errors.New("user not found")
			},
		}

		service := &UserService{userRepo: mockRepo}

		user, err := service.GetByID(999)

		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestUserService_Login(t *testing.T) {
	t.Run("Successful login", func(t *testing.T) {
		hashedPassword, _ := utils.HashPassword("secret")

		mockRepo := &MockUserRepository{
			GetPasswordHashByEmailFunc: func(email string) (string, *models.User, error) {
				return hashedPassword, &models.User{
					ID:        1,
					Email:     email,
					FirstName: "John",
					LastName:  "Doe",
					CreatedAt: time.Now(),
				}, nil
			},
		}

		service := &UserService{userRepo: mockRepo}

		user, err := service.Login("test@example.com", "secret")

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "test@example.com", user.Email)
	})

	t.Run("User not found", func(t *testing.T) {
		mockRepo := &MockUserRepository{
			GetPasswordHashByEmailFunc: func(email string) (string, *models.User, error) {
				return "", nil, errors.New("user not found")
			},
		}

		service := &UserService{userRepo: mockRepo}

		user, err := service.Login("notfound@example.com", "password")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, constants.ErrInvalidCredentials, err)
	})

	t.Run("Wrong password", func(t *testing.T) {
		hashedPassword, _ := utils.HashPassword("secret")

		mockRepo := &MockUserRepository{
			GetPasswordHashByEmailFunc: func(email string) (string, *models.User, error) {
				return hashedPassword, &models.User{
					ID:        1,
					Email:     email,
					FirstName: "John",
					LastName:  "Doe",
					CreatedAt: time.Now(),
				}, nil
			},
		}

		service := &UserService{userRepo: mockRepo}

		user, err := service.Login("test@example.com", "wrongpassword")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, constants.ErrInvalidCredentials, err)
	})
}

func TestUserService_UpdateUser(t *testing.T) {
	newEmail := "newemail@example.com"
	newFirstName := "Jane"

	t.Run("Successfully update user", func(t *testing.T) {
		mockRepo := &MockUserRepository{
			EmailExistsByOtherUserFunc: func(email string, userID int) (bool, error) {
				return false, nil
			},
			UpdateFunc: func(userID int, email, firstName, lastName *string) (*models.User, error) {
				return &models.User{
					ID:        userID,
					Email:     *email,
					FirstName: *firstName,
					LastName:  "Doe",
					CreatedAt: time.Now(),
				}, nil
			},
		}

		service := &UserService{userRepo: mockRepo}

		user, err := service.UpdateUser(1, &newEmail, &newFirstName, nil)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, newEmail, user.Email)
		assert.Equal(t, newFirstName, user.FirstName)
	})

	t.Run("Email already exists for other user", func(t *testing.T) {
		mockRepo := &MockUserRepository{
			EmailExistsByOtherUserFunc: func(email string, userID int) (bool, error) {
				return true, nil
			},
		}

		service := &UserService{userRepo: mockRepo}

		user, err := service.UpdateUser(1, &newEmail, nil, nil)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, constants.ErrEmailAlreadyExists, err)
	})

	t.Run("Update with empty string fields skips checking", func(t *testing.T) {
		emptyEmail := ""
		mockRepo := &MockUserRepository{
			UpdateFunc: func(userID int, email, firstName, lastName *string) (*models.User, error) {
				return &models.User{
					ID:        userID,
					Email:     "original@example.com",
					FirstName: *firstName,
					LastName:  "Doe",
					CreatedAt: time.Now(),
				}, nil
			},
		}

		service := &UserService{userRepo: mockRepo}

		user, err := service.UpdateUser(1, &emptyEmail, &newFirstName, nil)

		assert.NoError(t, err)
		assert.NotNil(t, user)
	})
}