package models

import "time"

type Bike struct {
	ID             int       `json:"id"`
	IsAvailable    bool      `json:"is_available"`
	Latitude       float64   `json:"latitude"`
	Longitude      float64   `json:"longitude"`
	PricePerMinute float64   `json:"price_per_minute"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}

func (b *Bike) TableName() string {
	return "bikes"
}