package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Nimirandad/bike-rental-service/internal/constants"
	"github.com/Nimirandad/bike-rental-service/internal/logger"
	"github.com/Nimirandad/bike-rental-service/internal/models"
	"github.com/Nimirandad/bike-rental-service/internal/services"
	"github.com/Nimirandad/bike-rental-service/internal/types"
	"github.com/Nimirandad/bike-rental-service/internal/utils"
)

type UserService interface {
	RegisterUser(email, password, firstName, lastName string) (*models.User, error)
	Login(email, password string) (*models.User, error)
	GetByID(userID int) (*models.User, error)
	UpdateUser(userID int, email, firstName, lastName *string) (*models.User, error)
}

type UserHandler struct {
	userService UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// RegisterUser godoc
// @Summary Register a new user
// @Description Register a new user with email, password, first name and last name
// @Tags users
// @Accept json
// @Produce json
// @Param user body types.RegisterUserRequest true "User registration data"
// @Success 200 {object} types.SuccessResponse{data=models.User} "User registered successfully"
// @Failure 400 {object} types.ErrorResponse "Invalid request payload"
// @Failure 409 {object} types.ErrorResponse "Email already exists"
// @Failure 500 {object} types.ErrorResponse "Internal server error"
// @Router /users/register [post]
func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	log := logger.Get()
	var req types.RegisterUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("Failed to decode register request")
		types.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	log.Info().Str("email", req.Email).Msg("User registration attempt")

	if validationErrors := utils.ValidateRegisterUserRequest(req.Email, req.Password, req.FirstName, req.LastName); len(validationErrors) > 0 {
		log.Warn().Interface("validation_errors", validationErrors).Msg("Registration validation failed")
		types.WriteValidationErrors(w, validationErrors)
		return
	}

	user, err := h.userService.RegisterUser(req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		if err == constants.ErrEmailAlreadyExists {
			log.Warn().Str("email", req.Email).Msg("Registration failed: email already exists")
			types.WriteError(w, http.StatusConflict, err.Error())
			return
		}
		log.Error().Err(err).Str("email", req.Email).Msg("Failed to register user")
		types.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Info().Int("user_id", user.ID).Str("email", user.Email).Msg("User registered successfully")
	types.WriteSuccess(w, "User registered successfully", user)
}

// LoginUser godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body types.LoginRequest true "Login credentials"
// @Success 200 {object} types.SuccessResponse{data=types.LoginResponse} "Login successful with JWT token"
// @Failure 400 {object} types.ErrorResponse "Invalid request payload"
// @Failure 401 {object} types.ErrorResponse "Invalid credentials"
// @Failure 500 {object} types.ErrorResponse "Internal server error"
// @Router /users/login [post]
func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	log := logger.Get()
	var req types.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("Failed to decode login request")
		types.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	log.Info().Str("email", req.Email).Msg("Login attempt")

	if validationErrors := utils.ValidateLoginRequest(req.Email, req.Password); len(validationErrors) > 0 {
		log.Warn().Interface("validation_errors", validationErrors).Msg("Login validation failed")
		types.WriteValidationErrors(w, validationErrors)
		return
	}

	user, err := h.userService.Login(req.Email, req.Password)
	if err != nil {
		if err == constants.ErrInvalidCredentials {
			log.Warn().Str("email", req.Email).Msg("Login failed: invalid credentials")
			types.WriteError(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}
		log.Error().Err(err).Str("email", req.Email).Msg("Login error")
		types.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	token, err := utils.GenerateJWT(user)
	if err != nil {
		log.Error().Err(err).Int("user_id", user.ID).Msg("Failed to generate JWT token")
		types.WriteError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	log.Info().Int("user_id", user.ID).Str("email", user.Email).Msg("Login successful")

	loginResponse := types.LoginResponse{
		Token: token,
	}

	types.WriteSuccess(w, "Login successful", loginResponse)
}

// GetUserProfile godoc
// @Summary Get user profile
// @Description Get authenticated user's profile information
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} types.SuccessResponse{data=models.User} "User profile retrieved successfully"
// @Failure 401 {object} types.ErrorResponse "Unauthorized - missing or invalid token"
// @Failure 404 {object} types.ErrorResponse "User not found"
// @Router /users/profile [get]
func (h *UserHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	log := logger.Get()

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		log.Warn().Msg("Attempt to get profile without authorization header")
		types.WriteError(w, http.StatusUnauthorized, "Authorization header is required")
		return
	}

	tokenString, err := utils.ExtractTokenFromHeader(authHeader)
	if err != nil {
		log.Warn().Err(err).Msg("Invalid authorization header format")
		types.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	claims, err := utils.ValidateJWT(tokenString)
	if err != nil {
		log.Warn().Err(err).Msg("Invalid or expired token for get profile")
		types.WriteError(w, http.StatusUnauthorized, "Invalid or expired token")
		return
	}

	log.Info().Int("user_id", claims.Sub).Msg("Fetching user profile")

	user, err := h.userService.GetByID(claims.Sub)
	if err != nil {
		log.Error().Err(err).Int("user_id", claims.Sub).Msg("User not found")
		types.WriteError(w, http.StatusNotFound, "User not found")
		return
	}

	log.Info().Int("user_id", claims.Sub).Str("email", user.Email).Msg("User profile retrieved successfully")
	types.WriteSuccess(w, "User profile retrieved successfully", user)
}

// UpdateUserProfile godoc
// @Summary Update user profile
// @Description Update authenticated user's profile (email, first name, last name)
// @Tags users
// @Accept json
// @Produce json
// @Param user body types.UpdateUserRequest true "Profile update data"
// @Security BearerAuth
// @Success 200 {object} types.SuccessResponse{data=models.User} "Profile updated successfully"
// @Failure 400 {object} types.ErrorResponse "Invalid request payload or validation errors"
// @Failure 401 {object} types.ErrorResponse "Unauthorized"
// @Failure 409 {object} types.ErrorResponse "Email already in use"
// @Failure 500 {object} types.ErrorResponse "Internal server error"
// @Router /users/profile [patch]
func (h *UserHandler) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	log := logger.Get()

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		log.Warn().Msg("Attempt to update profile without authorization header")
		types.WriteError(w, http.StatusUnauthorized, "Authorization header is required")
		return
	}

	tokenString, err := utils.ExtractTokenFromHeader(authHeader)
	if err != nil {
		log.Warn().Err(err).Msg("Invalid authorization header format")
		types.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	claims, err := utils.ValidateJWT(tokenString)
	if err != nil {
		log.Warn().Err(err).Msg("Invalid or expired token for update profile")
		types.WriteError(w, http.StatusUnauthorized, "Invalid or expired token")
		return
	}

	var req types.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Warn().Err(err).Int("user_id", claims.Sub).Msg("Failed to decode update profile request")
		types.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Email == nil && req.FirstName == nil && req.LastName == nil {
		log.Warn().Int("user_id", claims.Sub).Msg("No fields provided for update")
		types.WriteError(w, http.StatusBadRequest, "At least one field must be provided for update")
		return
	}

	if validationErrors := utils.ValidateUpdateUserRequest(req.Email, req.FirstName, req.LastName); len(validationErrors) > 0 {
		log.Warn().Int("user_id", claims.Sub).Interface("errors", validationErrors).Msg("Update user validation failed")
		types.WriteValidationErrors(w, validationErrors)
		return
	}

	log.Info().Int("user_id", claims.Sub).Msg("Attempting to update user profile")

	user, err := h.userService.UpdateUser(claims.Sub, req.Email, req.FirstName, req.LastName)
	if err != nil {
		if err == constants.ErrEmailAlreadyExists {
			log.Warn().Int("user_id", claims.Sub).Msg("Email already in use")
			types.WriteError(w, http.StatusConflict, "Email already in use")
			return
		}
		log.Error().Err(err).Int("user_id", claims.Sub).Msg("Error updating profile")
		types.WriteError(w, http.StatusInternalServerError, "Error updating profile")
		return
	}

	log.Info().Int("user_id", claims.Sub).Str("email", user.Email).Msg("Profile updated successfully")
	types.WriteSuccess(w, "Profile updated successfully", user)
}