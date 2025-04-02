package repositories

import (
	"database/sql"
	"errors"
	"log"
	"time"
)

// ParticipantRepository handles database operations related to event participants
type ParticipantRepository struct {
	DB *sql.DB
}

// NewParticipantRepository creates a new participant repository instance
func NewParticipantRepository(db *sql.DB) *ParticipantRepository {
	return &ParticipantRepository{DB: db}
}

// JoinEvent adds a user as a participant to an event
func (r *ParticipantRepository) JoinEvent(userID, eventID int) error {
	// First check if the event exists and is open
	eventQuery := `
	SELECT status,capacity FROM events WHERE id = $1
	`
	var status string
	var capacity int
	err := r.DB.QueryRow(eventQuery, eventID).Scan(&status, &capacity)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("event not found")
		}
		return err
	}

	if status != "open" {
		return errors.New("event is not open for registration")
	}

	// Check if the event is full
	countQuery := `
	SELECT COUNT(*) FROM participants WHERE event_id = $1
	`
	var count int
	err = r.DB.QueryRow(countQuery, eventID).Scan(&count)
	if err != nil {
		log.Printf("Error checking participant count: %v", err)
		return err
	}

	if count >= capacity {
		return errors.New("event is at full capacity")
	}

	// Check if user is already a participant
	checkQuery := `
	SELECT id FROM participants WHERE user_id = $1 AND event_id = $2
	`
	var participantID int
	err = r.DB.QueryRow(checkQuery, userID, eventID).Scan(&participantID)
	if err == nil {
		return errors.New("user is already a participant of this event")
	} else if err != sql.ErrNoRows {
		log.Printf("Error checking existing participant: %v", err)
		return err
	}

	// Check if user has reached the maximum number of active events
	// For this example, we'll set a limit of 5 active events per user
	activeEventsQuery := `
	SELECT COUNT(*) FROM participants p
	JOIN events e ON p.event_id = e.id
	WHERE p.user_id = $1 AND e.end_time > $2
	`
	var activeCount int
	err = r.DB.QueryRow(activeEventsQuery, userID, time.Now()).Scan(&activeCount)
	if err != nil {
		log.Printf("Error checking active events count: %v", err)
		return err
	}

	if activeCount >= 5 {
		return errors.New("user has reached the maximum number of active events")
	}

	// Add user as participant
	insertQuery := `
	INSERT INTO participants (user_id, event_id, joined_at)
	VALUES ($1, $2, $3)
	RETURNING id
	`

	var newID int
	err = r.DB.QueryRow(insertQuery, userID, eventID, time.Now()).Scan(&newID)
	if err != nil {
		log.Printf("Error adding participant: %v", err)
		return err
	}

	return nil
}

// LeaveEvent removes a user as a participant from an event
func (r *ParticipantRepository) LeaveEvent(userID, eventID int) error {
	// Check if the user is a participant
	checkQuery := `
	SELECT id FROM participants WHERE user_id = $1 AND event_id = $2
	`
	var participantID int
	err := r.DB.QueryRow(checkQuery, userID, eventID).Scan(&participantID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user is not a participant of this event")
		}
		log.Printf("Error checking participant: %v", err)
		return err
	}

	// Remove participant
	deleteQuery := `
	DELETE FROM participants
	WHERE user_id = $1 AND event_id = $2
	RETURNING id
	`

	var deletedID int
	err = r.DB.QueryRow(deleteQuery, userID, eventID).Scan(&deletedID)
	if err != nil {
		log.Printf("Error removing participant: %v", err)
		return err
	}

	return nil
}

// IsParticipant checks if a user is a participant of an event
func (r *ParticipantRepository) IsParticipant(userID, eventID int) (bool, error) {
	query := `
	SELECT id FROM participants WHERE user_id = $1 AND event_id = $2
	`

	var id int
	err := r.DB.QueryRow(query, userID, eventID).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		log.Printf("Error checking participant: %v", err)
		return false, err
	}

	return true, nil
}

// GetParticipantCount returns the number of participants for an event
func (r *ParticipantRepository) GetParticipantCount(eventID int) (int, error) {
	query := `
	SELECT COUNT(*) FROM participants WHERE event_id = $1
	`

	var count int
	err := r.DB.QueryRow(query, eventID).Scan(&count)
	if err != nil {
		log.Printf("Error getting participant count: %v", err)
		return 0, err
	}

	return count, nil
}
