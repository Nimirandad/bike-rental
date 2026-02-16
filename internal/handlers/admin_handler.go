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

type AdminService interface {
	CreateBike(latitude, longitude, pricePerMinute float64) (*models.Bike, error)
	GetAllBikes(page, limit int) ([]*models.Bike, int, error)
	UpdateBike(bikeID int, latitude, longitude *float64, isAvailable *bool, pricePerMinute *float64) (*models.Bike, error)
	GetAllUsers(page, limit int) ([]*models.User, int, error)
	GetUserByID(userID int) (*models.User, error)
	UpdateUser(userID int, email, firstName, lastName, hashedPassword *string) (*models.User, error)
	GetAllRentals(page, limit int) ([]*models.Rental, int, error)
	GetRentalByID(rentalID int) (*models.Rental, error)
	UpdateRental(rentalID int, status *string) (*models.Rental, error)
}

type AdminHandler struct {
	adminService AdminService
}

func NewAdminHandler(adminService *services.AdminService) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
	}
}

// AddBike godoc
// @Summary Add a new bike (Admin)
// @Description Create a new bike with location and pricing (requires admin authentication)
// @Tags admin
// @Accept json
// @Produce json
// @Param bike body types.AddBikeRequest true "Bike data with coordinates and price"
// @Security BasicAuth
// @Success 200 {object} types.SuccessResponse{data=models.Bike} "Bike added successfully"
// @Failure 400 {object} types.ErrorResponse "Invalid coordinates or price"
// @Failure 401 {object} types.ErrorResponse "Unauthorized - admin credentials required"
// @Failure 500 {object} types.ErrorResponse "Internal server error"
// @Router /admin/bikes [post]
func (h *AdminHandler) AddBike(w http.ResponseWriter, r *http.Request) {
	log := logger.Get()
	authHeader := r.Header.Get("Authorization")
	if err := utils.ValidateAdminBasicAuth(authHeader); err != nil {
		log.Warn().Err(err).Msg("Unauthorized admin access attempt")
		w.Header().Set("WWW-Authenticate", `Basic realm="Admin Access"`)
		types.WriteError(w, http.StatusUnauthorized, "Unauthorized: "+err.Error())
		return
	}

	var req types.AddBikeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("Failed to decode add bike request")
		types.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	log.Info().Float64("latitude", req.Latitude).Float64("longitude", req.Longitude).Msg("Admin attempting to add bike")

	if req.Latitude < constants.MinLatitude || req.Latitude > constants.MaxLatitude {
		log.Warn().Float64("latitude", req.Latitude).Msg("Invalid latitude value")
		types.WriteError(w, http.StatusBadRequest, "Latitude must be between -90 and 90")
		return
	}

	if req.Longitude < constants.MinLongitude || req.Longitude > constants.MaxLongitude {
		log.Warn().Float64("longitude", req.Longitude).Msg("Invalid longitude value")
		types.WriteError(w, http.StatusBadRequest, "Longitude must be between -180 and 180")
		return
	}

	pricePerMinute := 0.5
	if req.PricePerMinute != nil {
		if *req.PricePerMinute <= 0 {
			log.Warn().Float64("price", *req.PricePerMinute).Msg("Invalid price per minute")
			types.WriteError(w, http.StatusBadRequest, "Price per minute must be greater than 0")
			return
		}
		pricePerMinute = *req.PricePerMinute
	}

	bike, err := h.adminService.CreateBike(req.Latitude, req.Longitude, pricePerMinute)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create bike")
		types.WriteError(w, http.StatusInternalServerError, "Error creating bike")
		return
	}

	log.Info().Int("bike_id", bike.ID).Float64("latitude", bike.Latitude).Float64("longitude", bike.Longitude).Msg("Bike created successfully")
	types.WriteSuccess(w, "Bike created successfully", bike)
}

