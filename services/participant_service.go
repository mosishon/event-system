package services

import (
	"github.com/event-system/repositories"
)

// ParticipantService handles participant related business logic
type ParticipantService struct {
	ParticipantRepo *repositories.ParticipantRepository
}

// NewParticipantService creates a new participant service instance
func NewParticipantService(participantRepo *repositories.ParticipantRepository) *ParticipantService {
	return &ParticipantService{ParticipantRepo: participantRepo}
}

// JoinEvent adds a user as a participant to an event
func (s *ParticipantService) JoinEvent(userID, eventID int) error {
	return s.ParticipantRepo.JoinEvent(userID, eventID)
}

// LeaveEvent removes a user as a participant from an event
func (s *ParticipantService) LeaveEvent(userID, eventID int) error {
	return s.ParticipantRepo.LeaveEvent(userID, eventID)
}

// IsParticipant checks if a user is a participant of an event
func (s *ParticipantService) IsParticipant(userID, eventID int) (bool, error) {
	return s.ParticipantRepo.IsParticipant(userID, eventID)
}

// GetParticipantCount returns the number of participants for an event
func (s *ParticipantService) GetParticipantCount(eventID int) (int, error) {
	return s.ParticipantRepo.GetParticipantCount(eventID)
}
