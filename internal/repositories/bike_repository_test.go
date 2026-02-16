package repositories

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestBikeRepository_CountAvailable(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewBikeRepository(db)

	t.Run("Successfully count available bikes", func(t *testing.T) {
		mock.ExpectQuery("SELECT COUNT").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

		count, err := repo.CountAvailable()

		assert.NoError(t, err)
		assert.Equal(t, 5, count)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Zero available bikes", func(t *testing.T) {
		mock.ExpectQuery("SELECT COUNT").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		count, err := repo.CountAvailable()

		assert.NoError(t, err)
		assert.Equal(t, 0, count)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Database error", func(t *testing.T) {
		mock.ExpectQuery("SELECT COUNT").
			WillReturnError(fmt.Errorf("database error"))

		count, err := repo.CountAvailable()

		assert.Error(t, err)
		assert.Equal(t, 0, count)
		assert.Contains(t, err.Error(), "error counting available bikes")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestBikeRepository_GetAvailable(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewBikeRepository(db)
	now := time.Now()

	t.Run("Successfully get available bikes - first page", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "is_available", "latitude", "longitude", "price_per_minute", "created_at", "updated_at"}).
			AddRow(1, 1, 40.7128, -74.0060, 0.5, now, now).
			AddRow(2, 1, 40.7138, -74.0070, 0.6, now, now)

		mock.ExpectQuery("SELECT id, is_available, latitude, longitude, price_per_minute, created_at, updated_at FROM bikes WHERE is_available = 1 LIMIT \\? OFFSET \\?").
			WithArgs(10, 0).
			WillReturnRows(rows)

		bikes, err := repo.GetAvailable(1, 10)

		assert.NoError(t, err)
		assert.Len(t, bikes, 2)
		assert.Equal(t, 1, bikes[0].ID)
		assert.True(t, bikes[0].IsAvailable)
		assert.Equal(t, 40.7128, bikes[0].Latitude)
		assert.Equal(t, -74.0060, bikes[0].Longitude)
		assert.Equal(t, 0.5, bikes[0].PricePerMinute)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Successfully get available bikes - second page", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "is_available", "latitude", "longitude", "price_per_minute", "created_at", "updated_at"}).
			AddRow(11, 1, 40.7200, -74.0100, 0.7, now, now)

		mock.ExpectQuery("SELECT id, is_available, latitude, longitude, price_per_minute, created_at, updated_at FROM bikes WHERE is_available = 1 LIMIT \\? OFFSET \\?").
			WithArgs(10, 10).
			WillReturnRows(rows)

		bikes, err := repo.GetAvailable(2, 10)

		assert.NoError(t, err)
		assert.Len(t, bikes, 1)
		assert.Equal(t, 11, bikes[0].ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Empty result", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "is_available", "latitude", "longitude", "price_per_minute", "created_at", "updated_at"})

		mock.ExpectQuery("SELECT id, is_available, latitude, longitude, price_per_minute, created_at, updated_at FROM bikes WHERE is_available = 1 LIMIT \\? OFFSET \\?").
			WithArgs(10, 0).
			WillReturnRows(rows)

		bikes, err := repo.GetAvailable(1, 10)

		assert.NoError(t, err)
		assert.Len(t, bikes, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Query error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, is_available, latitude, longitude, price_per_minute, created_at, updated_at FROM bikes WHERE is_available = 1 LIMIT \\? OFFSET \\?").
			WithArgs(10, 0).
			WillReturnError(fmt.Errorf("database error"))

		bikes, err := repo.GetAvailable(1, 10)

		assert.Error(t, err)
		assert.Nil(t, bikes)
		assert.Contains(t, err.Error(), "error querying available bikes")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Scan error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "is_available", "latitude", "longitude", "price_per_minute", "created_at", "updated_at"}).
			AddRow(1, "invalid", 40.7128, -74.0060, 0.5, now, now)

		mock.ExpectQuery("SELECT id, is_available, latitude, longitude, price_per_minute, created_at, updated_at FROM bikes WHERE is_available = 1 LIMIT \\? OFFSET \\?").
			WithArgs(10, 0).
			WillReturnRows(rows)

		bikes, err := repo.GetAvailable(1, 10)

		assert.Error(t, err)
		assert.Nil(t, bikes)
		assert.Contains(t, err.Error(), "error scanning bike")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestBikeRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewBikeRepository(db)
	now := time.Now()

	t.Run("Bike found - available", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, is_available, latitude, longitude, price_per_minute, created_at, updated_at FROM bikes WHERE id = ?").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "is_available", "latitude", "longitude", "price_per_minute", "created_at", "updated_at"}).
				AddRow(1, 1, 40.7128, -74.0060, 0.5, now, now))

		bike, err := repo.GetByID(1)

		assert.NoError(t, err)
		assert.NotNil(t, bike)
		assert.Equal(t, 1, bike.ID)
		assert.True(t, bike.IsAvailable)
		assert.Equal(t, 40.7128, bike.Latitude)
		assert.Equal(t, -74.0060, bike.Longitude)
		assert.Equal(t, 0.5, bike.PricePerMinute)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Bike found - not available", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, is_available, latitude, longitude, price_per_minute, created_at, updated_at FROM bikes WHERE id = ?").
			WithArgs(2).
			WillReturnRows(sqlmock.NewRows([]string{"id", "is_available", "latitude", "longitude", "price_per_minute", "created_at", "updated_at"}).
				AddRow(2, 0, 40.7138, -74.0070, 0.6, now, now))

		bike, err := repo.GetByID(2)

		assert.NoError(t, err)
		assert.NotNil(t, bike)
		assert.Equal(t, 2, bike.ID)
		assert.False(t, bike.IsAvailable)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Bike not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, is_available, latitude, longitude, price_per_minute, created_at, updated_at FROM bikes WHERE id = ?").
			WithArgs(999).
			WillReturnError(sql.ErrNoRows)

		bike, err := repo.GetByID(999)

		assert.Error(t, err)
		assert.Nil(t, bike)
		assert.Contains(t, err.Error(), "not found")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Database error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, is_available, latitude, longitude, price_per_minute, created_at, updated_at FROM bikes WHERE id = ?").
			WithArgs(1).
			WillReturnError(fmt.Errorf("database error"))

		bike, err := repo.GetByID(1)

		assert.Error(t, err)
		assert.Nil(t, bike)
		assert.Contains(t, err.Error(), "error finding bike")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestBikeRepository_UpdateAvailability(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewBikeRepository(db)

	t.Run("Update to available", func(t *testing.T) {
		mock.ExpectExec("UPDATE bikes SET is_available = \\?, updated_at = CURRENT_TIMESTAMP WHERE id = \\?").
			WithArgs(1, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.UpdateAvailability(1, true)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Update to not available", func(t *testing.T) {
		mock.ExpectExec("UPDATE bikes SET is_available = \\?, updated_at = CURRENT_TIMESTAMP WHERE id = \\?").
			WithArgs(0, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.UpdateAvailability(1, false)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Database error", func(t *testing.T) {
		mock.ExpectExec("UPDATE bikes SET is_available = \\?, updated_at = CURRENT_TIMESTAMP WHERE id = \\?").
			WithArgs(1, 1).
			WillReturnError(fmt.Errorf("database error"))

		err := repo.UpdateAvailability(1, true)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error updating bike availability")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}