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

type MockUserService2 struct {
	RegisterUserFunc func(email, password, firstName, lastName string) (*models.User, error)
	LoginFunc        func(email, password string) (*models.User, error)
	GetByIDFunc      func(userID int) (*models.User, error)
	UpdateUserFunc   func(userID int, email, firstName, lastName *string) (*models.User, error)
}

func (m *MockUserService2) RegisterUser(email, password, firstName, lastName string) (*models.User, error) {
	return m.RegisterUserFunc(email, password, firstName, lastName)
}

func (m *MockUserService2) Login(email, password string) (*models.User, error) {
	return m.LoginFunc(email, password)
}

func (m *MockUserService2) GetByID(userID int) (*models.User, error) {
	return m.GetByIDFunc(userID)
}

func (m *MockUserService2) UpdateUser(userID int, email, firstName, lastName *string) (*models.User, error) {
	return m.UpdateUserFunc(userID, email, firstName, lastName)
}

func TestUserHandler_RegisterUser_Success(t *testing.T) {
	mockService := &MockUserService2{
		RegisterUserFunc: func(email, password, firstName, lastName string) (*models.User, error) {
			return &models.User{ID: 1, Email: email, FirstName: firstName, LastName: lastName}, nil
		},
	}

	handler := &UserHandler{userService: mockService}
	body, _ := json.Marshal(map[string]string{"email": "test@example.com", "password": "password123", "first_name": "John", "last_name": "Doe"})
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.RegisterUser(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserHandler_RegisterUser_InvalidJSON(t *testing.T) {
	handler := &UserHandler{}
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewReader([]byte("invalid")))
	w := httptest.NewRecorder()

	handler.RegisterUser(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_RegisterUser_ValidationError(t *testing.T) {
	handler := &UserHandler{}
	body, _ := json.Marshal(map[string]string{"email": "invalid", "password": "123", "first_name": "J", "last_name": "D"})
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.RegisterUser(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_RegisterUser_EmailExists(t *testing.T) {
	mockService := &MockUserService2{
		RegisterUserFunc: func(email, password, firstName, lastName string) (*models.User, error) {
			return nil, constants.ErrEmailAlreadyExists
		},
	}

	handler := &UserHandler{userService: mockService}
	body, _ := json.Marshal(map[string]string{"email": "test@example.com", "password": "password123", "first_name": "John", "last_name": "Doe"})
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.RegisterUser(w, req)
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestUserHandler_RegisterUser_InternalError(t *testing.T) {
	mockService := &MockUserService2{
		RegisterUserFunc: func(email, password, firstName, lastName string) (*models.User, error) {
			return nil, errors.New("database error")
		},
	}

	handler := &UserHandler{userService: mockService}
	body, _ := json.Marshal(map[string]string{"email": "test@example.com", "password": "password123", "first_name": "John", "last_name": "Doe"})
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.RegisterUser(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUserHandler_LoginUser_Success(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	mockService := &MockUserService2{
		LoginFunc: func(email, password string) (*models.User, error) {
			return &models.User{ID: 1, Email: email, FirstName: "John", LastName: "Doe"}, nil
		},
	}

	handler := &UserHandler{userService: mockService}
	body, _ := json.Marshal(map[string]string{"email": "test@example.com", "password": "password123"})
	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.LoginUser(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserHandler_LoginUser_InvalidCredentials(t *testing.T) {
	mockService := &MockUserService2{
		LoginFunc: func(email, password string) (*models.User, error) {
			return nil, constants.ErrInvalidCredentials
		},
	}

	handler := &UserHandler{userService: mockService}
	body, _ := json.Marshal(map[string]string{"email": "test@example.com", "password": "wrong"})
	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.LoginUser(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUserHandler_LoginUser_ServiceError(t *testing.T) {
	mockService := &MockUserService2{
		LoginFunc: func(email, password string) (*models.User, error) {
			return nil, errors.New("database error")
		},
	}

	handler := &UserHandler{userService: mockService}
	body, _ := json.Marshal(map[string]string{"email": "test@example.com", "password": "password123"})
	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.LoginUser(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUserHandler_GetUserProfile_Success(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	testUser := &models.User{ID: 1, Email: "test@example.com", FirstName: "John", LastName: "Doe"}
	token, _ := utils.GenerateJWT(testUser)

	mockService := &MockUserService2{
		GetByIDFunc: func(userID int) (*models.User, error) {
			return testUser, nil
		},
	}

	handler := &UserHandler{userService: mockService}
	req := httptest.NewRequest(http.MethodGet, "/api/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler.GetUserProfile(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserHandler_GetUserProfile_NoAuthHeader(t *testing.T) {
	handler := &UserHandler{}
	req := httptest.NewRequest(http.MethodGet, "/api/profile", nil)
	w := httptest.NewRecorder()

	handler.GetUserProfile(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUserHandler_GetUserProfile_InvalidToken(t *testing.T) {
	handler := &UserHandler{}
	req := httptest.NewRequest(http.MethodGet, "/api/profile", nil)
	req.Header.Set("Authorization", "Bearer invalid")
	w := httptest.NewRecorder()

	handler.GetUserProfile(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUserHandler_GetUserProfile_UserNotFound(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	testUser := &models.User{ID: 1, Email: "test@example.com", FirstName: "John", LastName: "Doe"}
	token, _ := utils.GenerateJWT(testUser)

	mockService := &MockUserService2{
		GetByIDFunc: func(userID int) (*models.User, error) {
			return nil, errors.New("not found")
		},
	}

	handler := &UserHandler{userService: mockService}
	req := httptest.NewRequest(http.MethodGet, "/api/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler.GetUserProfile(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUserHandler_UpdateUserProfile_Success(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	testUser := &models.User{ID: 1, Email: "test@example.com", FirstName: "John", LastName: "Doe"}
	token, _ := utils.GenerateJWT(testUser)

	mockService := &MockUserService2{
		UpdateUserFunc: func(userID int, email, firstName, lastName *string) (*models.User, error) {
			return &models.User{ID: userID, Email: "test@example.com", FirstName: "Johnny", LastName: "Doe"}, nil
		},
	}

	handler := &UserHandler{userService: mockService}
	body, _ := json.Marshal(map[string]string{"first_name": "Johnny"})
	req := httptest.NewRequest(http.MethodPut, "/api/profile", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler.UpdateUserProfile(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserHandler_UpdateUserProfile_NoFields(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	testUser := &models.User{ID: 1, Email: "test@example.com", FirstName: "John", LastName: "Doe"}
	token, _ := utils.GenerateJWT(testUser)

	handler := &UserHandler{}
	body, _ := json.Marshal(map[string]interface{}{})
	req := httptest.NewRequest(http.MethodPut, "/api/profile", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler.UpdateUserProfile(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_UpdateUserProfile_EmailExists(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	testUser := &models.User{ID: 1, Email: "test@example.com", FirstName: "John", LastName: "Doe"}
	token, _ := utils.GenerateJWT(testUser)

	mockService := &MockUserService2{
		UpdateUserFunc: func(userID int, email, firstName, lastName *string) (*models.User, error) {
			return nil, constants.ErrEmailAlreadyExists
		},
	}

	handler := &UserHandler{userService: mockService}
	body, _ := json.Marshal(map[string]interface{}{"email": "existing@example.com"})
	req := httptest.NewRequest(http.MethodPut, "/api/profile", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler.UpdateUserProfile(w, req)
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestUserHandler_UpdateUserProfile_InternalError(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	testUser := &models.User{ID: 1, Email: "test@example.com", FirstName: "John", LastName: "Doe"}
	token, _ := utils.GenerateJWT(testUser)

	mockService := &MockUserService2{
		UpdateUserFunc: func(userID int, email, firstName, lastName *string) (*models.User, error) {
			return nil, errors.New("database error")
		},
	}

	handler := &UserHandler{userService: mockService}
	body, _ := json.Marshal(map[string]string{"first_name": "Johnny"})
	req := httptest.NewRequest(http.MethodPut, "/api/profile", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler.UpdateUserProfile(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUserHandler_LoginUser_ValidationError(t *testing.T) {
	handler := &UserHandler{}
	body, _ := json.Marshal(map[string]string{"email": "invalid", "password": "123"})
	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.LoginUser(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_LoginUser_InvalidJSON(t *testing.T) {
	handler := &UserHandler{}
	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader([]byte("invalid")))
	w := httptest.NewRecorder()

	handler.LoginUser(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_UpdateUserProfile_InvalidJSON(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	testUser := &models.User{ID: 1, Email: "test@example.com"}
	token, _ := utils.GenerateJWT(testUser)

	handler := &UserHandler{}
	req := httptest.NewRequest(http.MethodPut, "/api/profile", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler.UpdateUserProfile(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_UpdateUserProfile_ValidationError(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	testUser := &models.User{ID: 1, Email: "test@example.com"}
	token, _ := utils.GenerateJWT(testUser)

	handler := &UserHandler{}
	body, _ := json.Marshal(map[string]string{"email": "invalid-email"})
	req := httptest.NewRequest(http.MethodPut, "/api/profile", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler.UpdateUserProfile(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_RegisterUser_EmptyFields(t *testing.T) {
	handler := &UserHandler{}
	body, _ := json.Marshal(map[string]string{
		"email":     "",
		"password":  "password123",
		"firstName": "John",
		"lastName":  "Doe",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.RegisterUser(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_RegisterUser_ShortPassword(t *testing.T) {
	handler := &UserHandler{}
	body, _ := json.Marshal(map[string]string{
		"email":     "test@example.com",
		"password":  "123",
		"firstName": "John",
		"lastName":  "Doe",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.RegisterUser(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_LoginUser_EmptyEmail(t *testing.T) {
	handler := &UserHandler{}
	body, _ := json.Marshal(map[string]string{
		"email":    "",
		"password": "password123",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.LoginUser(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
