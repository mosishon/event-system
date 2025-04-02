package models

import "time"

// swagger:model
type ParticipantStatusResponse struct {
	IsParticipant bool `json:"is_participant"`
}

// swagger:model
type ParticipantCountResponse struct {
	Count int `json:"count"`
}

// swagger:model
type ParticipantResponse struct {
	UserID   int       `json:"user_id"`
	EventID  int       `json:"event_id"`
	JoinedAt time.Time `json:"joined_at"`
}
