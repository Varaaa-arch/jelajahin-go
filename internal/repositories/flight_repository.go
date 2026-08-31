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
func (r *FlightRepository) SearchFlights(origin, destination, departureData string) ([]models.Flight, error){
	var flights []models.Flight
	query := r.db

	// Filt by status
	query = query.Where("status = ?", "scheduled")

	// Filt by origin airport
	if origin != "" {
		// Join dengan route table untuk cek origin_airport_id
		query = query.Joins("JOIN routes ON flight.routes_id = routes.id").
			Where("routes.origin_airport_id = ?", origin)
	}

	// Filt by destination airport 
	if destination != "" {
		// Join dengan route table untuk cek destination_airport_id
		query = query.Joins("JOIN routes ON flight.routes_id = routes.id").
			Where("routes.destination_airport_id = ?", destination)
	}

	// Filt by departure_date
	if departureData != "" {
		query = query.Where("departure_date = ?", departureData)
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