package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Nimirandad/bike-rental-service/internal/models"
)

type RentalRepository struct {
	db *sql.DB
}

func NewRentalRepository(db *sql.DB) *RentalRepository {
	return &RentalRepository{db: db}
}

func (r *RentalRepository) HasActiveRental(userID int) (bool, error) {
	var count int
	err := r.db.QueryRow(
		"SELECT COUNT(*) FROM rentals WHERE user_id = ? AND status = 'running'",
		userID,
	).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking active rentals: %w", err)
	}
	return count > 0, nil
}

func (r *RentalRepository) Create(userID, bikeID int, startLat, startLong float64) (*models.Rental, error) {
	result, err := r.db.Exec(
		`INSERT INTO rentals (user_id, bike_id, status, start_time, start_latitude, start_longitude) 
		VALUES (?, ?, 'running', ?, ?, ?)`,
		userID, bikeID, time.Now(), startLat, startLong,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating rental: %w", err)
	}

	rentalID, _ := result.LastInsertId()
	return r.GetByID(int(rentalID))
}

func (r *RentalRepository) GetByID(rentalID int) (*models.Rental, error) {
	var rental models.Rental
	var endTime sql.NullTime
	var endLat, endLong sql.NullFloat64
	var durationMinutes sql.NullInt64
	var cost sql.NullFloat64

	err := r.db.QueryRow(
		`SELECT id, user_id, bike_id, status, start_time, end_time, start_latitude, 
		start_longitude, end_latitude, end_longitude, duration_minutes, cost, created_at, updated_at 
		FROM rentals WHERE id = ?`,
		rentalID,
	).Scan(
		&rental.ID, &rental.UserID, &rental.BikeID, &rental.Status,
		&rental.StartTime, &endTime, &rental.StartLatitude,
		&rental.StartLongitude, &endLat, &endLong,
		&durationMinutes, &cost,
		&rental.CreatedAt, &rental.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("rental with id %d not found", rentalID)
	}
	if err != nil {
		return nil, fmt.Errorf("error finding rental: %w", err)
	}

	if endTime.Valid {
		rental.EndTime = endTime.Time
	}
	if endLat.Valid {
		rental.EndLatitude = endLat.Float64
	}
	if endLong.Valid {
		rental.EndLongitude = endLong.Float64
	}
	if durationMinutes.Valid {
		dur := int(durationMinutes.Int64)
		rental.DurationMinutes = &dur
	}
	if cost.Valid {
		c := cost.Float64
		rental.Cost = &c
	}

	return &rental, nil
}

func (r *RentalRepository) GetActiveRentalsByUser(userID int, page, limit int) ([]*models.Rental, error) {
	offset := (page - 1) * limit

	rows, err := r.db.Query(
		`SELECT id, user_id, bike_id, status, start_time, end_time, start_latitude, 
		start_longitude, end_latitude, end_longitude, duration_minutes, cost, created_at, updated_at 
		FROM rentals WHERE user_id = ? ORDER BY id DESC LIMIT ? OFFSET ?`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("error querying rentals: %w", err)
	}
	defer rows.Close()

	rentals := []*models.Rental{}
	for rows.Next() {
		var rental models.Rental
		var endTime sql.NullTime
		var endLat, endLong sql.NullFloat64
		var durationMinutes sql.NullInt64
		var cost sql.NullFloat64

		err := rows.Scan(
			&rental.ID, &rental.UserID, &rental.BikeID, &rental.Status,
			&rental.StartTime, &endTime, &rental.StartLatitude,
			&rental.StartLongitude, &endLat, &endLong,
			&durationMinutes, &cost,
			&rental.CreatedAt, &rental.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning rental: %w", err)
		}

		if endTime.Valid {
			rental.EndTime = endTime.Time
		}
		if endLat.Valid {
			rental.EndLatitude = endLat.Float64
		}
		if endLong.Valid {
			rental.EndLongitude = endLong.Float64
		}
		if durationMinutes.Valid {
			dur := int(durationMinutes.Int64)
			rental.DurationMinutes = &dur
		}
		if cost.Valid {
			c := cost.Float64
			rental.Cost = &c
		}

		rentals = append(rentals, &rental)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rentals: %w", err)
	}

	return rentals, nil
}

func (r *RentalRepository) CountByUser(userID int) (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM rentals WHERE user_id = ?", userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting rentals: %w", err)
	}
	return count, nil
}

func (r *RentalRepository) GetActiveRentalByUser(userID int) (*models.Rental, error) {
	var rental models.Rental
	var endTime sql.NullTime
	var endLat, endLong sql.NullFloat64
	var durationMinutes sql.NullInt64
	var cost sql.NullFloat64

	err := r.db.QueryRow(
		`SELECT id, user_id, bike_id, status, start_time, end_time, start_latitude, 
		start_longitude, end_latitude, end_longitude, duration_minutes, cost, created_at, updated_at 
		FROM rentals WHERE user_id = ? AND status = 'running' LIMIT 1`,
		userID,
	).Scan(
		&rental.ID, &rental.UserID, &rental.BikeID, &rental.Status,
		&rental.StartTime, &endTime, &rental.StartLatitude,
		&rental.StartLongitude, &endLat, &endLong,
		&durationMinutes, &cost,
		&rental.CreatedAt, &rental.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error finding active rental: %w", err)
	}

	if endTime.Valid {
		rental.EndTime = endTime.Time
	}
	if endLat.Valid {
		rental.EndLatitude = endLat.Float64
	}
	if endLong.Valid {
		rental.EndLongitude = endLong.Float64
	}
	if durationMinutes.Valid {
		dur := int(durationMinutes.Int64)
		rental.DurationMinutes = &dur
	}
	if cost.Valid {
		c := cost.Float64
		rental.Cost = &c
	}

	return &rental, nil
}

func (r *RentalRepository) EndRental(rentalID int, endLat, endLong float64, durationMinutes int, cost float64) (*models.Rental, error) {
	_, err := r.db.Exec(
		`UPDATE rentals SET status = 'ended', end_time = ?, end_latitude = ?, 
		end_longitude = ?, duration_minutes = ?, cost = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		time.Now(), endLat, endLong, durationMinutes, cost, rentalID,
	)
	if err != nil {
		return nil, fmt.Errorf("error ending rental: %w", err)
	}

	return r.GetByID(rentalID)
}