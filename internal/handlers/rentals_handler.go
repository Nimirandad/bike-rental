package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Nimirandad/bike-rental-service/internal/constants"
	"github.com/Nimirandad/bike-rental-service/internal/logger"
	"github.com/Nimirandad/bike-rental-service/internal/models"
	"github.com/Nimirandad/bike-rental-service/internal/services"
	"github.com/Nimirandad/bike-rental-service/internal/types"
	"github.com/Nimirandad/bike-rental-service/internal/utils"
)

type RentalService interface {
	StartRental(userID, bikeID int) (*models.Rental, error)
	EndRental(userID int, endLat, endLong float64) (*models.Rental, error)
	GetRentalHistory(userID, page, limit int) ([]*models.Rental, int, error)
}

type RentalHandler struct {
	rentalService RentalService
}

func NewRentalHandler(rentalService *services.RentalService) *RentalHandler {
	return &RentalHandler{rentalService: rentalService}
}

// StartRental godoc
// @Summary Start a bike rental
// @Description Start a new rental for the authenticated user with specified bike
// @Tags rentals
// @Accept json
// @Produce json
// @Param rental body types.StartRentalRequest true "Rental start data with bike_id"
// @Security BearerAuth
// @Success 200 {object} types.SuccessResponse{data=models.Rental} "Rental started successfully"
// @Failure 400 {object} types.ErrorResponse "Invalid request payload or bike_id"
// @Failure 401 {object} types.ErrorResponse "Unauthorized"
// @Failure 404 {object} types.ErrorResponse "Bike not found"
// @Failure 409 {object} types.ErrorResponse "User has active rental or bike not available"
// @Failure 500 {object} types.ErrorResponse "Internal server error"
// @Router /rentals/start [post]
func (h *RentalHandler) StartRental(w http.ResponseWriter, r *http.Request) {
	log := logger.Get()
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		log.Warn().Msg("Start rental: missing authorization header")
		types.WriteError(w, http.StatusUnauthorized, "Authorization header is required")
		return
	}

	tokenString, err := utils.ExtractTokenFromHeader(authHeader)
	if err != nil {
		log.Warn().Err(err).Msg("Start rental: invalid authorization format")
		types.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	claims, err := utils.ValidateJWT(tokenString)
	if err != nil {
		log.Warn().Err(err).Msg("Start rental: invalid or expired token")
		types.WriteError(w, http.StatusUnauthorized, "Invalid or expired token")
		return
	}

	userID := claims.Sub

	var req types.StartRentalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Int("user_id", userID).Msg("Failed to decode start rental request")
		types.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	log.Info().Int("user_id", userID).Int("bike_id", req.BikeID).Msg("Attempting to start rental")

	if req.BikeID <= 0 {
		log.Warn().Int("user_id", userID).Int("bike_id", req.BikeID).Msg("Invalid bike_id")
		types.WriteError(w, http.StatusBadRequest, "Valid bike_id is required")
		return
	}

	rental, err := h.rentalService.StartRental(userID, req.BikeID)
	if err != nil {
		if err == constants.ErrUserHasActiveRental {
			log.Warn().Int("user_id", userID).Msg("User already has active rental")
			types.WriteError(w, http.StatusConflict, "You already have an active rental. Please return your current bike before renting another.")
			return
		}
		if err == constants.ErrBikeNotAvailable {
			log.Warn().Int("bike_id", req.BikeID).Msg("Bike not available")
			types.WriteError(w, http.StatusConflict, "This bike is already rented by another user")
			return
		}
		if err == constants.ErrBikeNotFound {
			log.Warn().Int("bike_id", req.BikeID).Msg("Bike not found")
			types.WriteError(w, http.StatusNotFound, "Bike not found")
			return
		}
		log.Error().Err(err).Int("user_id", userID).Int("bike_id", req.BikeID).Msg("Failed to start rental")
		types.WriteError(w, http.StatusInternalServerError, "Error starting rental")
		return
	}

	log.Info().Int("rental_id", rental.ID).Int("user_id", userID).Int("bike_id", req.BikeID).Msg("Rental started successfully")
	types.WriteSuccess(w, "Rental started successfully", rental)
}

