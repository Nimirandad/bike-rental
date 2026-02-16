package types

type RegisterUserRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	Email     *string `json:"email,omitempty"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
}

type AddBikeRequest struct {
	Latitude       float64  `json:"latitude"`
	Longitude      float64  `json:"longitude"`
	PricePerMinute *float64 `json:"price_per_minute,omitempty"`
}

type UpdateBikeRequest struct {
	Latitude       *float64 `json:"latitude,omitempty"`
	Longitude      *float64 `json:"longitude,omitempty"`
	IsAvailable    *bool    `json:"is_available,omitempty"`
	PricePerMinute *float64 `json:"price_per_minute,omitempty"`
}

type AdminUpdateUserRequest struct {
	Email     *string `json:"email,omitempty"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Password  *string `json:"password,omitempty"`
}

type UpdateRentalRequest struct {
	Status *string `json:"status,omitempty"`
}

type StartRentalRequest struct {
	BikeID int `json:"bike_id"`
}

type EndRentalRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}