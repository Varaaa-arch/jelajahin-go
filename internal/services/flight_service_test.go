package services

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/Varaaa-arch/jelajahin-go/internal/models"
	"gorm.io/gorm"
)

func newTestRedis(t *testing.T) *redis.Client {
	t.Helper()
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	t.Cleanup(mr.Close)
	return redis.NewClient(&redis.Options{Addr: mr.Addr()})
}

// Mock repository untuk testing
type MockFlightRepository struct {
	flights []models.Flight
}

func (m *MockFlightRepository) SearchFlights(origin, destination, departureDate string) ([]models.Flight, error) {
	return m.flights, nil
}

func (m *MockFlightRepository) GetFlightByID(id string) (*models.Flight, error) {
	for _, f := range m.flights {
		if f.ID == id {
			return &f, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *MockFlightRepository) GetAvailableSeats(flightID string) ([]models.FlightSeat, error) {
	return []models.FlightSeat{}, nil
}

// Test SearchFlights function
func TestSearchFlights(t *testing.T) {
	// Mock data
	mockFlights := []models.Flight{
		{
			ID:           "1",
			FlightNumber: "GA100",
			BasePrice:    1000000,
			Status:       "scheduled",
		},
	}

	mockRepo := &MockFlightRepository{flights: mockFlights}
	service := NewFlightService(mockRepo, newTestRedis(t))

	// Test search
	flights, err := service.SearchFlights("CGK", "DPS", "2024-09-15")

	if err != nil {
		t.Errorf("SearchFlights failed: %v", err)
	}

	if len(flights) != 1 {
		t.Errorf("Expected 1 flight, got %d", len(flights))
	}

	if flights[0].FlightNumber != "GA100" {
		t.Errorf("Expected flight GA100, got %s", flights[0].FlightNumber)
	}

	// Search again should hit Redis cache (ab) — verify it still returns flights
	cached, err := service.SearchFlights("CGK", "DPS", "2024-09-15")
	if err != nil || len(cached) != 1 || cached[0].FlightNumber != "GA100" {
		t.Errorf("Cached result mismatch: err=%v flights=%+v", err, cached)
	}
}

// Test GetFlightByID
func TestGetFlightByID(t *testing.T) {
	mockFlights := []models.Flight{
		{
			ID:           "1",
			FlightNumber: "GA100",
		},
	}

	mockRepo := &MockFlightRepository{flights: mockFlights}
	service := NewFlightService(mockRepo, newTestRedis(t))

	flight, err := service.GetFlightByID("1")

	if err != nil {
		t.Errorf("GetFlightByID failed: %v", err)
	}

	if flight.FlightNumber != "GA100" {
		t.Errorf("Expected GA100, got %s", flight.FlightNumber)
	}
}
