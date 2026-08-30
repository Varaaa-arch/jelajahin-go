package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/Varaaa-arch/jelajahin-go/internal/models"
)

type FlightHandler struct {
	db *gorm.DB
}

func NewFlightHandler(db *gorm.DB) *FlightHandler {
	return &FlightHandler{db: db}
}

func (h *FlightHandler) SearchFlights(c *gin.Context) {
	origin := c.Query("origin")
	destination := c.Query("destination")
	departureDate := c.Query("departure_date")

	var flights []models.Flight

	query := h.db.
		Joins("JOIN routes ON routes.id = flights.route_id").
		Joins("JOIN airports origin_airport ON origin_airport.id = routes.origin_airport_id").
		Joins("JOIN airports destination_airport ON destination_airport.id = routes.destination_airport_id").
		Where("flights.status = ?", "scheduled")

	if origin != "" {
		query = query.Where("origin_airport.code = ?", origin)
	}
	if destination != "" {
		query = query.Where("destination_airport.code = ?", destination)
	}
	if departureDate != "" {
		query = query.Where("flights.departure_date = ?", departureDate)
	}

	if err := query.Find(&flights).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to fetch flights"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   flights,
		"status": "success",
	})
}

func (h *FlightHandler) GetFlightByID(c *gin.Context) {
	id := c.Param("id")

	var flight models.Flight
	if err := h.db.First(&flight, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Flight not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   flight,
		"status": "success",
	})
}
