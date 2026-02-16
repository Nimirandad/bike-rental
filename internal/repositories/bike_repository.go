package repositories

import (
	"database/sql"
	"fmt"

	"github.com/Nimirandad/bike-rental-service/internal/models"
)

type BikeRepository struct {
	db *sql.DB
}

func NewBikeRepository(db *sql.DB) *BikeRepository {
	return &BikeRepository{db: db}
}

func (r *BikeRepository) CountAvailable() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM bikes WHERE is_available = 1").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting available bikes: %w", err)
	}
	return count, nil
}

func (r *BikeRepository) GetAvailable(page, limit int) ([]*models.Bike, error) {
	offset := (page - 1) * limit

	rows, err := r.db.Query(
		"SELECT id, is_available, latitude, longitude, price_per_minute, created_at, updated_at FROM bikes WHERE is_available = 1 LIMIT ? OFFSET ?",
		limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("error querying available bikes: %w", err)
	}
	defer rows.Close()

	bikes := []*models.Bike{}
	for rows.Next() {
		var bike models.Bike
		var isAvailable int

		err := rows.Scan(&bike.ID, &isAvailable, &bike.Latitude, &bike.Longitude, &bike.PricePerMinute, &bike.CreatedAt, &bike.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning bike: %w", err)
		}

		bike.IsAvailable = isAvailable == 1
		bikes = append(bikes, &bike)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating bikes: %w", err)
	}

	return bikes, nil
}

func (r *BikeRepository) GetByID(bikeID int) (*models.Bike, error) {
	var bike models.Bike
	var isAvailable int

	err := r.db.QueryRow(
		"SELECT id, is_available, latitude, longitude, price_per_minute, created_at, updated_at FROM bikes WHERE id = ?",
		bikeID,
	).Scan(&bike.ID, &isAvailable, &bike.Latitude, &bike.Longitude, &bike.PricePerMinute, &bike.CreatedAt, &bike.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("bike with id %d not found", bikeID)
	}
	if err != nil {
		return nil, fmt.Errorf("error finding bike: %w", err)
	}

	bike.IsAvailable = isAvailable == 1
	return &bike, nil
}

func (r *BikeRepository) UpdateAvailability(bikeID int, isAvailable bool) error {
	availableInt := 0
	if isAvailable {
		availableInt = 1
	}

	_, err := r.db.Exec(
		"UPDATE bikes SET is_available = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		availableInt, bikeID,
	)
	if err != nil {
		return fmt.Errorf("error updating bike availability: %w", err)
	}

	return nil
}