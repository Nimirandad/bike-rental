package repositories

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Nimirandad/bike-rental-service/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(email, hashedPassword, firstName, lastName string) (*models.User, error) {
	result, err := r.db.Exec(
		"INSERT INTO users (email, hashed_password, first_name, last_name) VALUES (?, ?, ?, ?)",
		email, hashedPassword, firstName, lastName,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	userID, _ := result.LastInsertId()
	return r.GetByID(int(userID))
}

func (r *UserRepository) GetByID(userID int) (*models.User, error) {
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

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User

	err := r.db.QueryRow(
		"SELECT id, email, first_name, last_name, created_at FROM users WHERE email = ?",
		email,
	).Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user with email %s not found", email)
	}
	if err != nil {
		return nil, fmt.Errorf("error finding user by email: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) EmailExists(email string) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking email existence: %w", err)
	}
	return exists, nil
}

func (r *UserRepository) GetPasswordHashByEmail(email string) (string, *models.User, error) {
	var user models.User
	var hashedPassword string

	err := r.db.QueryRow(
		"SELECT id, email, hashed_password, first_name, last_name, created_at FROM users WHERE email = ?",
		email,
	).Scan(&user.ID, &user.Email, &hashedPassword, &user.FirstName, &user.LastName, &user.CreatedAt)

	if err == sql.ErrNoRows {
		return "", nil, fmt.Errorf("user with email %s not found", email)
	}
	if err != nil {
		return "", nil, fmt.Errorf("error finding user credentials: %w", err)
	}

	return hashedPassword, &user, nil
}

func (r *UserRepository) EmailExistsByOtherUser(email string, userID int) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ? AND id != ?)", email, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking email existence: %w", err)
	}
	return exists, nil
}

func (r *UserRepository) Update(userID int, email, firstName, lastName *string) (*models.User, error) {
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

	if len(updates) == 0 {
		return r.GetByID(userID)
	}

	updates = append(updates, "updated_at = CURRENT_TIMESTAMP")

	query += strings.Join(updates, ", ") + " WHERE id = ?"
	args = append(args, userID)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	return r.GetByID(userID)
}