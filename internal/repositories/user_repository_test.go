package repositories

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)

	t.Run("Successful user creation", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO users").
			WithArgs("test@example.com", "hashedpwd", "John", "Doe").
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery("SELECT id, email, first_name, last_name, created_at FROM users WHERE id = ?").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "created_at"}).
				AddRow(1, "test@example.com", "John", "Doe", time.Now()))

		user, err := repo.Create("test@example.com", "hashedpwd", "John", "Doe")

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, 1, user.ID)
		assert.Equal(t, "test@example.com", user.Email)
		assert.Equal(t, "John", user.FirstName)
		assert.Equal(t, "Doe", user.LastName)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Database error on insert", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO users").
			WithArgs("test@example.com", "hashedpwd", "John", "Doe").
			WillReturnError(fmt.Errorf("database error"))

		user, err := repo.Create("test@example.com", "hashedpwd", "John", "Doe")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "error creating user")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)
	now := time.Now()

	t.Run("User found", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, email, first_name, last_name, created_at FROM users WHERE id = ?").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "created_at"}).
				AddRow(1, "test@example.com", "John", "Doe", now))

		user, err := repo.GetByID(1)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, 1, user.ID)
		assert.Equal(t, "test@example.com", user.Email)
		assert.Equal(t, "John", user.FirstName)
		assert.Equal(t, "Doe", user.LastName)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("User not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, email, first_name, last_name, created_at FROM users WHERE id = ?").
			WithArgs(999).
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetByID(999)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "not found")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Database error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, email, first_name, last_name, created_at FROM users WHERE id = ?").
			WithArgs(1).
			WillReturnError(fmt.Errorf("database error"))

		user, err := repo.GetByID(1)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "error finding user")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)
	now := time.Now()

	t.Run("User found", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, email, first_name, last_name, created_at FROM users WHERE email = ?").
			WithArgs("test@example.com").
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "created_at"}).
				AddRow(1, "test@example.com", "John", "Doe", now))

		user, err := repo.GetByEmail("test@example.com")

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, 1, user.ID)
		assert.Equal(t, "test@example.com", user.Email)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("User not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, email, first_name, last_name, created_at FROM users WHERE email = ?").
			WithArgs("notfound@example.com").
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetByEmail("notfound@example.com")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "not found")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Database error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, email, first_name, last_name, created_at FROM users WHERE email = ?").
			WithArgs("test@example.com").
			WillReturnError(fmt.Errorf("database error"))

		user, err := repo.GetByEmail("test@example.com")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "error finding user by email")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepository_EmailExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)

	t.Run("Email exists", func(t *testing.T) {
		mock.ExpectQuery("SELECT EXISTS").
			WithArgs("test@example.com").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		exists, err := repo.EmailExists("test@example.com")

		assert.NoError(t, err)
		assert.True(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Email does not exist", func(t *testing.T) {
		mock.ExpectQuery("SELECT EXISTS").
			WithArgs("notfound@example.com").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		exists, err := repo.EmailExists("notfound@example.com")

		assert.NoError(t, err)
		assert.False(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Database error", func(t *testing.T) {
		mock.ExpectQuery("SELECT EXISTS").
			WithArgs("test@example.com").
			WillReturnError(fmt.Errorf("database error"))

		exists, err := repo.EmailExists("test@example.com")

		assert.Error(t, err)
		assert.False(t, exists)
		assert.Contains(t, err.Error(), "error checking email existence")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepository_GetPasswordHashByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)
	now := time.Now()

	t.Run("User found with password", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, email, hashed_password, first_name, last_name, created_at FROM users WHERE email = ?").
			WithArgs("test@example.com").
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "hashed_password", "first_name", "last_name", "created_at"}).
				AddRow(1, "test@example.com", "$2a$10$hashedpassword", "John", "Doe", now))

		hash, user, err := repo.GetPasswordHashByEmail("test@example.com")

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "$2a$10$hashedpassword", hash)
		assert.Equal(t, 1, user.ID)
		assert.Equal(t, "test@example.com", user.Email)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("User not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, email, hashed_password, first_name, last_name, created_at FROM users WHERE email = ?").
			WithArgs("notfound@example.com").
			WillReturnError(sql.ErrNoRows)

		hash, user, err := repo.GetPasswordHashByEmail("notfound@example.com")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Empty(t, hash)
		assert.Contains(t, err.Error(), "not found")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Database error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, email, hashed_password, first_name, last_name, created_at FROM users WHERE email = ?").
			WithArgs("test@example.com").
			WillReturnError(fmt.Errorf("database error"))

		hash, user, err := repo.GetPasswordHashByEmail("test@example.com")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Empty(t, hash)
		assert.Contains(t, err.Error(), "error finding user credentials")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepository_EmailExistsByOtherUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)

	t.Run("Email exists for other user", func(t *testing.T) {
		mock.ExpectQuery("SELECT EXISTS").
			WithArgs("test@example.com", 1).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		exists, err := repo.EmailExistsByOtherUser("test@example.com", 1)

		assert.NoError(t, err)
		assert.True(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Email does not exist for other user", func(t *testing.T) {
		mock.ExpectQuery("SELECT EXISTS").
			WithArgs("test@example.com", 1).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		exists, err := repo.EmailExistsByOtherUser("test@example.com", 1)

		assert.NoError(t, err)
		assert.False(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Database error", func(t *testing.T) {
		mock.ExpectQuery("SELECT EXISTS").
			WithArgs("test@example.com", 1).
			WillReturnError(fmt.Errorf("database error"))

		exists, err := repo.EmailExistsByOtherUser("test@example.com", 1)

		assert.Error(t, err)
		assert.False(t, exists)
		assert.Contains(t, err.Error(), "error checking email existence")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)
	now := time.Now()

	t.Run("Update all fields", func(t *testing.T) {
		newEmail := "newemail@example.com"
		newFirstName := "Jane"
		newLastName := "Smith"

		mock.ExpectExec("UPDATE users SET email = \\?, first_name = \\?, last_name = \\?, updated_at = CURRENT_TIMESTAMP WHERE id = \\?").
			WithArgs(newEmail, newFirstName, newLastName, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectQuery("SELECT id, email, first_name, last_name, created_at FROM users WHERE id = ?").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "created_at"}).
				AddRow(1, newEmail, newFirstName, newLastName, now))

		user, err := repo.Update(1, &newEmail, &newFirstName, &newLastName)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, newEmail, user.Email)
		assert.Equal(t, newFirstName, user.FirstName)
		assert.Equal(t, newLastName, user.LastName)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Update only email", func(t *testing.T) {
		newEmail := "newemail@example.com"

		mock.ExpectExec("UPDATE users SET email = \\?, updated_at = CURRENT_TIMESTAMP WHERE id = \\?").
			WithArgs(newEmail, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectQuery("SELECT id, email, first_name, last_name, created_at FROM users WHERE id = ?").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "created_at"}).
				AddRow(1, newEmail, "John", "Doe", now))

		user, err := repo.Update(1, &newEmail, nil, nil)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, newEmail, user.Email)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("No fields to update", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, email, first_name, last_name, created_at FROM users WHERE id = ?").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "created_at"}).
				AddRow(1, "test@example.com", "John", "Doe", now))

		user, err := repo.Update(1, nil, nil, nil)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Update error", func(t *testing.T) {
		newEmail := "newemail@example.com"

		mock.ExpectExec("UPDATE users SET email = \\?, updated_at = CURRENT_TIMESTAMP WHERE id = \\?").
			WithArgs(newEmail, 1).
			WillReturnError(fmt.Errorf("database error"))

		user, err := repo.Update(1, &newEmail, nil, nil)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "error updating user")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}