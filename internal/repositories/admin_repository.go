package repositories

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Nimirandad/bike-rental-service/internal/models"
)

type AdminRepository struct {
	db *sql.DB
}

func NewAdminRepository(db *sql.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) CreateBike(latitude, longitude, pricePerMinute float64) (*models.Bike, error) {
	result, err := r.db.Exec(
		"INSERT INTO bikes (is_available, latitude, longitude, price_per_minute) VALUES (?, ?, ?, ?)",
		1, latitude, longitude, pricePerMinute,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating bike: %w", err)
	}

	bikeID, _ := result.LastInsertId()
	return r.GetBikeByID(int(bikeID))
}

func (r *AdminRepository) GetBikeByID(bikeID int) (*models.Bike, error) {
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

func (r *AdminRepository) CountAll() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM bikes").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting bikes: %w", err)
	}
	return count, nil
}

func (r *AdminRepository) GetAllBikes(page, limit int) ([]*models.Bike, error) {
	offset := (page - 1) * limit

	rows, err := r.db.Query(
		"SELECT id, is_available, latitude, longitude, price_per_minute, created_at, updated_at FROM bikes ORDER BY id ASC LIMIT ? OFFSET ?",
		limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("error querying bikes: %w", err)
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

func (r *AdminRepository) UpdateBike(bikeID int, latitude, longitude *float64, isAvailable *bool, pricePerMinute *float64) (*models.Bike, error) {
	_, err := r.GetBikeByID(bikeID)
	if err != nil {
		return nil, err
	}

	query := "UPDATE bikes SET "
	params := []interface{}{}
	updates := []string{}

	if latitude != nil {
		updates = append(updates, "latitude = ?")
		params = append(params, *latitude)
	}

	if longitude != nil {
		updates = append(updates, "longitude = ?")
		params = append(params, *longitude)
	}

	if isAvailable != nil {
		updates = append(updates, "is_available = ?")
		var availableInt int
		if *isAvailable {
			availableInt = 1
		} else {
			availableInt = 0
		}
		params = append(params, availableInt)
	}

	if pricePerMinute != nil {
		updates = append(updates, "price_per_minute = ?")
		params = append(params, *pricePerMinute)
	}

	updates = append(updates, "updated_at = CURRENT_TIMESTAMP")

	if len(params) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	for i, update := range updates {
		if i > 0 {
			query += ", "
		}
		query += update
	}
	query += " WHERE id = ?"
	params = append(params, bikeID)

	_, err = r.db.Exec(query, params...)
	if err != nil {
		return nil, fmt.Errorf("error updating bike: %w", err)
	}

	return r.GetBikeByID(bikeID)
}

func (r *AdminRepository) GetAllUsers(page, limit int) ([]*models.User, error) {
	offset := (page - 1) * limit

	rows, err := r.db.Query(
		"SELECT id, email, first_name, last_name, created_at FROM users ORDER BY id ASC LIMIT ? OFFSET ?",
		limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("error querying users: %w", err)
	}
	defer rows.Close()

	users := []*models.User{}
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning user: %w", err)
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	return users, nil
}

func (r *AdminRepository) CountAllUsers() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting users: %w", err)
	}
	return count, nil
}

func (r *AdminRepository) GetUserByID(userID int) (*models.User, error) {
	var user models.User

	err := r.db.QueryRow(
		"SELECT id, email, first_name, last_name, created_at FROM users WHERE id = ?",
		userID,
	).Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user with id %d not found", userID)
	}
	if err != nil {
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	return &user, nil
}

func (r *AdminRepository) EmailExistsByOtherUser(email string, userID int) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ? AND id != ?)", email, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking email existence: %w", err)
	}
	return exists, nil
}

func (r *AdminRepository) UpdateUser(userID int, email, firstName, lastName, hashedPassword *string) (*models.User, error) {
	query := "UPDATE users SET "
	args := []interface{}{}
	updates := []string{}

	if email != nil {
		updates = append(updates, "email = ?")
		args = append(args, *email)
	}

	if firstName != nil {
		updates = append(updates, "first_name = ?")
		args = append(args, *firstName)
	}

	if lastName != nil {
		updates = append(updates, "last_name = ?")
		args = append(args, *lastName)
	}

	if hashedPassword != nil {
		updates = append(updates, "hashed_password = ?")
		args = append(args, *hashedPassword)
	}

	if len(updates) == 0 {
		return r.GetUserByID(userID)
	}

	updates = append(updates, "updated_at = CURRENT_TIMESTAMP")

	query += strings.Join(updates, ", ") + " WHERE id = ?"
	args = append(args, userID)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	return r.GetUserByID(userID)
}

func (r *AdminRepository) GetAllRentals(page, limit int) ([]*models.Rental, error) {
	offset := (page - 1) * limit

	rows, err := r.db.Query(
		`SELECT id, user_id, bike_id, status, start_time, end_time, start_latitude, 
		start_longitude, end_latitude, end_longitude, duration_minutes, cost, created_at, updated_at 
		FROM rentals ORDER BY id ASC LIMIT ? OFFSET ?`,
		limit, offset,
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

func (r *AdminRepository) CountAllRentals() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM rentals").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting rentals: %w", err)
	}
	return count, nil
}

func (r *AdminRepository) GetRentalByID(rentalID int) (*models.Rental, error) {
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

func (r *AdminRepository) UpdateRental(rentalID int, status *string) (*models.Rental, error) {
	if status == nil {
		return r.GetRentalByID(rentalID)
	}

	_, err := r.db.Exec(
		"UPDATE rentals SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		*status, rentalID,
	)
	if err != nil {
		return nil, fmt.Errorf("error updating rental: %w", err)
	}

	return r.GetRentalByID(rentalID)
}