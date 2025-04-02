package main

import (
	"fmt"
	"log"
	"os"

	"github.com/event-system/config"
	"github.com/event-system/database"
	"github.com/event-system/docs"
	_ "github.com/event-system/docs"
	"github.com/event-system/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

// @title Event Management System API
// @version 1.0
// @description This is an event management system API with JWT authentication
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@event-system.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Initialize Swagger documentation
	docs.SwaggerInfo.Title = "Event Management System API"
	docs.SwaggerInfo.Description = "This is an event management system API with JWT authentication"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	// Initialize database connection
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create tables if they don't exist
	err = database.CreateTables(db)
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: config.ErrorHandler,
	})

	// Middleware
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))

	// CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Setup routes
	routes.SetupRoutes(app, db)

	// Get port from environment variable or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Server starting on port %s\n", port)
	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
