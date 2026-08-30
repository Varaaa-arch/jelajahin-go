package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"github.com/Varaaa-arch/jelajahin-go/internal/handlers"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB, redis *redis.Client) {
	// Flight routes
	flightHandler := handlers.NewFlightHandler(db)
	
	api := router.Group("/api")
	{
		flights := api.Group("/flights")
		{
			flights.GET("", flightHandler.SearchFlights)
			flights.GET("/:id", flightHandler.GetFlightByID)
		}
	}
}
