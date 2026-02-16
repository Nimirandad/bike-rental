package constants

import "errors"

// Geographical boundaries for bike locations
const (
	MinLatitude  = -90
	MaxLatitude  = 90
	MinLongitude = -180
	MaxLongitude = 180

	DefaultPage  = 1
	DefaultLimit = 20
	MaxLimit     = 100
)

// User Service Errors
var (
	ErrEmailAlreadyExists = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// Rental Service Errors
var (
	ErrBikeNotAvailable    = errors.New("bike is not available for rental")
	ErrUserHasActiveRental = errors.New("user already has an active rental")
	ErrBikeNotFound        = errors.New("bike not found")
	ErrNoActiveRental      = errors.New("you don't have an active rental")
	ErrEndLocationTooFar   = errors.New("end location must be within 5km of start location")
)
