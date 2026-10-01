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

	// Initialize seat dependencies
	seatService := services.NewSeatService(redis)
	seatHandler := handlers.NewSeatHandler(seatService)

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

		// Seat routes (public untuk lock/unlock agar flow web-session + guest tetap jalan.
		// user_id dikirim di body dan kepemilikan lock dicek di Redis via SeatService.
		// Jika Authorization tersedia, tetap divalidasi; jika tidak, lanjut sebagai guest.)
		seats := api.Group("/seats")
		{
			seats.POST("/lock", seatHandler.LockSeat)
			seats.POST("/unlock", seatHandler.UnlockSeat)
			seats.POST("/lock-multiple", seatHandler.LockMultipleSeats)
			seats.POST("/extend", seatHandler.ExtendSeatHold)
			seats.POST("/confirm", seatHandler.ConfirmSeats)
			seats.GET("/:flight_id/:seat_id/lock-status", seatHandler.CheckSeatLock)
		}
	}
}
