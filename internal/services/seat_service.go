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
	// ErrNotLockOwner dikembalikan saat user bukan pemilik lock.
	ErrNotLockOwner = errors.New("not the lock owner")
	// ErrSeatAlreadyBooked dikembalikan saat kursi sudah dikonfirmasi/booked.
	ErrSeatAlreadyBooked = errors.New("seat already booked")
)

// BookedSeatTTL durasi tanda kursi sudah terpakai (1 tahun).
// Kursi booked disimpan terpisah dari lock supaya tidak teroverwrite.
const BookedSeatTTL = 365 * 24 * time.Hour

// DefaultLockTTL durasi default lock kursi.
const DefaultLockTTL = 15 * time.Minute

// MaxHoldTTL batas maksimal perpanjangan tahan kursi (48 jam).
const MaxHoldTTL = 48 * time.Hour

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
// Hanya pemilik lock (userID yang sama) yang diizinkan.
// Mengembalikan ErrSeatNotLocked saat kursi sebenarnya tidak dikunci,
// dan ErrNotLockOwner saat userID bukan pemilik lock.
func (s *SeatService) UnlockSeat(flightID, seatID, userID string) error {
	ctx := context.Background()

	lockKey := fmt.Sprintf("lock:seat:%s:%s", flightID, seatID)

	// Ambil siapa pemilik lock
	lockedBy, err := s.redis.Get(ctx, lockKey).Result()
	if err != nil {
		if err == redis.Nil {
			return ErrSeatNotLocked
		}
		log.Printf("Error checking seat before unlock: %v\n", err)
		return err
	}

	// Validasi kepemilikan
	if lockedBy != userID {
		log.Printf("Unlock rejected: seat %s locked by %s, requested by %s\n", seatID, lockedBy, userID)
		return ErrNotLockOwner
	}

	err = s.redis.Del(ctx, lockKey).Err()
	if err != nil {
		log.Printf("Error unlocking seat: %v\n", err)
		return err
	}

	log.Printf("Seat %s unlocked by owner %s\n", seatID, userID)
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

	// If any seat failed to lock, unlock all yang berhasil (internal rollback, bypass ownership check)
	if len(failedSeats) > 0 {
		for _, seatID := range lockedSeats {
			s.forceUnlockSeat(flightID, seatID)
		}
		log.Printf("Failed to lock some seats: %v\n", failedSeats)
		return false, failedSeats, ErrSeatAlreadyLocked
	}

	log.Printf("Successfully locked %d seats\n", len(lockedSeats))
	return true, lockedSeats, nil
}

// forceUnlockSeat adalah helper internal yang menghapus lock tanpa validasi ownership.
// Digunakan hanya untuk rollback saat LockMultipleSeats gagal sebagian.
func (s *SeatService) forceUnlockSeat(flightID, seatID string) {
	ctx := context.Background()
	lockKey := fmt.Sprintf("lock:seat:%s:%s", flightID, seatID)
	if err := s.redis.Del(ctx, lockKey).Err(); err != nil {
		log.Printf("Rollback: error unlocking seat %s: %v\n", seatID, err)
	}
}

// ConfirmSeat menandai kursi sebagai booked (setelah payment sukses).
// Menghapus lock dan menyimpan booked key secara permanen (TTL 1 tahun).
func (s *SeatService) ConfirmSeat(flightID, seatID, userID string) error {
	ctx := context.Background()

	lockKey := fmt.Sprintf("lock:seat:%s:%s", flightID, seatID)
	bookedKey := fmt.Sprintf("booked:seat:%s:%s", flightID, seatID)

	// Cek apakah sudah booked
	alreadyBooked, err := s.redis.Exists(ctx, bookedKey).Result()
	if err != nil {
		return err
	}
	if alreadyBooked > 0 {
		return ErrSeatAlreadyBooked
	}

	// Tandai sebagai booked (TTL panjang = 1 tahun)
	if err := s.redis.Set(ctx, bookedKey, userID, BookedSeatTTL).Err(); err != nil {
		log.Printf("Error confirming seat %s: %v\n", seatID, err)
		return err
	}

	// Hapus lock (jika masih ada)
	s.redis.Del(ctx, lockKey)

	log.Printf("Seat %s (flight %s) confirmed/booked by user %s\n", seatID, flightID, userID)
	return nil
}

