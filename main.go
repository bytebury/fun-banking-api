package main

import (
	"golfer/config"
	"golfer/database"
	"golfer/middleware"
	"golfer/models"
	"golfer/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Unable to load .env configuration")
	}

	// Setup configurations/constants
	setupConfigs()
	// Database setup
	database.SetupConnection()
	// Run migrations
	runMigrations()
	// Router setup
	router := gin.Default()
	// Middleware
	router.Use(middleware.CorsMiddleware())
	// Routes
	routes.SetupRoutes(router)

	// Start the application
	router.Run()
}

func runMigrations() {
	database.DB.AutoMigrate(&models.User{})
	database.DB.AutoMigrate(&models.Bank{})
	database.DB.AutoMigrate(&models.Customer{})
}

func setupConfigs() {
	config.JwtKey = []byte(os.Getenv("JWT_SECRET"))
	config.AppBaseURL = os.Getenv("APP_BASE_URL")
}
