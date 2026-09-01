package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type SeatService struct {
	redis *redis.Client
}

func NewSeatService(redis *redis.Client) *SeatService {
	return &SeatService{
		redis: redis,
	}
}

// LockSeat mengunci kursi selama 15 menit
func (s *SeatService) LockSeat(flightID, seatID, userID string) (bool, string, error) {
	ctx := context.Background()

	// Generate lock key
	lockKey := fmt.Sprintf("lock:seat:%s:%s", flightID, seatID)

	// Try to set lock (SET NX = only set if not exists)
	result, err := s.redis.SetNX(ctx, lockKey, userID, 15*time.Minute).Result()
	if err != nil {
		log.Printf("Error locking seat: %v\n", err)
		return false, "", err
	}

	if !result {
		// Lock failed - seat already locked
		// Get who locked it
		lockedBy, err := s.redis.Get(ctx, lockKey).Result()
		if err != nil {
			lockedBy = "unknown"
		}
		log.Printf("Seat %s already locked by %s\n", seatID, lockedBy)
		return false, lockedBy, nil
	}

	// Lock successful
	log.Printf("Seat %s locked by user %s (TTL: 15 min)\n", seatID, userID)
	
	// Retrieve TTL to confirm the lock was applied (value not needed here)
	_, _ = s.redis.TTL(ctx, lockKey).Result()

	return true, userID, nil
}

// UnlockSeat menghapus lock (manual unlock)
func (s *SeatService) UnlockSeat(flightID, seatID string) error {
	ctx := context.Background()

	lockKey := fmt.Sprintf("lock:seat:%s:%s", flightID, seatID)

	err := s.redis.Del(ctx, lockKey).Err()
	if err != nil {
		log.Printf("Error unlocking seat: %v\n", err)
		return err
	}

	log.Printf("Seat %s unlocked\n", seatID)
	return nil
}

// CheckSeatLock cek apakah kursi sedang locked
func (s *SeatService) CheckSeatLock(flightID, seatID string) (bool, string, error) {
	ctx := context.Background()

	lockKey := fmt.Sprintf("lock:seat:%s:%s", flightID, seatID)

	// Try to get lock value
	lockedBy, err := s.redis.Get(ctx, lockKey).Result()
	if err != nil {
		if err == redis.Nil {
			// Lock not found, seat is free
			return false, "", nil
		}
		log.Printf("Error checking seat lock: %v\n", err)
		return false, "", err
	}

	// Seat is locked
	return true, lockedBy, nil
}

// GetSeatLockTTL get remaining time untuk lock
func (s *SeatService) GetSeatLockTTL(flightID, seatID string) (time.Duration, error) {
	ctx := context.Background()

	lockKey := fmt.Sprintf("lock:seat:%s:%s", flightID, seatID)

	ttl, err := s.redis.TTL(ctx, lockKey).Result()
	if err != nil {
		log.Printf("Error getting seat lock TTL: %v\n", err)
		return 0, err
	}

	return ttl, nil
}

// LockMultipleSeats lock multiple seats sekaligus (untuk booking dengan multiple passengers)
func (s *SeatService) LockMultipleSeats(flightID string, seatIDs []string, userID string) (bool, []string, error) {
	ctx := context.Background()

	lockedSeats := []string{}
	failedSeats := []string{}

	for _, seatID := range seatIDs {
		lockKey := fmt.Sprintf("lock:seat:%s:%s", flightID, seatID)

		result, err := s.redis.SetNX(ctx, lockKey, userID, 15*time.Minute).Result()
		if err != nil {
			failedSeats = append(failedSeats, seatID)
			continue
		}

		if result {
			lockedSeats = append(lockedSeats, seatID)
		} else {
			failedSeats = append(failedSeats, seatID)
		}
	}

	// If any seat failed to lock, unlock all yang berhasil
	if len(failedSeats) > 0 {
		for _, seatID := range lockedSeats {
			s.UnlockSeat(flightID, seatID)
		}
		log.Printf("Failed to lock some seats: %v\n", failedSeats)
		return false, failedSeats, nil
	}

	log.Printf("Successfully locked %d seats\n", len(lockedSeats))
	return true, lockedSeats, nil
}
