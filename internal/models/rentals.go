package models

import "time"

type Rental struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	BikeID          int       `json:"bike_id"`
	Status          string    `json:"status"`
	StartTime       time.Time `json:"start_time"`
	EndTime         time.Time `json:"end_time,omitempty"`
	StartLatitude   float64   `json:"start_latitude"`
	StartLongitude  float64   `json:"start_longitude"`
	EndLatitude     float64   `json:"end_latitude,omitempty"`
	EndLongitude    float64   `json:"end_longitude,omitempty"`
	DurationMinutes *int      `json:"duration_minutes,omitempty"`
	Cost            *float64  `json:"cost,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (r *Rental) TableName() string {
	return "rentals"
}