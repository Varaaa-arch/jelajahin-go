package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"github.com/Varaaa-arch/jelajahin-go/internal/handlers"
	"github.com/Varaaa-arch/jelajahin-go/internal/repositories"
	"github.com/Varaaa-arch/jelajahin-go/internal/services"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB, redis *redis.Client) {
	// Initialize flight dependencies
	flightRepo := repositories.NewFlightRepository(db)
	flightService := services.NewFlightService(flightRepo, redis)
	flightHandler := handlers.NewFlightHandler(flightRepo, flightService)

	// API v1 routes
	api := router.Group("/api/v1")
	{
		// Flight routes
		flights := api.Group("/flights")
		{
			flights.GET("/search", flightHandler.SearchFlights)
			flights.GET("/:id", flightHandler.GetFlightByID)
			flights.GET("/:id/seats", flightHandler.GetAvailableSeats)
		}
	}
}
