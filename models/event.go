package models

import "time"

// رویداد رو تو سیستم نشون میده
type Event struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Capacity    int       `json:"capacity"`
	OrganizerID int       `json:"organizer_id"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ساختار درخواست ساخت/آپدیت رویداد
type EventRequest struct {
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	StartTime   time.Time `json:"start_time" validate:"required"`
	EndTime     time.Time `json:"end_time" validate:"required,gtfield=StartTime"`
	Capacity    int       `json:"capacity" validate:"required,gt=0"`
}

// ساختار پاسخ رویداد
type EventResponse struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Capacity    int       `json:"capacity"`
	OrganizerID int       `json:"organizer_id"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ساختار پاسخ رویداد همراه با شرکت‌کننده‌هاش
type EventWithParticipantsResponse struct {
	Event        EventResponse  `json:"event"`
	Participants []UserResponse `json:"participants"`
}

type Participant struct {
	ID       int       `json:"id"`
	UserID   int       `json:"user_id"`
	EventID  int       `json:"event_id"`
	JoinedAt time.Time `json:"joined_at"`
}
