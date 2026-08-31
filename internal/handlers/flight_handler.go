package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/Varaaa-arch/jelajahin-go/internal/repositories"
	"github.com/Varaaa-arch/jelajahin-go/internal/services"
	"gorm.io/gorm"
)

type FlightHandler struct {
	repo    *repositories.FlightRepository
	service *services.FlightService
}

func NewFlightHandler(repo *repositories.FlightRepository, service *services.FlightService) *FlightHandler {
	return &FlightHandler{
		repo:    repo,
		service: service,
	}
}

// SearchFlights endpoint: GET /api/v1/flights/search
// Query params: origin, destination, departure_date
func (h *FlightHandler) SearchFlights(c *gin.Context) {
	origin := c.Query("origin")
	destination := c.Query("destination")
	departureDate := c.Query("departure_date")

	// Validate params (minimal validation)
	if origin == "" || destination == "" || departureDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing required parameters: origin, destination, departure_date",
		})
		return
	}

	// Call service dengan caching
	flights, err := h.service.SearchFlights(origin, destination, departureDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to search flights",
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, gin.H{
		"data": flights,
		"message": "Flights retrieved successfully",
		"total": len(flights),
	})
}

// GetFlightByID endpoint: GET /api/v1/flights/:id
func (h *FlightHandler) GetFlightByID(c *gin.Context) {
	id := c.Param("id")

	flight, err := h.service.GetFlightByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Flight not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch flight",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": flight,
		"message": "Flight retrieved successfully",
	})
}

// GetAvailableSeats endpoint: GET /api/v1/flights/:id/seats
func (h *FlightHandler) GetAvailableSeats(c *gin.Context) {
	flightID := c.Param("id")

	seats, err := h.repo.GetAvailableSeats(flightID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch available seats",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": seats,
		"total": len(seats),
		"message": "Available seats retrieved successfully",
	})
}
