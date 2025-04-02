package repositories

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/event-system/models"
)

// EventRepository handles database operations related to events
type EventRepository struct {
	DB *sql.DB
}

// NewEventRepository creates a new event repository instance
func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{DB: db}
}

// Create inserts a new event into the database
func (r *EventRepository) Create(event *models.Event) error {
	query := `
	INSERT INTO events (name, description, location, start_time, end_time, capacity, organizer_id, status, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id
	`

	now := time.Now()
	event.CreatedAt = now
	event.UpdatedAt = now
	if event.Status == "" {
		event.Status = "open"
	}

	err := r.DB.QueryRow(
		query,
		event.Name,
		event.Description,
		event.Location,
		event.StartTime,
		event.EndTime,
		event.Capacity,
		event.OrganizerID,
		event.Status,
		event.CreatedAt,
		event.UpdatedAt,
	).Scan(&event.ID)

	if err != nil {
		log.Printf("Error creating event: %v", err)
		return err
	}

	return nil
}

// GetByID retrieves an event by ID
func (r *EventRepository) GetByID(id int) (*models.Event, error) {
	query := `
	SELECT id, name, description, location, start_time, end_time, capacity, organizer_id, status, created_at, updated_at
	FROM events
	WHERE id = $1
	`

	event := &models.Event{}
	err := r.DB.QueryRow(query, id).Scan(
		&event.ID,
		&event.Name,
		&event.Description,
		&event.Location,
		&event.StartTime,
		&event.EndTime,
		&event.Capacity,
		&event.OrganizerID,
		&event.Status,
		&event.CreatedAt,
		&event.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("event not found")
		}
		log.Printf("Error getting event by ID: %v", err)
		return nil, err
	}

	return event, nil
}

// Update updates an existing event
func (r *EventRepository) Update(event *models.Event) error {
	query := `
	UPDATE events
	SET name = $1, description = $2, location = $3, start_time = $4, end_time = $5, 
	    capacity = $6, status = $7, updated_at = $8
	WHERE id = $9 AND organizer_id = $10
	RETURNING id
	`

	event.UpdatedAt = time.Now()

	var id int
	err := r.DB.QueryRow(
		query,
		event.Name,
		event.Description,
		event.Location,
		event.StartTime,
		event.EndTime,
		event.Capacity,
		event.Status,
		event.UpdatedAt,
		event.ID,
		event.OrganizerID,
	).Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("event not found or you are not the organizer")
		}
		log.Printf("Error updating event: %v", err)
		return err
	}

	return nil
}

// Delete deletes an event by ID if it has no participants
func (r *EventRepository) Delete(id int, organizerID int) error {
	// First check if the event has participants
	var count int
	participantsQuery := `
	SELECT COUNT(*) FROM participants WHERE event_id = $1
	`
	err := r.DB.QueryRow(participantsQuery, id).Scan(&count)
	if err != nil {
		log.Printf("Error checking participants: %v", err)
		return err
	}

	if count > 0 {
		return errors.New("cannot delete event with participants")
	}

	// Delete the event
	query := `
	DELETE FROM events
	WHERE id = $1 AND organizer_id = $2
	RETURNING id
	`

	var deletedID int
	err = r.DB.QueryRow(query, id, organizerID).Scan(&deletedID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("event not found or you are not the organizer")
		}
		log.Printf("Error deleting event: %v", err)
		return err
	}

	return nil
}

// GetAllPublic retrieves all open events
func (r *EventRepository) GetAllPublic() ([]models.Event, error) {
	query := `
	SELECT id, name, description, location, start_time, end_time, capacity, organizer_id, status, created_at, updated_at
	FROM events
	WHERE status = 'open'
	ORDER BY start_time ASC
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		log.Printf("Error getting public events: %v", err)
		return nil, err
	}
	defer rows.Close()

	events := []models.Event{}
	for rows.Next() {
		event := models.Event{}
		err := rows.Scan(
			&event.ID,
			&event.Name,
			&event.Description,
			&event.Location,
			&event.StartTime,
			&event.EndTime,
			&event.Capacity,
			&event.OrganizerID,
			&event.Status,
			&event.CreatedAt,
			&event.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning event: %v", err)
			return nil, err
		}
		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating events: %v", err)
		return nil, err
	}

	return events, nil
}

// GetByOrganizer retrieves all events created by a specific organizer
func (r *EventRepository) GetByOrganizer(organizerID int) ([]models.Event, error) {
	query := `
	SELECT id, name, description, location, start_time, end_time, capacity, organizer_id, status, created_at, updated_at
	FROM events
	WHERE organizer_id = $1
	ORDER BY start_time ASC
	`

	rows, err := r.DB.Query(query, organizerID)
	if err != nil {
		log.Printf("Error getting events by organizer: %v", err)
		return nil, err
	}
	defer rows.Close()

	events := []models.Event{}
	for rows.Next() {
		event := models.Event{}
		err := rows.Scan(
			&event.ID,
			&event.Name,
			&event.Description,
			&event.Location,
			&event.StartTime,
			&event.EndTime,
			&event.Capacity,
			&event.OrganizerID,
			&event.Status,
			&event.CreatedAt,
			&event.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning event: %v", err)
			return nil, err
		}
		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating events: %v", err)
		return nil, err
	}

	return events, nil
}

// GetParticipantsByEventID retrieves all participants for a specific event
func (r *EventRepository) GetParticipantsByEventID(eventID int) ([]models.UserResponse, error) {
	query := `
	SELECT u.id, u.username, u.email, u.created_at
	FROM users u
	JOIN participants p ON u.id = p.user_id
	WHERE p.event_id = $1
	`

	rows, err := r.DB.Query(query, eventID)
	if err != nil {
		log.Printf("Error getting participants: %v", err)
		return nil, err
	}
	defer rows.Close()

	participants := []models.UserResponse{}
	for rows.Next() {
		participant := models.UserResponse{}
		err := rows.Scan(
			&participant.ID,
			&participant.Username,
			&participant.Email,
			&participant.CreatedAt,
		)
		if err != nil {
			log.Printf("Error scanning participant: %v", err)
			return nil, err
		}
		participants = append(participants, participant)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating participants: %v", err)
		return nil, err
	}

	return participants, nil
}

// GetEventsByParticipant retrieves all events a user is participating in
func (r *EventRepository) GetEventsByParticipant(userID int) ([]models.Event, error) {
	query := `
	SELECT e.id, e.name, e.description, e.location, e.start_time, e.end_time, e.capacity, e.organizer_id, e.status, e.created_at, e.updated_at
	FROM events e
	JOIN participants p ON e.id = p.event_id
	WHERE p.user_id = $1
	ORDER BY e.start_time ASC
	`

	rows, err := r.DB.Query(query, userID)
	if err != nil {
		log.Printf("Error getting events by participant: %v", err)
		return nil, err
	}
	defer rows.Close()

	events := []models.Event{}
	for rows.Next() {
		event := models.Event{}
		err := rows.Scan(
			&event.ID,
			&event.Name,
			&event.Description,
			&event.Location,
			&event.StartTime,
			&event.EndTime,
			&event.Capacity,
			&event.OrganizerID,
			&event.Status,
			&event.CreatedAt,
			&event.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning event: %v", err)
			return nil, err
		}
		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating events: %v", err)
		return nil, err
	}

	return events, nil
}