// EndRental godoc
// @Summary End a bike rental
// @Description End the active rental for authenticated user with end location
// @Tags rentals
// @Accept json
// @Produce json
// @Param rental body types.EndRentalRequest true "End location coordinates"
// @Security BearerAuth
// @Success 200 {object} types.SuccessResponse{data=models.Rental} "Rental ended successfully with cost"
// @Failure 400 {object} types.ErrorResponse "Invalid coordinates or location too far from start"
// @Failure 401 {object} types.ErrorResponse "Unauthorized"
// @Failure 409 {object} types.ErrorResponse "No active rental to end"
// @Failure 500 {object} types.ErrorResponse "Internal server error"
// @Router /rentals/end [post]
func (h *RentalHandler) EndRental(w http.ResponseWriter, r *http.Request) {
	log := logger.Get()

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		log.Warn().Msg("Attempt to end rental without authorization header")
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
		log.Warn().Err(err).Msg("Invalid or expired token for end rental")
		types.WriteError(w, http.StatusUnauthorized, "Invalid or expired token")
		return
	}

	userID := claims.Sub

	var req types.EndRentalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Warn().Err(err).Int("user_id", userID).Msg("Failed to decode end rental request")
		types.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Latitude < -90 || req.Latitude > 90 {
		log.Warn().Float64("latitude", req.Latitude).Int("user_id", userID).Msg("Invalid latitude for end rental")
		types.WriteError(w, http.StatusBadRequest, "Latitude must be between -90 and 90")
		return
	}

	if req.Longitude < -180 || req.Longitude > 180 {
		log.Warn().Float64("longitude", req.Longitude).Int("user_id", userID).Msg("Invalid longitude for end rental")
		types.WriteError(w, http.StatusBadRequest, "Longitude must be between -180 and 180")
		return
	}

	log.Info().Int("user_id", userID).Float64("latitude", req.Latitude).Float64("longitude", req.Longitude).Msg("Attempting to end rental")

	rental, err := h.rentalService.EndRental(userID, req.Latitude, req.Longitude)
	if err != nil {
		if err == constants.ErrNoActiveRental {
			log.Warn().Int("user_id", userID).Msg("User has no active rental to end")
			types.WriteError(w, http.StatusConflict, "You don't have an active rental to end")
			return
		}
		if err == constants.ErrEndLocationTooFar {
			log.Warn().Int("user_id", userID).Float64("latitude", req.Latitude).Float64("longitude", req.Longitude).Msg("End location too far from start")
			types.WriteError(w, http.StatusBadRequest, "End location must be within 5km of the start location")
			return
		}
		log.Error().Err(err).Int("user_id", userID).Msg("Failed to end rental")
		types.WriteError(w, http.StatusInternalServerError, "Error ending rental")
		return
	}

	logEvent := log.Info().Int("rental_id", rental.ID).Int("user_id", userID)
	if rental.Cost != nil {
		logEvent = logEvent.Float64("cost", *rental.Cost)
	}
	logEvent.Msg("Rental ended successfully")
	types.WriteSuccess(w, "Rental ended successfully", rental)
}

// GetRentalHistory godoc
// @Summary Get rental history
// @Description Get paginated rental history for authenticated user
// @Tags rentals
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page (max 100)" default(20)
// @Security BearerAuth
// @Success 200 {object} types.PaginatedResponse{data=[]models.Rental} "Rental history retrieved successfully"
// @Failure 401 {object} types.ErrorResponse "Unauthorized"
// @Failure 500 {object} types.ErrorResponse "Internal server error"
// @Router /rentals/history [get]
func (h *RentalHandler) GetRentalHistory(w http.ResponseWriter, r *http.Request) {
	log := logger.Get()

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		log.Warn().Msg("Attempt to get rental history without authorization header")
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
		log.Warn().Err(err).Msg("Invalid or expired token for rental history")
		types.WriteError(w, http.StatusUnauthorized, "Invalid or expired token")
		return
	}

	userID := claims.Sub

	page := 1
	if pageParam := r.URL.Query().Get("page"); pageParam != "" {
		if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
			page = p
		}
	}

	limit := 20
	if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	log.Info().Int("user_id", userID).Int("page", page).Int("limit", limit).Msg("Fetching rental history")

	rentals, total, err := h.rentalService.GetRentalHistory(userID, page, limit)
	if err != nil {
		log.Error().Err(err).Int("user_id", userID).Int("page", page).Int("limit", limit).Msg("Error retrieving rental history")
		types.WriteError(w, http.StatusInternalServerError, "Error retrieving rental history")
		return
	}

	log.Info().Int("user_id", userID).Int("total", total).Int("returned", len(rentals)).Int("page", page).Int("limit", limit).Msg("Rental history retrieved successfully")
	types.WritePaginatedSuccess(w, "Rental history retrieved successfully", rentals, total, page, limit)
}
