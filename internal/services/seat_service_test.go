package services

import (
	"errors"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestLockSeat(t *testing.T) {
	// Use Redis test instance (make sure Redis running)
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()

	service := NewSeatService(client)
	flightID := "svc-flight-1"
	seatID := "SVC-12A"
	userID := "user-123"

	// Test lock success
	locked, lockedBy, err := service.LockSeat(flightID, seatID, userID)

	if err != nil {
		t.Errorf("LockSeat failed: %v", err)
	}

	if !locked {
		t.Errorf("Expected seat to be locked, got: %v", locked)
	}

	if lockedBy != userID {
		t.Errorf("Expected locked by %s, got %s", userID, lockedBy)
	}

	// Cleanup
	service.UnlockSeat(flightID, seatID, userID)
}

func TestLockSeatAlreadyLocked(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()

	service := NewSeatService(client)
	flightID := "svc-flight-2"
	seatID := "SVC-12B"
	userID1 := "user-123"
	userID2 := "user-456"

	// Lock dengan user 1
	service.LockSeat(flightID, seatID, userID1)

	// Try lock dengan user 2 (harus fail)
	locked, lockedBy, err := service.LockSeat(flightID, seatID, userID2)

	if !errors.Is(err, ErrSeatAlreadyLocked) {
		t.Errorf("Expected ErrSeatAlreadyLocked, got: %v", err)
	}

	if locked {
		t.Errorf("Expected seat to be already locked, got: %v", locked)
	}

	if lockedBy != userID1 {
		t.Errorf("Expected locked by %s, got %s", userID1, lockedBy)
	}

	// Cleanup
	service.UnlockSeat(flightID, seatID, userID1)
}

func TestCheckSeatLock(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()

	service := NewSeatService(client)
	flightID := "svc-flight-3"
	seatID := "SVC-12C"
	userID := "user-123"

	// Lock seat
	service.LockSeat(flightID, seatID, userID)

	// Check if locked
	isLocked, lockedBy, err := service.CheckSeatLock(flightID, seatID)

	if err != nil {
		t.Errorf("CheckSeatLock failed: %v", err)
	}

	if !isLocked {
		t.Errorf("Expected seat to be locked")
	}

	if lockedBy != userID {
		t.Errorf("Expected locked by %s, got %s", userID, lockedBy)
	}

	// Cleanup
	service.UnlockSeat(flightID, seatID, userID)
}

func TestGetSeatLockTTL(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()

	service := NewSeatService(client)
	flightID := "svc-flight-4"
	seatID := "SVC-12D"
	userID := "user-123"

	// Lock seat
	service.LockSeat(flightID, seatID, userID)

	// Get TTL
	ttl, err := service.GetSeatLockTTL(flightID, seatID)

	if err != nil {
		t.Errorf("GetSeatLockTTL failed: %v", err)
	}

	// TTL should be close to 15 minutes
	if ttl <= 0 {
		t.Errorf("Expected TTL > 0, got %v", ttl)
	}

	if ttl > 15*time.Minute {
		t.Errorf("Expected TTL <= 15 min, got %v", ttl)
	}

	// Cleanup
	service.UnlockSeat(flightID, seatID, userID)
}

func TestLockMultipleSeats(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()

	service := NewSeatService(client)
	flightID := "svc-flight-5"
	seatIDs := []string{"SVC-12E", "SVC-12F", "SVC-12G"}
	userID := "user-123"

	// Lock multiple seats
	locked, lockedSeats, err := service.LockMultipleSeats(flightID, seatIDs, userID)

	if err != nil {
		t.Errorf("LockMultipleSeats failed: %v", err)
	}

	if !locked {
		t.Errorf("Expected all seats to be locked")
	}

	if len(lockedSeats) != len(seatIDs) {
		t.Errorf("Expected %d locked seats, got %d", len(seatIDs), len(lockedSeats))
	}

	// Cleanup
	for _, seatID := range seatIDs {
		service.UnlockSeat(flightID, seatID, userID)
	}
}