// UpdateBike godoc
// @Summary Update bike (Admin)
// @Description Update bike details like location, availability, or price (requires admin authentication)
// @Tags admin
// @Accept json
// @Produce json
// @Param bike-id path int true "Bike ID"
// @Param bike body types.UpdateBikeRequest true "Bike update data"
// @Security BasicAuth
// @Success 200 {object} types.SuccessResponse{data=models.Bike} "Bike updated successfully"
// @Failure 400 {object} types.ErrorResponse "Invalid bike ID or update data"
// @Failure 401 {object} types.ErrorResponse "Unauthorized"
// @Failure 404 {object} types.ErrorResponse "Bike not found"
// @Failure 500 {object} types.ErrorResponse "Internal server error"
// @Router /admin/bikes/{bike-id} [patch]
func (h *AdminHandler) UpdateBike(w http.ResponseWriter, r *http.Request) {
	log := logger.Get()

	authHeader := r.Header.Get("Authorization")
	if err := utils.ValidateAdminBasicAuth(authHeader); err != nil {
		log.Warn().Err(err).Msg("Unauthorized admin attempt to update bike")
		w.Header().Set("WWW-Authenticate", `Basic realm="Admin Access"`)
		types.WriteError(w, http.StatusUnauthorized, "Unauthorized: "+err.Error())
		return
	}

	bikeIDStr := r.PathValue("bike-id")
	bikeID, err := strconv.Atoi(bikeIDStr)
	if err != nil || bikeID <= 0 {
		log.Warn().Str("bike_id", bikeIDStr).Msg("Invalid bike ID in update request")
		types.WriteError(w, http.StatusBadRequest, "Invalid bike ID")
		return
	}

	var req types.UpdateBikeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Warn().Err(err).Int("bike_id", bikeID).Msg("Failed to decode update bike request")
		types.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Latitude == nil && req.Longitude == nil && req.IsAvailable == nil && req.PricePerMinute == nil {
		log.Warn().Int("bike_id", bikeID).Msg("No fields provided for update")
		types.WriteError(w, http.StatusBadRequest, "At least one field must be provided for update")
		return
	}

	if req.Latitude != nil {
		if *req.Latitude < constants.MinLatitude || *req.Latitude > constants.MaxLatitude {
			log.Warn().Float64("latitude", *req.Latitude).Int("bike_id", bikeID).Msg("Invalid latitude for bike update")
			types.WriteError(w, http.StatusBadRequest, "Latitude must be between -90 and 90")
			return
		}
	}

	if req.Longitude != nil {
		if *req.Longitude < constants.MinLongitude || *req.Longitude > constants.MaxLongitude {
			log.Warn().Float64("longitude", *req.Longitude).Int("bike_id", bikeID).Msg("Invalid longitude for bike update")
			types.WriteError(w, http.StatusBadRequest, "Longitude must be between -180 and 180")
			return
		}
	}

	if req.PricePerMinute != nil {
		if *req.PricePerMinute <= 0 {
			log.Warn().Float64("price", *req.PricePerMinute).Int("bike_id", bikeID).Msg("Invalid price for bike update")
			types.WriteError(w, http.StatusBadRequest, "Price per minute must be greater than 0")
			return
		}
	}

	log.Info().Int("bike_id", bikeID).Msg("Admin attempting to update bike")

	bike, err := h.adminService.UpdateBike(bikeID, req.Latitude, req.Longitude, req.IsAvailable, req.PricePerMinute)
	if err != nil {
		if err.Error() == "bike with id "+bikeIDStr+" not found" {
			log.Warn().Int("bike_id", bikeID).Msg("Bike not found for update")
			types.WriteError(w, http.StatusNotFound, "Bike not found")
			return
		}
		log.Error().Err(err).Int("bike_id", bikeID).Msg("Error updating bike")
		types.WriteError(w, http.StatusInternalServerError, "Error updating bike")
		return
	}

	log.Info().Int("bike_id", bikeID).Msg("Bike updated successfully by admin")
	types.WriteSuccess(w, "Bike updated successfully", bike)
}

