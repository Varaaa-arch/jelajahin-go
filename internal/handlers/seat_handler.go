package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/Varaaa-arch/jelajahin-go/internal/services"
)

type SeatHandler struct {
	seatService *services.SeatService
}

func NewSeatHandler(seatService *services.SeatService) *SeatHandler {
	return &SeatHandler{
		seatService: seatService,
	}
}

// LockSeatRequest body untuk POST /api/v1/seats/lock
type LockSeatRequest struct {
	FlightID string `json:"flight_id" binding:"required"`
	SeatID   string `json:"seat_id" binding:"required"`
	UserID   string `json:"user_id" binding:"required"`
}

// LockSeatResponse response dari lock seat
type LockSeatResponse struct {
	Success    bool      `json:"success"`
	Message    string    `json:"message"`
	LockToken  string    `json:"lock_token"`
	LockedBy   string    `json:"locked_by,omitempty"`
	ExpiryTime time.Time `json:"expiry_time,omitempty"`
	TTLSeconds int64     `json:"ttl_seconds,omitempty"`
}

// LockSeat endpoint: POST /api/v1/seats/lock
func (h *SeatHandler) LockSeat(c *gin.Context) {
	var req LockSeatRequest

	// Validate request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
			"code":    http.StatusBadRequest,
			"details": err.Error(),
		})
		return
	}

	// Try to lock seat
	_, lockedBy, err := h.seatService.LockSeat(req.FlightID, req.SeatID, req.UserID)
	if err != nil {
		if errors.Is(err, services.ErrSeatAlreadyLocked) {
			c.JSON(http.StatusConflict, gin.H{
				"success":   false,
				"message":   "Seat already locked",
				"code":      http.StatusConflict,
				"locked_by": lockedBy,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to lock seat",
			"code":    http.StatusInternalServerError,
			"details": err.Error(),
		})
		return
	}

	// Lock successful - generate lock token
	lockToken := req.FlightID + ":" + req.SeatID + ":" + req.UserID
	expiryTime := time.Now().Add(services.DefaultLockTTL)

	c.JSON(http.StatusOK, LockSeatResponse{
		Success:    true,
		Message:    "Seat locked successfully",
		LockToken:  lockToken,
		LockedBy:   req.UserID,
		ExpiryTime: expiryTime,
		TTLSeconds: int64(services.DefaultLockTTL.Seconds()), // 15 minutes = 900 seconds
	})
}

// UnlockSeat endpoint: POST /api/v1/seats/unlock
type UnlockSeatRequest struct {
	FlightID string `json:"flight_id" binding:"required"`
	SeatID   string `json:"seat_id" binding:"required"`
	UserID   string `json:"user_id" binding:"required"`
}

func (h *SeatHandler) UnlockSeat(c *gin.Context) {
	var req UnlockSeatRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body: flight_id, seat_id, dan user_id wajib diisi",
			"code":    http.StatusBadRequest,
		})
		return
	}

	err := h.seatService.UnlockSeat(req.FlightID, req.SeatID, req.UserID)
	if err != nil {
		if errors.Is(err, services.ErrSeatNotLocked) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Seat is not locked",
				"code":    http.StatusNotFound,
			})
			return
		}
		if errors.Is(err, services.ErrNotLockOwner) {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "You are not the owner of this lock",
				"code":    http.StatusForbidden,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to unlock seat",
			"code":    http.StatusInternalServerError,
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Seat unlocked successfully",
	})
}

// ConfirmSeatsRequest body untuk POST /api/v1/seats/confirm
type ConfirmSeatsRequest struct {
	FlightID string   `json:"flight_id" binding:"required"`
	SeatIDs  []string `json:"seat_ids" binding:"required,min=1"`
	UserID   string   `json:"user_id" binding:"required"`
}

// ConfirmSeats endpoint: POST /api/v1/seats/confirm
// Dipanggil oleh Laravel setelah payment berhasil.
// Menandai kursi sebagai booked di Redis dan menghapus lock.
func (h *SeatHandler) ConfirmSeats(c *gin.Context) {
	var req ConfirmSeatsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body: flight_id, seat_ids, dan user_id wajib diisi",
			"code":    http.StatusBadRequest,
			"details": err.Error(),
		})
		return
	}

	failedSeats, err := h.seatService.ConfirmMultipleSeats(req.FlightID, req.SeatIDs, req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":      false,
			"message":      "Some seats could not be confirmed",
			"code":         http.StatusInternalServerError,
			"failed_seats": failedSeats,
			"details":      err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":          true,
		"message":          "Seats confirmed successfully",
		"confirmed_seats":  req.SeatIDs,
		"flight_id":        req.FlightID,
	})
}

// CheckSeatLock endpoint: GET /api/v1/seats/:flight_id/:seat_id/lock-status
func (h *SeatHandler) CheckSeatLock(c *gin.Context) {
	flightID := c.Param("flight_id")
	seatID := c.Param("seat_id")

	isLocked, lockedBy, err := h.seatService.CheckSeatLock(flightID, seatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to check seat lock",
			"code":    http.StatusInternalServerError,
			"details": err.Error(),
		})
		return
	}

	if !isLocked {
		c.JSON(http.StatusOK, gin.H{
			"success":   true,
			"is_locked": false,
			"message":   "Seat is available",
		})
		return
	}

	// Get TTL
	ttl, _ := h.seatService.GetSeatLockTTL(flightID, seatID)

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"is_locked":  true,
		"locked_by":  lockedBy,
		"ttl_seconds": int64(ttl.Seconds()),
	})
}

// LockMultipleSeats endpoint: POST /api/v1/seats/lock-multiple
type LockMultipleSeatsRequest struct {
	FlightID string   `json:"flight_id" binding:"required"`
	SeatIDs  []string `json:"seat_ids" binding:"required,min=1"`
	UserID   string   `json:"user_id" binding:"required"`
}

func (h *SeatHandler) LockMultipleSeats(c *gin.Context) {
	var req LockMultipleSeatsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
			"code":    http.StatusBadRequest,
		})
		return
	}

	locked, lockedSeats, err := h.seatService.LockMultipleSeats(req.FlightID, req.SeatIDs, req.UserID)
	if err != nil {
		if errors.Is(err, services.ErrSeatAlreadyLocked) {
			c.JSON(http.StatusConflict, gin.H{
				"success":     false,
				"message":     "Some seats are already locked",
				"code":        http.StatusConflict,
				"failed_seats": lockedSeats,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to lock seats",
			"code":    http.StatusInternalServerError,
			"details": err.Error(),
		})
		return
	}

	if !locked {
		c.JSON(http.StatusConflict, gin.H{
			"success":     false,
			"message":     "Some seats are already locked",
			"code":        http.StatusConflict,
			"failed_seats": lockedSeats,
		})
		return
	}

	expiryTime := time.Now().Add(services.DefaultLockTTL)

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"message":     "Seats locked successfully",
		"locked_seats": lockedSeats,
		"expiry_time": expiryTime,
		"ttl_seconds": int64(services.DefaultLockTTL.Seconds()),
	})
}
