package controllers

import (
	"strconv"

	"github.com/event-system/models"
	"github.com/event-system/services"
	"github.com/gofiber/fiber/v2"
)

// EventController handles event related HTTP requests
type EventController struct {
	EventService *services.EventService
}

// NewEventController creates a new event controller instance
func NewEventController(eventService *services.EventService) *EventController {
	return &EventController{EventService: eventService}
}

// CreateEvent handles event creation
// @Summary Create a new event
// @Description Create a new event with name, description, location, start time, end time, and capacity
// @Tags events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param event body models.EventRequest true "Event creation data"
// @Success 201 {object} models.EventResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /events [post]
func (c *EventController) CreateEvent(ctx *fiber.Ctx) error {
	// Get user ID from context
	userID, ok := ctx.Locals("userID").(int)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	// Parse request body
	req := new(models.EventRequest)
	if err := ctx.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body"+err.Error())
	}

	// Validate request
	if req.Name == "" || req.StartTime.IsZero() || req.EndTime.IsZero() || req.Capacity <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Name, start time, end time, and capacity are required")
	}

	// Check if end time is after start time
	if !req.EndTime.After(req.StartTime) {
		return fiber.NewError(fiber.StatusBadRequest, "End time must be after start time")
	}

	// Create event
	event, err := c.EventService.CreateEvent(*req, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Return response
	ctx.Status(fiber.StatusCreated)
	return ctx.JSON(event)
}

// GetEvent handles getting a single event by ID
// @Summary Get an event
// @Description Get an event by ID
// @Tags events
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Success 200 {object} models.EventResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /events/{id} [get]
func (c *EventController) GetEvent(ctx *fiber.Ctx) error {
	// Get event ID from path
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid event ID")
	}

	// Get event
	event, err := c.EventService.GetEventByID(id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	// Return response
	return ctx.JSON(event)
}

// UpdateEvent handles updating an event
// @Summary Update an event
// @Description Update an event with name, description, location, start time, end time, and capacity
// @Tags events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Event ID"
// @Param event body models.EventRequest true "Event update data"
// @Success 200 {object} models.EventResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /events/{id} [put]
func (c *EventController) UpdateEvent(ctx *fiber.Ctx) error {
	// Get user ID from context
	userID, ok := ctx.Locals("userID").(int)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	// Get event ID from path
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid event ID")
	}

	// Parse request body
	req := new(models.EventRequest)
	if err := ctx.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if req.Name == "" || req.StartTime.IsZero() || req.EndTime.IsZero() || req.Capacity <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Name, start time, end time, and capacity are required")
	}

	// Check if end time is after start time
	if !req.EndTime.After(req.StartTime) {
		return fiber.NewError(fiber.StatusBadRequest, "End time must be after start time")
	}

	// Update event
	event, err := c.EventService.UpdateEvent(id, *req, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Return response
	return ctx.JSON(event)
}

// CloseEvent change event status to close
// @Summary Close an event
// @Description Close an event by ID
// @Tags events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Event ID"
// @Success 200 {object} models.EventResponse
// @Router /events/{id}/close [post]
func (c *EventController) CloseEvent(ctx *fiber.Ctx) error {
	// Get user ID from context
	userID, ok := ctx.Locals("userID").(int)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	// Get event ID from path
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid event ID")
	}

	// Update event
	event, err := c.EventService.CloseEvent(userID, id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Return response
	return ctx.JSON(event)
}

// OpenEvent change event status to open
// @Summary Open an event
// @Description Open an event by ID
// @Tags events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Event ID"
// @Success 200 {object} models.EventResponse
// @Router /events/{id}/open [post]
func (c *EventController) OpenEvent(ctx *fiber.Ctx) error {
	// Get user ID from context
	userID, ok := ctx.Locals("userID").(int)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	// Get event ID from path
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid event ID")
	}

	// Update event
	event, err := c.EventService.OpenEvent(userID, id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Return response
	return ctx.JSON(event)
}

// DeleteEvent handles deleting an event
// @Summary Delete an event
// @Description Delete an event by ID
// @Tags events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Event ID"
// @Success 200 {object} models.MessageResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /events/{id} [delete]
func (c *EventController) DeleteEvent(ctx *fiber.Ctx) error {
	// Get user ID from context
	userID, ok := ctx.Locals("userID").(int)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	// Get event ID from path
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid event ID")
	}

	// Delete event
	err = c.EventService.DeleteEvent(id, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Return response
	return ctx.JSON(fiber.Map{
		"message": "Event deleted successfully",
	})
}

// GetAllPublicEvents handles getting all open events
// @Summary Get all open events
// @Description Get all open events
// @Tags events
// @Accept json
// @Produce json
// @Success 200 {array} models.EventResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /events/public [get]
func (c *EventController) GetAllPublicEvents(ctx *fiber.Ctx) error {
	// Get events
	events, err := c.EventService.GetAllPublicEvents()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Return response
	return ctx.JSON(events)
}

// GetMyEvents handles getting all events created by the current user
// @Summary Get my events
// @Description Get all events created by the current user
// @Tags events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.EventResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /events/my/ [get]
func (c *EventController) GetMyEvents(ctx *fiber.Ctx) error {

	// Get user ID from context
	userID, ok := ctx.Locals("userID").(int)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}
	// Get events
	events, err := c.EventService.GetEventsByOrganizer(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Return response
	return ctx.JSON(events)
}

// GetEventWithParticipants handles getting an event with its participants
// @Summary Get event with participants
// @Description Get an event with its participants
// @Tags events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Event ID"
// @Success 200 {object} models.EventWithParticipantsResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /events/{id}/participants [get]
func (c *EventController) GetEventWithParticipants(ctx *fiber.Ctx) error {
	// Get user ID from context
	userID, ok := ctx.Locals("userID").(int)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	// Get event ID from path
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid event ID")
	}

	// Get event with participants
	event, err := c.EventService.GetEventWithParticipants(id, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Return response
	return ctx.JSON(event)
}

// GetMyParticipatingEvents handles getting all events the current user is participating in
// @Summary Get my participating events
// @Description Get all events the current user is participating in
// @Tags events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.EventResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /events/participating [get]
func (c *EventController) GetMyParticipatingEvents(ctx *fiber.Ctx) error {
	// Get user ID from context
	userID, ok := ctx.Locals("userID").(int)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	// Get events
	events, err := c.EventService.GetEventsByParticipant(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Return response
	return ctx.JSON(events)
}