// GetAllBikes godoc
// @Summary Get all bikes (Admin)
// @Description Get paginated list of all bikes in the system (requires admin authentication)
// @Tags admin
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page (max 100)" default(20)
// @Security BasicAuth
// @Success 200 {object} types.PaginatedResponse{data=[]models.Bike} "All bikes retrieved successfully"
// @Failure 401 {object} types.ErrorResponse "Unauthorized"
// @Failure 500 {object} types.ErrorResponse "Internal server error"
// @Router /admin/bikes [get]
func (h *AdminHandler) GetAllBikes(w http.ResponseWriter, r *http.Request) {
	log := logger.Get()

	authHeader := r.Header.Get("Authorization")
	if err := utils.ValidateAdminBasicAuth(authHeader); err != nil {
		log.Warn().Err(err).Msg("Unauthorized admin attempt to get all bikes")
		w.Header().Set("WWW-Authenticate", `Basic realm="Admin Access"`)
		types.WriteError(w, http.StatusUnauthorized, "Unauthorized: "+err.Error())
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

	log.Info().Int("page", page).Int("limit", limit).Msg("Admin fetching all bikes")

	bikes, total, err := h.adminService.GetAllBikes(page, limit)
	if err != nil {
		log.Error().Err(err).Int("page", page).Int("limit", limit).Msg("Error retrieving bikes for admin")
		types.WriteError(w, http.StatusInternalServerError, "Error retrieving bikes")
		return
	}

	log.Info().Int("total", total).Int("returned", len(bikes)).Int("page", page).Int("limit", limit).Msg("All bikes retrieved successfully by admin")
	types.WritePaginatedSuccess(w, "All bikes retrieved successfully", bikes, total, page, limit)
}

// GetAllUsers godoc
// @Summary Get all users (Admin)
// @Description Get paginated list of all users in the system (requires admin authentication)
// @Tags admin
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page (max 100)" default(20)
// @Security BasicAuth
// @Success 200 {object} types.PaginatedResponse{data=[]models.User} "All users retrieved successfully"
// @Failure 401 {object} types.ErrorResponse "Unauthorized"
// @Failure 500 {object} types.ErrorResponse "Internal server error"
// @Router /admin/users [get]
func (h *AdminHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	log := logger.Get()

	authHeader := r.Header.Get("Authorization")
	if err := utils.ValidateAdminBasicAuth(authHeader); err != nil {
		log.Warn().Err(err).Msg("Unauthorized admin attempt to get all users")
		w.Header().Set("WWW-Authenticate", `Basic realm="Admin Access"`)
		types.WriteError(w, http.StatusUnauthorized, "Unauthorized: "+err.Error())
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

	log.Info().Int("page", page).Int("limit", limit).Msg("Admin fetching all users")

	users, total, err := h.adminService.GetAllUsers(page, limit)
	if err != nil {
		log.Error().Err(err).Int("page", page).Int("limit", limit).Msg("Error retrieving users for admin")
		types.WriteError(w, http.StatusInternalServerError, "Error retrieving users")
		return
	}

	log.Info().Int("total", total).Int("returned", len(users)).Int("page", page).Int("limit", limit).Msg("All users retrieved successfully by admin")
	types.WritePaginatedSuccess(w, "All users retrieved successfully", users, total, page, limit)
}

// GetUserDetails godoc
// @Summary Get user details (Admin)
// @Description Get detailed information about a specific user (requires admin authentication)
// @Tags admin
// @Accept json
// @Produce json
// @Param user-id path int true "User ID"
// @Security BasicAuth
// @Success 200 {object} types.SuccessResponse{data=models.User} "User details retrieved successfully"
// @Failure 400 {object} types.ErrorResponse "Invalid user ID"
// @Failure 401 {object} types.ErrorResponse "Unauthorized"
// @Failure 404 {object} types.ErrorResponse "User not found"
// @Router /admin/users/{user-id} [get]
func (h *AdminHandler) GetUserDetails(w http.ResponseWriter, r *http.Request) {
	log := logger.Get()

	authHeader := r.Header.Get("Authorization")
	if err := utils.ValidateAdminBasicAuth(authHeader); err != nil {
		log.Warn().Err(err).Msg("Unauthorized admin attempt to get user details")
		w.Header().Set("WWW-Authenticate", `Basic realm="Admin Access"`)
		types.WriteError(w, http.StatusUnauthorized, "Unauthorized: "+err.Error())
		return
	}

	userIDStr := r.PathValue("user-id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		log.Warn().Str("user_id", userIDStr).Msg("Invalid user ID for admin get details")
		types.WriteError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	log.Info().Int("user_id", userID).Msg("Admin fetching user details")

	user, err := h.adminService.GetUserByID(userID)
	if err != nil {
		log.Warn().Err(err).Int("user_id", userID).Msg("User not found for admin")
		types.WriteError(w, http.StatusNotFound, "User not found")
		return
	}

	log.Info().Int("user_id", userID).Str("email", user.Email).Msg("User details retrieved successfully by admin")
	types.WriteSuccess(w, "User details retrieved successfully", user)
}

// UpdateUser godoc
// @Summary Update user (Admin)
// @Description Update user details including email, name, or password (requires admin authentication)
// @Tags admin
// @Accept json
// @Produce json
// @Param user-id path int true "User ID"
// @Param user body types.AdminUpdateUserRequest true "User update data"
// @Security BasicAuth
// @Success 200 {object} types.SuccessResponse{data=models.User} "User updated successfully"
// @Failure 400 {object} types.ErrorResponse "Invalid user ID or update data"
// @Failure 401 {object} types.ErrorResponse "Unauthorized"
// @Failure 409 {object} types.ErrorResponse "Email already exists"
// @Failure 500 {object} types.ErrorResponse "Internal server error"
// @Router /admin/users/{user-id} [patch]
func (h *AdminHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	log := logger.Get()

	authHeader := r.Header.Get("Authorization")
	if err := utils.ValidateAdminBasicAuth(authHeader); err != nil {
		log.Warn().Err(err).Msg("Unauthorized admin attempt to update user")
		w.Header().Set("WWW-Authenticate", `Basic realm="Admin Access"`)
		types.WriteError(w, http.StatusUnauthorized, "Unauthorized: "+err.Error())
		return
	}

	userIDStr := r.PathValue("user-id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		log.Warn().Str("user_id", userIDStr).Msg("Invalid user ID for admin update")
		types.WriteError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req types.AdminUpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Warn().Err(err).Int("user_id", userID).Msg("Failed to decode admin update user request")
		types.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Email == nil && req.FirstName == nil && req.LastName == nil && req.Password == nil {
		log.Warn().Int("user_id", userID).Msg("No fields provided for user update by admin")
		types.WriteError(w, http.StatusBadRequest, "At least one field must be provided for update")
		return
	}

	var hashedPassword *string
	if req.Password != nil {
		hashed, err := utils.HashPassword(*req.Password)
		if err != nil {
			log.Error().Err(err).Int("user_id", userID).Msg("Error hashing password for admin user update")
			types.WriteError(w, http.StatusInternalServerError, "Error processing password")
			return
		}
		hashedPassword = &hashed
	}

	log.Info().Int("user_id", userID).Msg("Admin attempting to update user")

	user, err := h.adminService.UpdateUser(userID, req.Email, req.FirstName, req.LastName, hashedPassword)
	if err != nil {
		if err == constants.ErrEmailAlreadyExists {
			log.Warn().Int("user_id", userID).Msg("Email already exists for admin user update")
			types.WriteError(w, http.StatusConflict, "Email already exists")
			return
		}
		log.Error().Err(err).Int("user_id", userID).Msg("Error updating user by admin")
		types.WriteError(w, http.StatusInternalServerError, "Error updating user")
		return
	}

	log.Info().Int("user_id", userID).Str("email", user.Email).Msg("User updated successfully by admin")
	types.WriteSuccess(w, "User updated successfully", user)
}

// GetAllRentals godoc
// @Summary Get all rentals (Admin)
// @Description Get paginated list of all rentals in the system (requires admin authentication)
// @Tags admin
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page (max 100)" default(20)
// @Security BasicAuth
// @Success 200 {object} types.PaginatedResponse{data=[]models.Rental} "All rentals retrieved successfully"
// @Failure 401 {object} types.ErrorResponse "Unauthorized"
// @Failure 500 {object} types.ErrorResponse "Internal server error"
// @Router /admin/rentals [get]
func (h *AdminHandler) GetAllRentals(w http.ResponseWriter, r *http.Request) {
	log := logger.Get()

	authHeader := r.Header.Get("Authorization")
	if err := utils.ValidateAdminBasicAuth(authHeader); err != nil {
		log.Warn().Err(err).Msg("Unauthorized admin attempt to get all rentals")
		w.Header().Set("WWW-Authenticate", `Basic realm="Admin Access"`)
		types.WriteError(w, http.StatusUnauthorized, "Unauthorized: "+err.Error())
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

	log.Info().Int("page", page).Int("limit", limit).Msg("Admin fetching all rentals")

	rentals, total, err := h.adminService.GetAllRentals(page, limit)
	if err != nil {
		log.Error().Err(err).Int("page", page).Int("limit", limit).Msg("Error retrieving rentals for admin")
		types.WriteError(w, http.StatusInternalServerError, "Error retrieving rentals")
		return
	}

	log.Info().Int("total", total).Int("returned", len(rentals)).Int("page", page).Int("limit", limit).Msg("All rentals retrieved successfully by admin")
	types.WritePaginatedSuccess(w, "All rentals retrieved successfully", rentals, total, page, limit)
}

// GetRentalDetails godoc
// @Summary Get rental details (Admin)
// @Description Get detailed information about a specific rental (requires admin authentication)
// @Tags admin
// @Accept json
// @Produce json
// @Param rental-id path int true "Rental ID"
// @Security BasicAuth
// @Success 200 {object} types.SuccessResponse{data=models.Rental} "Rental details retrieved successfully"
// @Failure 400 {object} types.ErrorResponse "Invalid rental ID"
// @Failure 401 {object} types.ErrorResponse "Unauthorized"
// @Failure 404 {object} types.ErrorResponse "Rental not found"
// @Router /admin/rentals/{rental-id} [get]
func (h *AdminHandler) GetRentalDetails(w http.ResponseWriter, r *http.Request) {
	log := logger.Get()

	authHeader := r.Header.Get("Authorization")
	if err := utils.ValidateAdminBasicAuth(authHeader); err != nil {
		log.Warn().Err(err).Msg("Unauthorized admin attempt to get rental details")
		w.Header().Set("WWW-Authenticate", `Basic realm="Admin Access"`)
		types.WriteError(w, http.StatusUnauthorized, "Unauthorized: "+err.Error())
		return
	}

	rentalIDStr := r.PathValue("rental-id")
	rentalID, err := strconv.Atoi(rentalIDStr)
	if err != nil || rentalID <= 0 {
		log.Warn().Str("rental_id", rentalIDStr).Msg("Invalid rental ID for admin get details")
		types.WriteError(w, http.StatusBadRequest, "Invalid rental ID")
		return
	}

	log.Info().Int("rental_id", rentalID).Msg("Admin fetching rental details")

	rental, err := h.adminService.GetRentalByID(rentalID)
	if err != nil {
		log.Warn().Err(err).Int("rental_id", rentalID).Msg("Rental not found for admin")
		types.WriteError(w, http.StatusNotFound, "Rental not found")
		return
	}

	log.Info().Int("rental_id", rentalID).Int("user_id", rental.UserID).Int("bike_id", rental.BikeID).Msg("Rental details retrieved successfully by admin")
	types.WriteSuccess(w, "Rental details retrieved successfully", rental)
}

// UpdateRental godoc
// @Summary Update rental (Admin)
// @Description Update rental status (requires admin authentication)
// @Tags admin
// @Accept json
// @Produce json
// @Param rental-id path int true "Rental ID"
// @Param rental body types.UpdateRentalRequest true "Rental update data with status"
// @Security BasicAuth
// @Success 200 {object} types.SuccessResponse{data=models.Rental} "Rental updated successfully"
// @Failure 400 {object} types.ErrorResponse "Invalid rental ID or status required"
// @Failure 401 {object} types.ErrorResponse "Unauthorized"
// @Failure 500 {object} types.ErrorResponse "Internal server error"
// @Router /admin/rentals/{rental-id} [patch]
func (h *AdminHandler) UpdateRental(w http.ResponseWriter, r *http.Request) {
	log := logger.Get()

	authHeader := r.Header.Get("Authorization")
	if err := utils.ValidateAdminBasicAuth(authHeader); err != nil {
		log.Warn().Err(err).Msg("Unauthorized admin attempt to update rental")
		w.Header().Set("WWW-Authenticate", `Basic realm="Admin Access"`)
		types.WriteError(w, http.StatusUnauthorized, "Unauthorized: "+err.Error())
		return
	}

	rentalIDStr := r.PathValue("rental-id")
	rentalID, err := strconv.Atoi(rentalIDStr)
	if err != nil || rentalID <= 0 {
		log.Warn().Str("rental_id", rentalIDStr).Msg("Invalid rental ID for admin update")
		types.WriteError(w, http.StatusBadRequest, "Invalid rental ID")
		return
	}

	var req types.UpdateRentalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Warn().Err(err).Int("rental_id", rentalID).Msg("Failed to decode admin update rental request")
		types.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Status == nil {
		log.Warn().Int("rental_id", rentalID).Msg("No status provided for rental update by admin")
		types.WriteError(w, http.StatusBadRequest, "Status field is required")
		return
	}

	log.Info().Int("rental_id", rentalID).Str("new_status", *req.Status).Msg("Admin attempting to update rental")

	rental, err := h.adminService.UpdateRental(rentalID, req.Status)
	if err != nil {
		log.Error().Err(err).Int("rental_id", rentalID).Str("status", *req.Status).Msg("Error updating rental by admin")
		types.WriteError(w, http.StatusInternalServerError, "Error updating rental")
		return
	}

	log.Info().Int("rental_id", rentalID).Str("status", rental.Status).Msg("Rental updated successfully by admin")
	types.WriteSuccess(w, "Rental updated successfully", rental)
}
