package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// Connect establishes a connection to the PostgreSQL database
func Connect() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Set default values if environment variables are not set
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5432"
	}
	if user == "" {
		user = "postgres"
	}
	if dbname == "" {
		dbname = "event_system"
	}

	// Create connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Check connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("Successfully connected to database")
	return db, nil
}

// CreateTables creates all necessary tables if they don't exist
func CreateTables(db *sql.DB) error {
	// Create users table
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	// Create events table
	eventsTable := `
	CREATE TABLE IF NOT EXISTS events (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		description TEXT,
		location VARCHAR(255),
		start_time TIMESTAMP NOT NULL,
		end_time TIMESTAMP NOT NULL,
		capacity INTEGER NOT NULL,
		organizer_id INTEGER NOT NULL REFERENCES users(id),
		status VARCHAR(20) DEFAULT 'open',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT check_capacity CHECK (capacity > 0),
		CONSTRAINT check_time CHECK (end_time > start_time)
	);
	`

	// Create participants table (for many-to-many relationship between users and events)
	participantsTable := `
	CREATE TABLE IF NOT EXISTS participants (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL REFERENCES users(id),
		event_id INTEGER NOT NULL REFERENCES events(id),
		joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT unique_participant UNIQUE (user_id, event_id)
	);
	`

	// Execute SQL statements
	_, err := db.Exec(usersTable)
	if err != nil {
		return err
	}

	_, err = db.Exec(eventsTable)
	if err != nil {
		return err
	}

	_, err = db.Exec(participantsTable)
	if err != nil {
		return err
	}

	log.Println("Database tables created successfully")
	return nil
}