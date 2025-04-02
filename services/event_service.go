package services

import (
	"errors"
	"log"

	"github.com/event-system/models"
	"github.com/event-system/repositories"
)

// EventService handles event related business logic
type EventService struct {
	EventRepo       *repositories.EventRepository
	ParticipantRepo *repositories.ParticipantRepository
}

// NewEventService creates a new event service instance
func NewEventService(eventRepo *repositories.EventRepository, participantRepo *repositories.ParticipantRepository) *EventService {
	return &EventService{
		EventRepo:       eventRepo,
		ParticipantRepo: participantRepo,
	}
}
func (s *EventService) CloseEvent(organizerID int, eventID int) (*models.EventResponse, error) {
	// Get existing event
	existingEvent, err := s.EventRepo.GetByID(eventID)
	if err != nil {
		return nil, err
	}
	// Check if user is the organizer
	if existingEvent.OrganizerID != organizerID {
		return nil, errors.New("you are not the organizer of this event")
	}
	// check if event is already closed
	if existingEvent.Status == "closed" {
		return nil, errors.New("event is already closed")
	}
	// Close the event
	existingEvent.Status = "closed"
	// Save updated event
	err = s.EventRepo.Update(existingEvent)
	if err != nil {
		log.Printf("Error closing event: %v", err)
		return nil, errors.New("error closing event")
	}
	// Return updated event
	return &models.EventResponse{
		ID:          existingEvent.ID,
		Name:        existingEvent.Name,
		Description: existingEvent.Description,
		Location:    existingEvent.Location,
		StartTime:   existingEvent.StartTime,
		EndTime:     existingEvent.EndTime,
		Capacity:    existingEvent.Capacity,
		OrganizerID: existingEvent.OrganizerID,
		Status:      existingEvent.Status,
		CreatedAt:   existingEvent.CreatedAt,
		UpdatedAt:   existingEvent.UpdatedAt,
	}, nil

}

func (s *EventService) OpenEvent(organizerID int, eventID int) (*models.EventResponse, error) {
	// Get existing event
	existingEvent, err := s.EventRepo.GetByID(eventID)
	if err != nil {
		return nil, err
	}
	// Check if user is the organizer
	if existingEvent.OrganizerID != organizerID {
		return nil, errors.New("you are not the organizer of this event")
	}
	// check if event is already closed
	if existingEvent.Status == "open" {
		return nil, errors.New("event is already open")
	}
	// Close the event
	existingEvent.Status = "open"
	// Save updated event
	err = s.EventRepo.Update(existingEvent)
	if err != nil {
		log.Printf("Error opening event: %v", err)
		return nil, errors.New("error opening event")
	}
	// Return updated event
	return &models.EventResponse{
		ID:          existingEvent.ID,
		Name:        existingEvent.Name,
		Description: existingEvent.Description,
		Location:    existingEvent.Location,
		StartTime:   existingEvent.StartTime,
		EndTime:     existingEvent.EndTime,
		Capacity:    existingEvent.Capacity,
		OrganizerID: existingEvent.OrganizerID,
		Status:      existingEvent.Status,
		CreatedAt:   existingEvent.CreatedAt,
		UpdatedAt:   existingEvent.UpdatedAt,
	}, nil

}

// CreateEvent creates a new event
func (s *EventService) CreateEvent(req models.EventRequest, organizerID int) (*models.EventResponse, error) {
	// Create event object
	event := &models.Event{
		Name:        req.Name,
		Description: req.Description,
		Location:    req.Location,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Capacity:    req.Capacity,
		OrganizerID: organizerID,
		Status:      "open",
	}

	// Save event to database
	err := s.EventRepo.Create(event)
	if err != nil {
		log.Printf("Error creating event: %v", err)
		return nil, errors.New("error creating event")
	}

	// اطلاعات رویداد رو برمیگردونه
	return &models.EventResponse{
		ID:          event.ID,
		Name:        event.Name,
		Description: event.Description,
		Location:    event.Location,
		StartTime:   event.StartTime,
		EndTime:     event.EndTime,
		Capacity:    event.Capacity,
		OrganizerID: event.OrganizerID,
		Status:      event.Status,
		CreatedAt:   event.CreatedAt,
		UpdatedAt:   event.UpdatedAt,
	}, nil
}

// GetEventByID retrieves an event by ID
func (s *EventService) GetEventByID(id int) (*models.EventResponse, error) {
	// Get event from database
	event, err := s.EventRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// اطلاعات رویداد رو برمیگردونه
	return &models.EventResponse{
		ID:          event.ID,
		Name:        event.Name,
		Description: event.Description,
		Location:    event.Location,
		StartTime:   event.StartTime,
		EndTime:     event.EndTime,
		Capacity:    event.Capacity,
		OrganizerID: event.OrganizerID,
		Status:      event.Status,
		CreatedAt:   event.CreatedAt,
		UpdatedAt:   event.UpdatedAt,
	}, nil
}

