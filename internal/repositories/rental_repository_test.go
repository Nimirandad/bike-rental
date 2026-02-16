package repositories

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestRentalRepository_HasActiveRental(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRentalRepository(db)

	t.Run("User has active rental", func(t *testing.T) {
		mock.ExpectQuery("SELECT COUNT").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		hasActive, err := repo.HasActiveRental(1)

		assert.NoError(t, err)
		assert.True(t, hasActive)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("User has no active rental", func(t *testing.T) {
		mock.ExpectQuery("SELECT COUNT").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		hasActive, err := repo.HasActiveRental(1)

		assert.NoError(t, err)
		assert.False(t, hasActive)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Database error", func(t *testing.T) {
		mock.ExpectQuery("SELECT COUNT").
			WithArgs(1).
			WillReturnError(fmt.Errorf("database error"))

		hasActive, err := repo.HasActiveRental(1)

		assert.Error(t, err)
		assert.False(t, hasActive)
		assert.Contains(t, err.Error(), "error checking active rentals")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRentalRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRentalRepository(db)
	now := time.Now()

	t.Run("Successfully create rental", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO rentals").
			WithArgs(1, 10, sqlmock.AnyArg(), 40.7128, -74.0060).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery("SELECT id, user_id, bike_id, status").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "bike_id", "status", "start_time", "end_time", "start_latitude", "start_longitude", "end_latitude", "end_longitude", "duration_minutes", "cost", "created_at", "updated_at"}).
				AddRow(1, 1, 10, "running", now, nil, 40.7128, -74.0060, nil, nil, nil, nil, now, now))

		rental, err := repo.Create(1, 10, 40.7128, -74.0060)

		assert.NoError(t, err)
		assert.NotNil(t, rental)
		assert.Equal(t, 1, rental.ID)
		assert.Equal(t, 1, rental.UserID)
		assert.Equal(t, 10, rental.BikeID)
		assert.Equal(t, "running", rental.Status)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Database error on insert", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO rentals").
			WithArgs(1, 10, sqlmock.AnyArg(), 40.7128, -74.0060).
			WillReturnError(fmt.Errorf("database error"))

		rental, err := repo.Create(1, 10, 40.7128, -74.0060)

		assert.Error(t, err)
		assert.Nil(t, rental)
		assert.Contains(t, err.Error(), "error creating rental")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRentalRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRentalRepository(db)
	now := time.Now()

	t.Run("Rental found - running status", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, user_id, bike_id, status").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "bike_id", "status", "start_time", "end_time", "start_latitude", "start_longitude", "end_latitude", "end_longitude", "duration_minutes", "cost", "created_at", "updated_at"}).
				AddRow(1, 1, 10, "running", now, nil, 40.7128, -74.0060, nil, nil, nil, nil, now, now))

		rental, err := repo.GetByID(1)

		assert.NoError(t, err)
		assert.NotNil(t, rental)
		assert.Equal(t, 1, rental.ID)
		assert.Equal(t, "running", rental.Status)
		assert.Nil(t, rental.DurationMinutes)
		assert.Nil(t, rental.Cost)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Rental found - ended status with cost", func(t *testing.T) {
		endTime := now.Add(30 * time.Minute)
		durationMinutes := 30
		cost := 15.0

		mock.ExpectQuery("SELECT id, user_id, bike_id, status").
			WithArgs(2).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "bike_id", "status", "start_time", "end_time", "start_latitude", "start_longitude", "end_latitude", "end_longitude", "duration_minutes", "cost", "created_at", "updated_at"}).
				AddRow(2, 1, 10, "ended", now, endTime, 40.7128, -74.0060, 40.7200, -74.0100, durationMinutes, cost, now, now))

		rental, err := repo.GetByID(2)

		assert.NoError(t, err)
		assert.NotNil(t, rental)
		assert.Equal(t, 2, rental.ID)
		assert.Equal(t, "ended", rental.Status)
		assert.Equal(t, endTime, rental.EndTime)
		assert.Equal(t, 40.7200, rental.EndLatitude)
		assert.Equal(t, -74.0100, rental.EndLongitude)
		assert.NotNil(t, rental.DurationMinutes)
		assert.Equal(t, durationMinutes, *rental.DurationMinutes)
		assert.NotNil(t, rental.Cost)
		assert.Equal(t, cost, *rental.Cost)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Rental not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, user_id, bike_id, status").
			WithArgs(999).
			WillReturnError(sql.ErrNoRows)

		rental, err := repo.GetByID(999)

		assert.Error(t, err)
		assert.Nil(t, rental)
		assert.Contains(t, err.Error(), "not found")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Database error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, user_id, bike_id, status").
			WithArgs(1).
			WillReturnError(fmt.Errorf("database error"))

		rental, err := repo.GetByID(1)

		assert.Error(t, err)
		assert.Nil(t, rental)
		assert.Contains(t, err.Error(), "error finding rental")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRentalRepository_GetActiveRentalsByUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRentalRepository(db)
	now := time.Now()

	t.Run("Successfully get rentals", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "user_id", "bike_id", "status", "start_time", "end_time", "start_latitude", "start_longitude", "end_latitude", "end_longitude", "duration_minutes", "cost", "created_at", "updated_at"}).
			AddRow(1, 1, 10, "running", now, nil, 40.7128, -74.0060, nil, nil, nil, nil, now, now).
			AddRow(2, 1, 11, "ended", now.Add(-1*time.Hour), now, 40.7128, -74.0060, 40.7200, -74.0100, 30, 15.0, now, now)

		mock.ExpectQuery("SELECT id, user_id, bike_id, status").
			WithArgs(1, 10, 0).
			WillReturnRows(rows)

		rentals, err := repo.GetActiveRentalsByUser(1, 1, 10)

		assert.NoError(t, err)
		assert.Len(t, rentals, 2)
		assert.Equal(t, 1, rentals[0].ID)
		assert.Equal(t, "running", rentals[0].Status)
		assert.Equal(t, 2, rentals[1].ID)
		assert.Equal(t, "ended", rentals[1].Status)
		assert.NotNil(t, rentals[1].DurationMinutes)
		assert.NotNil(t, rentals[1].Cost)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Empty result", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "user_id", "bike_id", "status", "start_time", "end_time", "start_latitude", "start_longitude", "end_latitude", "end_longitude", "duration_minutes", "cost", "created_at", "updated_at"})

		mock.ExpectQuery("SELECT id, user_id, bike_id, status").
			WithArgs(1, 10, 0).
			WillReturnRows(rows)

		rentals, err := repo.GetActiveRentalsByUser(1, 1, 10)

		assert.NoError(t, err)
		assert.Len(t, rentals, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Query error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, user_id, bike_id, status").
			WithArgs(1, 10, 0).
			WillReturnError(fmt.Errorf("database error"))

		rentals, err := repo.GetActiveRentalsByUser(1, 1, 10)

		assert.Error(t, err)
		assert.Nil(t, rentals)
		assert.Contains(t, err.Error(), "error querying rentals")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRentalRepository_CountByUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRentalRepository(db)

	t.Run("Successfully count rentals", func(t *testing.T) {
		mock.ExpectQuery("SELECT COUNT").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

		count, err := repo.CountByUser(1)

		assert.NoError(t, err)
		assert.Equal(t, 5, count)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Zero rentals", func(t *testing.T) {
		mock.ExpectQuery("SELECT COUNT").
			WithArgs(2).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		count, err := repo.CountByUser(2)

		assert.NoError(t, err)
		assert.Equal(t, 0, count)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Database error", func(t *testing.T) {
		mock.ExpectQuery("SELECT COUNT").
			WithArgs(1).
			WillReturnError(fmt.Errorf("database error"))

		count, err := repo.CountByUser(1)

		assert.Error(t, err)
		assert.Equal(t, 0, count)
		assert.Contains(t, err.Error(), "error counting rentals")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRentalRepository_GetActiveRentalByUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRentalRepository(db)
	now := time.Now()

	t.Run("Active rental found", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, user_id, bike_id, status").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "bike_id", "status", "start_time", "end_time", "start_latitude", "start_longitude", "end_latitude", "end_longitude", "duration_minutes", "cost", "created_at", "updated_at"}).
				AddRow(1, 1, 10, "running", now, nil, 40.7128, -74.0060, nil, nil, nil, nil, now, now))

		rental, err := repo.GetActiveRentalByUser(1)

		assert.NoError(t, err)
		assert.NotNil(t, rental)
		assert.Equal(t, 1, rental.ID)
		assert.Equal(t, "running", rental.Status)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("No active rental", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, user_id, bike_id, status").
			WithArgs(2).
			WillReturnError(sql.ErrNoRows)

		rental, err := repo.GetActiveRentalByUser(2)

		assert.NoError(t, err)
		assert.Nil(t, rental)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Database error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, user_id, bike_id, status").
			WithArgs(1).
			WillReturnError(fmt.Errorf("database error"))

		rental, err := repo.GetActiveRentalByUser(1)

		assert.Error(t, err)
		assert.Nil(t, rental)
		assert.Contains(t, err.Error(), "error finding active rental")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRentalRepository_EndRental(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRentalRepository(db)
	now := time.Now()

	t.Run("Successfully end rental", func(t *testing.T) {
		mock.ExpectExec("UPDATE rentals SET status = 'ended'").
			WithArgs(sqlmock.AnyArg(), 40.7200, -74.0100, 30, 15.0, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectQuery("SELECT id, user_id, bike_id, status").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "bike_id", "status", "start_time", "end_time", "start_latitude", "start_longitude", "end_latitude", "end_longitude", "duration_minutes", "cost", "created_at", "updated_at"}).
				AddRow(1, 1, 10, "ended", now, now, 40.7128, -74.0060, 40.7200, -74.0100, 30, 15.0, now, now))

		rental, err := repo.EndRental(1, 40.7200, -74.0100, 30, 15.0)

		assert.NoError(t, err)
		assert.NotNil(t, rental)
		assert.Equal(t, 1, rental.ID)
		assert.Equal(t, "ended", rental.Status)
		assert.Equal(t, 40.7200, rental.EndLatitude)
		assert.Equal(t, -74.0100, rental.EndLongitude)
		assert.NotNil(t, rental.DurationMinutes)
		assert.Equal(t, 30, *rental.DurationMinutes)
		assert.NotNil(t, rental.Cost)
		assert.Equal(t, 15.0, *rental.Cost)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Update error", func(t *testing.T) {
		mock.ExpectExec("UPDATE rentals SET status = 'ended'").
			WithArgs(sqlmock.AnyArg(), 40.7200, -74.0100, 30, 15.0, 1).
			WillReturnError(fmt.Errorf("database error"))

		rental, err := repo.EndRental(1, 40.7200, -74.0100, 30, 15.0)

		assert.Error(t, err)
		assert.Nil(t, rental)
		assert.Contains(t, err.Error(), "error ending rental")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}