package routes

import (
	"database/sql"

	"github.com/event-system/controllers"
	"github.com/event-system/middleware"
	"github.com/event-system/repositories"
	"github.com/event-system/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(app *fiber.App, db *sql.DB) {
	// Create repositories
	userRepo := repositories.NewUserRepository(db)
	eventRepo := repositories.NewEventRepository(db)
	participantRepo := repositories.NewParticipantRepository(db)

	// Create services
	authService := services.NewAuthService(userRepo)
	eventService := services.NewEventService(eventRepo, participantRepo)
	participantService := services.NewParticipantService(participantRepo)

	// Create controllers
	authController := controllers.NewAuthController(authService)
	eventController := controllers.NewEventController(eventService)
	participantController := controllers.NewParticipantController(participantService)

	// Protected middleware
	protectedMiddleware := middleware.Protected(authService)
	// API routes
	api := app.Group("/api")

	// Swagger documentation
	api.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/register", authController.Register)
	auth.Post("/login", authController.Login)
	auth.Get("/profile", protectedMiddleware, authController.GetProfile)

	// Public event routes
	events := api.Group("/events")
	events.Get("/public", eventController.GetAllPublicEvents)
	events.Get("/:id<int>", eventController.GetEvent)
	events.Get("/:id<int>/participant-count", participantController.GetParticipantCount)

	// Protected event routes
	events.Post("/", protectedMiddleware, eventController.CreateEvent)
	events.Put("/:id<int>", protectedMiddleware, eventController.UpdateEvent)
	events.Post("/:id/close", protectedMiddleware, eventController.CloseEvent)
	events.Post("/:id/open", protectedMiddleware, eventController.OpenEvent)
	events.Delete("/:id<int>", protectedMiddleware, eventController.DeleteEvent)
	events.Get("/my/", protectedMiddleware, eventController.GetMyEvents)
	events.Get("/participating", protectedMiddleware, eventController.GetMyParticipatingEvents)
	events.Get("/:id<int>/participants", protectedMiddleware, eventController.GetEventWithParticipants)

	// Participant routes
	events.Post("/:id<int>/join", protectedMiddleware, participantController.JoinEvent)
	events.Post("/:id<int>/leave", protectedMiddleware, participantController.LeaveEvent)
	events.Get("/:id<int>/is-participant", protectedMiddleware, participantController.IsParticipant)

	// Add request logger middleware for API routes
	api.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))
}
