package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// Sentinel error untuk membedakan kegagalan bisnis vs teknis.
var (
	// ErrSeatAlreadyLocked dikembalikan saat kursi sudah dikunci user lain.
	ErrSeatAlreadyLocked = errors.New("seat already locked")
	// ErrSeatNotLocked dikembalikan saat mencoba unlock/check kursi yang tidak dikunci.
	ErrSeatNotLocked = errors.New("seat is not locked")
)

// DefaultLockTTL durasi default lock kursi.
const DefaultLockTTL = 15 * time.Minute

type SeatService struct {
	redis *redis.Client
}

func NewSeatService(redis *redis.Client) *SeatService {
	return &SeatService{
		redis: redis,
	}
}

// LockSeat mengunci kursi selama 15 menit.
// Saat kursi sudah dikunci user lain, mengembalikan (false, lockedBy, ErrSeatAlreadyLocked).
func (s *SeatService) LockSeat(flightID, seatID, userID string) (bool, string, error) {
	ctx := context.Background()

	// Generate lock key
	lockKey := fmt.Sprintf("lock:seat:%s:%s", flightID, seatID)

	// Try to set lock (SET NX = only set if not exists)
	result, err := s.redis.SetNX(ctx, lockKey, userID, DefaultLockTTL).Result()
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
		return false, lockedBy, ErrSeatAlreadyLocked
	}

	// Lock successful
	log.Printf("Seat %s locked by user %s (TTL: %s)\n", seatID, userID, DefaultLockTTL)

	// Jadwalkan notifikasi mendekati / saat lock kedaluwarsa.
	s.watchLockExpiry(flightID, seatID, userID)

	return true, userID, nil
}

// UnlockSeat menghapus lock (manual unlock).
// Mengembalikan ErrSeatNotLocked saat kursi sebenarnya tidak dikunci.
func (s *SeatService) UnlockSeat(flightID, seatID string) error {
	ctx := context.Background()

	lockKey := fmt.Sprintf("lock:seat:%s:%s", flightID, seatID)

	exists, err := s.redis.Exists(ctx, lockKey).Result()
	if err != nil {
		log.Printf("Error checking seat before unlock: %v\n", err)
		return err
	}
	if exists == 0 {
		return ErrSeatNotLocked
	}

	err = s.redis.Del(ctx, lockKey).Err()
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

		result, err := s.redis.SetNX(ctx, lockKey, userID, DefaultLockTTL).Result()
		if err != nil {
			failedSeats = append(failedSeats, seatID)
			continue
		}

		if result {
			lockedSeats = append(lockedSeats, seatID)
			s.watchLockExpiry(flightID, seatID, userID)
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
		return false, failedSeats, ErrSeatAlreadyLocked
	}

	log.Printf("Successfully locked %d seats\n", len(lockedSeats))
	return true, lockedSeats, nil
}

// watchLockExpiry menjalankan goroutine yang mencatat notifikasi saat lock
// hampir kedaluwarsa dan saat kedaluwarsa (expiry notification).
func (s *SeatService) watchLockExpiry(flightID, seatID, userID string) {
	lockKey := fmt.Sprintf("lock:seat:%s:%s", flightID, seatID)

	// Notifikasi saat TTL menyisakan < 3 menit.
	warnAfter := DefaultLockTTL - (3 * time.Minute)
	time.AfterFunc(warnAfter, func() {
		ctx := context.Background()
		ttl, err := s.redis.TTL(ctx, lockKey).Result()
		if err != nil || ttl <= 0 {
			return // lock sudah kedaluwarsa/dihapus
		}
		log.Printf("[LOCK-EXPIRY] Seat %s (flight %s) locked by %s will expire in ~%s\n",
			seatID, flightID, userID, ttl)
	})

	// Notifikasi saat lock kedaluwarsa.
	time.AfterFunc(DefaultLockTTL+time.Second, func() {
		ctx := context.Background()
		exists, err := s.redis.Exists(ctx, lockKey).Result()
		if err != nil {
			return
		}
		if exists == 0 {
			log.Printf("[LOCK-EXPIRED] Seat %s (flight %s) lock expired and released\n", seatID, flightID)
		}
	})
}
