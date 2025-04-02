package repositories

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/event-system/models"
)

// UserRepository handles database operations related to users
type UserRepository struct {
	DB *sql.DB
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// Create inserts a new user into the database
func (r *UserRepository) Create(user *models.User) error {
	query := `
	INSERT INTO users (username, email, password, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id
	`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	err := r.DB.QueryRow(
		query,
		user.Username,
		user.Email,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	if err != nil {
		log.Printf("Error creating user: %v", err)
		return err
	}

	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id int) (*models.User, error) {
	query := `
	SELECT id, username, email, password, created_at, updated_at
	FROM users
	WHERE id = $1
	`

	user := &models.User{}
	err := r.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		log.Printf("Error getting user by ID: %v", err)
		return nil, err
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := `
	SELECT id, username, email, password, created_at, updated_at
	FROM users
	WHERE email = $1
	`

	user := &models.User{}
	err := r.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		log.Printf("Error getting user by email: %v", err)
		return nil, err
	}

	return user, nil
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	query := `
	SELECT id, username, email, password, created_at, updated_at
	FROM users
	WHERE username = $1
	`

	user := &models.User{}
	err := r.DB.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		log.Printf("Error getting user by username: %v", err)
		return nil, err
	}

	return user, nil
}
