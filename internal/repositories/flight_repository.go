package repositories

import (
	"github.com/Varaaa-arch/jelajahin-go/internal/models"
	"gorm.io/gorm"
)

type FlightRepository struct {
	db *gorm.DB 
}

func NewFlightRepository(db *gorm.DB) *FlightRepository {
	return &FlightRepository{db: db}
}

// Search filter dengan filter origin, destinasi, departure_data
func (r *FlightRepository) SearchFlights(origin, destination, departureData string) ([]models.Flight, error) {
	var flights []models.Flight
	query := r.db

	// Filt by status
	query = query.Where("flights.status = ?", "scheduled")

	// Filt by origin & destination airport (single join dengan routes)
	if origin != "" || destination != "" {
		query = query.Joins("JOIN routes ON flights.route_id = routes.id")
		if origin != "" {
			query = query.Joins("JOIN airports origin_airport ON origin_airport.id = routes.origin_airport_id").
				Where("origin_airport.code = ?", origin)
		}
		if destination != "" {
			query = query.Joins("JOIN airports dest_airport ON dest_airport.id = routes.destination_airport_id").
				Where("dest_airport.code = ?", destination)
		}
	}

	// Filt by departure_date
	if departureData != "" {
		query = query.Where("flights.departure_date = ?", departureData)
	}

	// Execute query
	if err := query.Find(&flights).Error; err != nil {
		return nil, err
	}

	return flights, nil
}

// GetFlightByID untuk get detail flight
func (r *FlightRepository) GetFlightByID(id string) (*models.Flight, error) {
	var flight models.Flight
	if err := r.db.First(&flight, "id = ?", id).Error; err != nil{
		return nil, err
	}

	return &flight, nil
}

// GetAvailableSeats untuk seat yang belum booked
func (r *FlightRepository) GetAvailableSeats(flightID string) ([]models.FlightSeat, error){
	var seats []models.FlightSeat

	if err := r.db.Where("flight_id = ? AND is_available = ?", flightID, true).Find(&seats).Error; err != nil {
		return nil, err
	}

	return seats, nil
}