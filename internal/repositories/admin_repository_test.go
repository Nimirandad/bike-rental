package repositories

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestAdminRepository_CreateBike(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewAdminRepository(db)
	now := time.Now()

	t.Run("Successfully create bike", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO bikes").
			WithArgs(1, 40.7128, -74.0060, 0.5).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery("SELECT id, is_available, latitude, longitude, price_per_minute, created_at, updated_at FROM bikes WHERE id = ?").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "is_available", "latitude", "longitude", "price_per_minute", "created_at", "updated_at"}).
				AddRow(1, 1, 40.7128, -74.0060, 0.5, now, now))

		bike, err := repo.CreateBike(40.7128, -74.0060, 0.5)

		assert.NoError(t, err)
		assert.NotNil(t, bike)
		assert.Equal(t, 1, bike.ID)
		assert.True(t, bike.IsAvailable)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAdminRepository_UpdateBike(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewAdminRepository(db)
	now := time.Now()

	newLat := 40.7200
	newPrice := 0.7
	available := false

	t.Run("Update multiple fields", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, is_available, latitude, longitude, price_per_minute, created_at, updated_at FROM bikes WHERE id = ?").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "is_available", "latitude", "longitude", "price_per_minute", "created_at", "updated_at"}).
				AddRow(1, 1, 40.7128, -74.0060, 0.5, now, now))

		mock.ExpectExec("UPDATE bikes SET").
			WithArgs(newLat, 0, newPrice, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectQuery("SELECT id, is_available, latitude, longitude, price_per_minute, created_at, updated_at FROM bikes WHERE id = ?").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "is_available", "latitude", "longitude", "price_per_minute", "created_at", "updated_at"}).
				AddRow(1, 0, newLat, -74.0060, newPrice, now, now))

		bike, err := repo.UpdateBike(1, &newLat, nil, &available, &newPrice)

		assert.NoError(t, err)
		assert.NotNil(t, bike)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("No fields to update", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, is_available, latitude, longitude, price_per_minute, created_at, updated_at FROM bikes WHERE id = ?").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "is_available", "latitude", "longitude", "price_per_minute", "created_at", "updated_at"}).
				AddRow(1, 1, 40.7128, -74.0060, 0.5, now, now))

		bike, err := repo.UpdateBike(1, nil, nil, nil, nil)

		assert.Error(t, err)
		assert.Nil(t, bike)
		assert.Contains(t, err.Error(), "no fields to update")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAdminRepository_GetAllBikes(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewAdminRepository(db)
	now := time.Now()

	t.Run("Get bikes with pagination", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "is_available", "latitude", "longitude", "price_per_minute", "created_at", "updated_at"}).
			AddRow(1, 1, 40.7128, -74.0060, 0.5, now, now).
			AddRow(2, 0, 40.7138, -74.0070, 0.6, now, now)

		mock.ExpectQuery("SELECT id, is_available, latitude, longitude, price_per_minute, created_at, updated_at FROM bikes ORDER BY id ASC LIMIT \\? OFFSET \\?").
			WithArgs(10, 0).
			WillReturnRows(rows)

		bikes, err := repo.GetAllBikes(1, 10)

		assert.NoError(t, err)
		assert.Len(t, bikes, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAdminRepository_GetAllUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewAdminRepository(db)
	now := time.Now()

	t.Run("Get users with pagination", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "created_at"}).
			AddRow(1, "test1@example.com", "John", "Doe", now).
			AddRow(2, "test2@example.com", "Jane", "Smith", now)

		mock.ExpectQuery("SELECT id, email, first_name, last_name, created_at FROM users ORDER BY id ASC LIMIT \\? OFFSET \\?").
			WithArgs(10, 0).
			WillReturnRows(rows)

		users, err := repo.GetAllUsers(1, 10)

		assert.NoError(t, err)
		assert.Len(t, users, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAdminRepository_UpdateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewAdminRepository(db)
	now := time.Now()

	newEmail := "newemail@example.com"
	newFirstName := "Jane"

	t.Run("Update user fields", func(t *testing.T) {
		mock.ExpectExec("UPDATE users SET email = \\?, first_name = \\?, updated_at = CURRENT_TIMESTAMP WHERE id = \\?").
			WithArgs(newEmail, newFirstName, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectQuery("SELECT id, email, first_name, last_name, created_at FROM users WHERE id = ?").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "created_at"}).
				AddRow(1, newEmail, newFirstName, "Doe", now))

		user, err := repo.UpdateUser(1, &newEmail, &newFirstName, nil, nil)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, newEmail, user.Email)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAdminRepository_GetAllRentals(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewAdminRepository(db)
	now := time.Now()

	t.Run("Get rentals with pagination", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "user_id", "bike_id", "status", "start_time", "end_time", "start_latitude", "start_longitude", "end_latitude", "end_longitude", "duration_minutes", "cost", "created_at", "updated_at"}).
			AddRow(1, 1, 10, "running", now, nil, 40.7128, -74.0060, nil, nil, nil, nil, now, now).
			AddRow(2, 2, 11, "ended", now, now, 40.7128, -74.0060, 40.7200, -74.0100, 30, 15.0, now, now)

		mock.ExpectQuery("SELECT id, user_id, bike_id, status").
			WithArgs(10, 0).
			WillReturnRows(rows)

		rentals, err := repo.GetAllRentals(1, 10)

		assert.NoError(t, err)
		assert.Len(t, rentals, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAdminRepository_UpdateRental(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewAdminRepository(db)
	now := time.Now()

	newStatus := "ended"

	t.Run("Update rental status", func(t *testing.T) {
		mock.ExpectExec("UPDATE rentals SET status = \\?, updated_at = CURRENT_TIMESTAMP WHERE id = \\?").
			WithArgs(newStatus, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectQuery("SELECT id, user_id, bike_id, status").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "bike_id", "status", "start_time", "end_time", "start_latitude", "start_longitude", "end_latitude", "end_longitude", "duration_minutes", "cost", "created_at", "updated_at"}).
				AddRow(1, 1, 10, newStatus, now, now, 40.7128, -74.0060, 40.7200, -74.0100, 30, 15.0, now, now))

		rental, err := repo.UpdateRental(1, &newStatus)

		assert.NoError(t, err)
		assert.NotNil(t, rental)
		assert.Equal(t, newStatus, rental.Status)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAdminRepository_CountMethods(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewAdminRepository(db)

	t.Run("CountAll bikes", func(t *testing.T) {
		mock.ExpectQuery("SELECT COUNT").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(10))

		count, err := repo.CountAll()

		assert.NoError(t, err)
		assert.Equal(t, 10, count)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("CountAllUsers", func(t *testing.T) {
		mock.ExpectQuery("SELECT COUNT").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

		count, err := repo.CountAllUsers()

		assert.NoError(t, err)
		assert.Equal(t, 5, count)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("CountAllRentals", func(t *testing.T) {
		mock.ExpectQuery("SELECT COUNT").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(20))

		count, err := repo.CountAllRentals()

		assert.NoError(t, err)
		assert.Equal(t, 20, count)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAdminRepository_EmailExistsByOtherUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewAdminRepository(db)

	t.Run("Email exists by another user", func(t *testing.T) {
		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM users WHERE email = \\? AND id != \\?\\)").
			WithArgs("test@example.com", 1).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(1))

		exists, err := repo.EmailExistsByOtherUser("test@example.com", 1)

		assert.NoError(t, err)
		assert.True(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Email does not exist", func(t *testing.T) {
		mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM users WHERE email = \\? AND id != \\?\\)").
			WithArgs("newuser@example.com", 2).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(0))

		exists, err := repo.EmailExistsByOtherUser("newuser@example.com", 2)

		assert.NoError(t, err)
		assert.False(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