// UpdateEvent updates an existing event
func (s *EventService) UpdateEvent(id int, req models.EventRequest, organizerID int) (*models.EventResponse, error) {
	// Get existing event
	existingEvent, err := s.EventRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check if user is the organizer
	if existingEvent.OrganizerID != organizerID {
		return nil, errors.New("you are not the organizer of this event")
	}

	// Update event fields
	existingEvent.Name = req.Name
	existingEvent.Description = req.Description
	existingEvent.Location = req.Location
	existingEvent.StartTime = req.StartTime
	existingEvent.EndTime = req.EndTime
	existingEvent.Capacity = req.Capacity

	// Save updated event
	err = s.EventRepo.Update(existingEvent)
	if err != nil {
		log.Printf("Error updating event: %v", err)
		return nil, errors.New("error updating event")
	}

	// Return updated event
	return &models.EventResponse{
		ID:          existingEvent.ID,
		Name:        existingEvent.Name,
		Description: existingEvent.Description,
		Location:    existingEvent.Location,
		StartTime:   existingEvent.StartTime,
		EndTime:     existingEvent.EndTime,
		Capacity:    existingEvent.Capacity,
		OrganizerID: existingEvent.OrganizerID,
		Status:      existingEvent.Status,
		CreatedAt:   existingEvent.CreatedAt,
		UpdatedAt:   existingEvent.UpdatedAt,
	}, nil
}

// DeleteEvent deletes an event
func (s *EventService) DeleteEvent(id int, organizerID int) error {
	// Delete event from database
	return s.EventRepo.Delete(id, organizerID)
}

// GetAllPublicEvents retrieves all public events
func (s *EventService) GetAllPublicEvents() ([]models.EventResponse, error) {
	// Get events from database
	events, err := s.EventRepo.GetAllPublic()
	if err != nil {
		return nil, err
	}

	// Convert to response format
	response := make([]models.EventResponse, len(events))
	for i, event := range events {
		response[i] = models.EventResponse{
			ID:          event.ID,
			Name:        event.Name,
			Description: event.Description,
			Location:    event.Location,
			StartTime:   event.StartTime,
			EndTime:     event.EndTime,
			Capacity:    event.Capacity,
			OrganizerID: event.OrganizerID,
			Status:      event.Status,
			CreatedAt:   event.CreatedAt,
			UpdatedAt:   event.UpdatedAt,
		}
	}

	return response, nil
}

// GetEventsByOrganizer retrieves all events created by a specific organizer
func (s *EventService) GetEventsByOrganizer(organizerID int) ([]models.EventResponse, error) {
	// Get events from database
	events, err := s.EventRepo.GetByOrganizer(organizerID)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	response := make([]models.EventResponse, len(events))
	for i, event := range events {
		response[i] = models.EventResponse{
			ID:          event.ID,
			Name:        event.Name,
			Description: event.Description,
			Location:    event.Location,
			StartTime:   event.StartTime,
			EndTime:     event.EndTime,
			Capacity:    event.Capacity,
			OrganizerID: event.OrganizerID,
			Status:      event.Status,
			CreatedAt:   event.CreatedAt,
			UpdatedAt:   event.UpdatedAt,
		}
	}

	return response, nil
}

// GetEventWithParticipants retrieves an event with its participants
func (s *EventService) GetEventWithParticipants(id int, organizerID int) (*models.EventWithParticipantsResponse, error) {
	// Get event from database
	event, err := s.EventRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check if user is the organizer
	if event.OrganizerID != organizerID {
		return nil, errors.New("you are not the organizer of this event")
	}

	// Get participants
	participants, err := s.EventRepo.GetParticipantsByEventID(id)
	if err != nil {
		return nil, err
	}

	// Create response
	response := &models.EventWithParticipantsResponse{
		Event: models.EventResponse{
			ID:          event.ID,
			Name:        event.Name,
			Description: event.Description,
			Location:    event.Location,
			StartTime:   event.StartTime,
			EndTime:     event.EndTime,
			Capacity:    event.Capacity,
			OrganizerID: event.OrganizerID,
			Status:      event.Status,
			CreatedAt:   event.CreatedAt,
			UpdatedAt:   event.UpdatedAt,
		},
		Participants: participants,
	}

	return response, nil
}

// GetEventsByParticipant retrieves all events a user is participating in
func (s *EventService) GetEventsByParticipant(userID int) ([]models.EventResponse, error) {
	// Get events from database
	events, err := s.EventRepo.GetEventsByParticipant(userID)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	response := make([]models.EventResponse, len(events))
	for i, event := range events {
		response[i] = models.EventResponse{
			ID:          event.ID,
			Name:        event.Name,
			Description: event.Description,
			Location:    event.Location,
			StartTime:   event.StartTime,
			EndTime:     event.EndTime,
			Capacity:    event.Capacity,
			OrganizerID: event.OrganizerID,
			Status:      event.Status,
			CreatedAt:   event.CreatedAt,
			UpdatedAt:   event.UpdatedAt,
		}
	}

	return response, nil
}
