package handlers

import (
	"net/http"
	"strconv"

	"github.com/Nimirandad/bike-rental-service/internal/constants"
	"github.com/Nimirandad/bike-rental-service/internal/logger"
	"github.com/Nimirandad/bike-rental-service/internal/models"
	"github.com/Nimirandad/bike-rental-service/internal/services"
	"github.com/Nimirandad/bike-rental-service/internal/types"
	"github.com/Nimirandad/bike-rental-service/internal/utils"
)

type BikeService interface {
	GetAvailableBikes(page, limit int) ([]*models.Bike, int, error)
}

type BikeHandler struct {
	bikeService BikeService
}

func NewBikeHandler(bikeService *services.BikeService) *BikeHandler {
	return &BikeHandler{bikeService: bikeService}
}

// GetAvailableBikes godoc
// @Summary List available bikes
// @Description Get paginated list of available bikes for rent
// @Tags bikes
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page (max 100)" default(20)
// @Security BearerAuth
// @Success 200 {object} types.PaginatedResponse{data=[]models.Bike} "List of available bikes"
// @Failure 401 {object} types.ErrorResponse "Unauthorized - missing or invalid token"
// @Failure 500 {object} types.ErrorResponse "Internal server error"
// @Router /bikes/available [get]
func (h *BikeHandler) GetAvailableBikes(w http.ResponseWriter, r *http.Request) {
	log := logger.Get()

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		log.Warn().Msg("Attempt to get bikes without authorization header")
		types.WriteError(w, http.StatusUnauthorized, "Authorization header is required")
		return
	}

	tokenString, err := utils.ExtractTokenFromHeader(authHeader)
	if err != nil {
		log.Warn().Err(err).Msg("Invalid authorization header format")
		types.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	_, err = utils.ValidateJWT(tokenString)
	if err != nil {
		log.Warn().Err(err).Msg("Invalid or expired token")
		types.WriteError(w, http.StatusUnauthorized, "Invalid or expired token")
		return
	}

	page := constants.DefaultPage
	if pageParam := r.URL.Query().Get("page"); pageParam != "" {
		if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
			page = p
		}
	}

	limit := constants.DefaultLimit
	if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil && l > 0 && l <= constants.MaxLimit {
			limit = l
		}
	}

	log.Info().Int("page", page).Int("limit", limit).Msg("Fetching available bikes")

	bikes, total, err := h.bikeService.GetAvailableBikes(page, limit)
	if err != nil {
		log.Error().Err(err).Int("page", page).Int("limit", limit).Msg("Error retrieving bikes")
		types.WriteError(w, http.StatusInternalServerError, "Error retrieving bikes")
		return
	}

	log.Info().Int("total", total).Int("returned", len(bikes)).Int("page", page).Int("limit", limit).Msg("Available bikes retrieved successfully")
	types.WritePaginatedSuccess(w, "Available bikes retrieved successfully", bikes, total, page, limit)
}
