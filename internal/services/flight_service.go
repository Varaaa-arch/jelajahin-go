package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Varaaa-arch/jelajahin-go/internal/models"
	"github.com/Varaaa-arch/jelajahin-go/internal/repositories"
	"github.com/redis/go-redis/v9"
)

type FlightService struct {
	repo *repositories.FlightRepository
	redis *redis.Client
}

func NewFlightService(repo *repositories.FlightRepository, redis *redis.Client) *FlightService {
	return &FlightService{repo: repo, redis: redis}
}

// SearchFlights dengan caching di Redis
func (s *FlightService) SearchFlights(origin, destination, departureDate string) ([]models.Flight, error) {
	ctx := context.Background()

	// Generate cache key
	cacheKey := fmt.Sprintf("flights:%s:%s:%s", origin, destination, departureDate)

	// Try to get from Redis cache
	cachedData, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		// Cache hit!
		var flights []models.Flight
		if err := json.Unmarshal([]byte(cachedData), &flights); err == nil {
			log.Printf("Cache hit for key: %s\n", cacheKey)
			return flights, nil
		}
	}

	// Cache miss, query database
	log.Printf("Cache miss for key: %s, querying database\n", cacheKey)
	flights, err := s.repo.SearchFlights(origin, destination, departureDate)
	if err != nil {
		return nil, err
	}

	// Store in Redis cache with TTL 1 hour
	flightData, err := json.Marshal(flights)
	if err == nil {
		err = s.redis.Set(ctx, cacheKey, flightData, time.Hour).Err()
		if err != nil {
			log.Printf("Failed to cache flights: %v\n", err)
		} else {
			log.Printf("Cached flights for key: %s (TTL: 1 hour)\n", cacheKey)
		}
	}

	return flights, nil
}

// GetFlightByID untuk detail flight (juga dengan cache)
func (s *FlightService) GetFlightByID(id string) (*models.Flight, error) {
	ctx := context.Background()

	// Cache key
	cacheKey := fmt.Sprintf("flight:%s", id)

	// Try cache
	cachedData, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var flight models.Flight
		if err := json.Unmarshal([]byte(cachedData), &flight); err == nil {
			log.Printf("Cache hit for flight: %s\n", id)
			return &flight, nil
		}
	}

	// Query database
	flight, err := s.repo.GetFlightByID(id)
	if err != nil {
		return nil, err
	}

	// Cache result
	flightData, err := json.Marshal(flight)
	if err == nil {
		s.redis.Set(ctx, cacheKey, flightData, time.Hour).Err()
	}

	return flight, nil
}

// InvalidateSearchCache untuk clear cache saat ada perubahan flight
func (s *FlightService) InvalidateSearchCache(origin, destination, departureDate string) error {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("flights:%s:%s:%s", origin, destination, departureDate)
	
	err := s.redis.Del(ctx, cacheKey).Err()
	if err != nil {
		log.Printf("Failed to invalidate cache: %v\n", err)
	}
	return err
}