// ConfirmMultipleSeats konfirmasi banyak kursi sekaligus.
// Mengembalikan daftar kursi yang gagal dikonfirmasi.
func (s *SeatService) ConfirmMultipleSeats(flightID string, seatIDs []string, userID string) ([]string, error) {
	failed := []string{}

	for _, seatID := range seatIDs {
		if err := s.ConfirmSeat(flightID, seatID, userID); err != nil {
			if err == ErrSeatAlreadyBooked {
				log.Printf("Seat %s already booked, skipping\n", seatID)
				// Already booked = idempotent, bukan error fatal
				continue
			}
			log.Printf("Error confirming seat %s: %v\n", seatID, err)
			failed = append(failed, seatID)
		}
	}

	if len(failed) > 0 {
		return failed, fmt.Errorf("failed to confirm %d seat(s): %v", len(failed), failed)
	}

	log.Printf("Successfully confirmed %d seats for flight %s\n", len(seatIDs), flightID)
	return nil, nil
}

// ExtendSeatLock memperpanjang (atau mengunci ulang) kursi selama menunggu
// persetujuan admin. Dipanggil Laravel setelah payment sukses.
// - Kursi yang sudah booked: dilewati (idempoten).
// - Lock milik user yang sama: TTL diperpanjang.
// - Lock hilang/kedaluwarsa: dikunci ulang bila belum booked.
// - Lock milik orang lain: gagal.
func (s *SeatService) ExtendSeatLock(flightID string, seatIDs []string, userID string, ttl time.Duration) ([]string, error) {
	if ttl <= 0 {
		ttl = DefaultLockTTL
	}
	if ttl > MaxHoldTTL {
		ttl = MaxHoldTTL
	}

	ctx := context.Background()
	failed := []string{}

	for _, seatID := range seatIDs {
		lockKey := fmt.Sprintf("lock:seat:%s:%s", flightID, seatID)
		bookedKey := fmt.Sprintf("booked:seat:%s:%s", flightID, seatID)

		// Sudah booked permanen -> anggap sukses, tidak perlu hold.
		booked, err := s.redis.Exists(ctx, bookedKey).Result()
		if err != nil {
			failed = append(failed, seatID)
			continue
		}
		if booked > 0 {
			continue
		}

		lockedBy, err := s.redis.Get(ctx, lockKey).Result()
		if err != nil && err != redis.Nil {
			failed = append(failed, seatID)
			continue
		}
		if err == nil && lockedBy != userID {
			failed = append(failed, seatID)
			continue
		}
		if err == nil {
			// Milik sendiri -> perpanjang TTL.
			if ok, err := s.redis.Expire(ctx, lockKey, ttl).Result(); err != nil || !ok {
				failed = append(failed, seatID)
				continue
			}
			log.Printf("Seat %s (flight %s) hold extended for user %s (TTL: %s)\n", seatID, flightID, userID, ttl)
			continue
		}

		// Lock hilang -> kunci ulang.
		ok, err := s.redis.SetNX(ctx, lockKey, userID, ttl).Result()
		if err != nil || !ok {
			failed = append(failed, seatID)
			continue
		}
		log.Printf("Seat %s (flight %s) re-locked for user %s (TTL: %s)\n", seatID, flightID, userID, ttl)
		s.watchLockExpiry(flightID, seatID, userID)
	}

	if len(failed) > 0 {
		return failed, fmt.Errorf("failed to extend %d seat(s): %v", len(failed), failed)
	}

	return nil, nil
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
