package controllers

import (
	"github.com/event-system/models"
	"github.com/event-system/services"
	"github.com/gofiber/fiber/v2"
)

// AuthController handles authentication related HTTP requests
type AuthController struct {
	AuthService *services.AuthService
}

// NewAuthController creates a new auth controller instance
func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{AuthService: authService}
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user with username, email, and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.RegisterRequest true "User registration data"
// @Success 201 {object} models.TokenResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /auth/register [post]
func (c *AuthController) Register(ctx *fiber.Ctx) error {
	// Parse request body
	req := new(models.RegisterRequest)
	if err := ctx.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Username, email, and password are required")
	}

	// Register user
	response, err := c.AuthService.Register(*req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Return response
	ctx.Status(fiber.StatusCreated)
	return ctx.JSON(response)
}

// Login handles user login
// @Summary Login a user
// @Description Login a user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.LoginRequest true "User login data"
// @Success 200 {object} models.TokenResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/login [post]
func (c *AuthController) Login(ctx *fiber.Ctx) error {
	// Parse request body
	req := new(models.LoginRequest)
	if err := ctx.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if req.Email == "" || req.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Email and password are required")
	}

	// Login user
	response, err := c.AuthService.Login(*req)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	// Return response
	return ctx.JSON(response)
}

// GetProfile handles getting the current user's profile
// @Summary Get user profile
// @Description Get the current user's profile
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.UserResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/profile [get]
// @Summary Get user profile
func (c *AuthController) GetProfile(ctx *fiber.Ctx) error {
	// Get user ID from context
	userID, ok := ctx.Locals("userID").(int)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	// Get user profile
	user, err := c.AuthService.GetUserByID(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Return response
	return ctx.JSON(user)
}
