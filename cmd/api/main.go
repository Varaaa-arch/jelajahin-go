package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/Varaaa-arch/jelajahin-go/pkg/cache"
	"github.com/Varaaa-arch/jelajahin-go/pkg/database"
	"github.com/Varaaa-arch/jelajahin-go/internal/routes"
	"github.com/Varaaa-arch/jelajahin-go/internal/middleware"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Initialize database
	db := database.InitPostgres()
	
	// Initialize Redis
	redis := cache.InitRedis()

	// Set Gin mode
	if os.Getenv("SERVER_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := gin.Default()

	// Middleware
	router.Use(middleware.CorsMiddleware())
	router.Use(middleware.RequestLoggerMiddleware())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"message": "API is running",
		})
	})

	// Setup routes with db and redis
	routes.SetupRoutes(router, db, redis)

	// Start server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
