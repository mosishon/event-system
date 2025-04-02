package controllers

import (
	"strconv"

	"github.com/event-system/models"
	"github.com/event-system/services"
	"github.com/gofiber/fiber/v2"
)

// ParticipantController handles participant related HTTP requests
type ParticipantController struct {
	ParticipantService *services.ParticipantService
}

// NewParticipantController creates a new participant controller instance
func NewParticipantController(participantService *services.ParticipantService) *ParticipantController {
	return &ParticipantController{ParticipantService: participantService}
}

// JoinEvent handles joining an event
// @Summary Join an event
// @Description Join an event as a participant
// @Tags participants
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Event ID"
// @Success 200 {object} models.MessageResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /events/{id}/join [post]
func (c *ParticipantController) JoinEvent(ctx *fiber.Ctx) error {
	// Get user ID from context
	userID, ok := ctx.Locals("userID").(int)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	// Get event ID from path
	eventID, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid event ID")
	}

	// Join event
	err = c.ParticipantService.JoinEvent(userID, eventID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Return response
	return ctx.JSON(fiber.Map{
		"message": "Successfully joined event",
	})
}

// LeaveEvent handles leaving an event
// @Summary Leave an event
// @Description Leave an event as a participant
// @Tags participants
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Event ID"
// @Success 200 {object} models.MessageResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /events/{id}/leave [post]
func (c *ParticipantController) LeaveEvent(ctx *fiber.Ctx) error {
	// Get user ID from context
	userID, ok := ctx.Locals("userID").(int)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	// Get event ID from path
	eventID, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid event ID")
	}

	// Leave event
	err = c.ParticipantService.LeaveEvent(userID, eventID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Return response
	return ctx.JSON(fiber.Map{
		"message": "Successfully left event",
	})
}

// IsParticipant handles checking if a user is a participant of an event
// @Summary Check if user is a participant
// @Description Check if the current user is a participant of an event
// @Tags participants
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Event ID"
// @Success 200 {object} models.ParticipantStatusResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /events/{id}/is-participant [get]
func (c *ParticipantController) IsParticipant(ctx *fiber.Ctx) error {
	// Get user ID from context
	userID, ok := ctx.Locals("userID").(int)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	// Get event ID from path
	eventID, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid event ID")
	}

	// Check if user is a participant
	isParticipant, err := c.ParticipantService.IsParticipant(userID, eventID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Return response
	return ctx.JSON(models.ParticipantStatusResponse{
		IsParticipant: isParticipant,
	})
}

// GetParticipantCount handles getting the number of participants for an event
// @Summary Get participant count
// @Description Get the number of participants for an event
// @Tags participants
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Success 200 {object} models.ParticipantCountResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /events/{id}/participant-count [get]
func (c *ParticipantController) GetParticipantCount(ctx *fiber.Ctx) error {
	// Get event ID from path
	eventID, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid event ID")
	}

	// Get participant count
	count, err := c.ParticipantService.GetParticipantCount(eventID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Return response
	return ctx.JSON(models.ParticipantCountResponse{
		Count: count,
	})
}
