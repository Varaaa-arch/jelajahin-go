package handlers

import (
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
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Try to lock seat
	locked, lockedBy, err := h.seatService.LockSeat(req.FlightID, req.SeatID, req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to lock seat",
			"details": err.Error(),
		})
		return
	}

	// If lock failed
	if !locked {
		c.JSON(http.StatusConflict, LockSeatResponse{
			Success:   false,
			Message:   "Seat already locked",
			LockedBy:  lockedBy,
		})
		return
	}

	// Lock successful - generate lock token
	lockToken := req.FlightID + ":" + req.SeatID + ":" + req.UserID
	expiryTime := time.Now().Add(15 * time.Minute)

	c.JSON(http.StatusOK, LockSeatResponse{
		Success:    true,
		Message:    "Seat locked successfully",
		LockToken:  lockToken,
		LockedBy:   req.UserID,
		ExpiryTime: expiryTime,
		TTLSeconds: 900, // 15 minutes = 900 seconds
	})
}

// UnlockSeat endpoint: POST /api/v1/seats/unlock
type UnlockSeatRequest struct {
	FlightID string `json:"flight_id" binding:"required"`
	SeatID   string `json:"seat_id" binding:"required"`
}

func (h *SeatHandler) UnlockSeat(c *gin.Context) {
	var req UnlockSeatRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	err := h.seatService.UnlockSeat(req.FlightID, req.SeatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to unlock seat",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Seat unlocked successfully",
	})
}

// CheckSeatLock endpoint: GET /api/v1/seats/:flight_id/:seat_id/lock-status
func (h *SeatHandler) CheckSeatLock(c *gin.Context) {
	flightID := c.Param("flight_id")
	seatID := c.Param("seat_id")

	isLocked, lockedBy, err := h.seatService.CheckSeatLock(flightID, seatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check seat lock",
		})
		return
	}

	if !isLocked {
		c.JSON(http.StatusOK, gin.H{
			"is_locked": false,
			"message":   "Seat is available",
		})
		return
	}

	// Get TTL
	ttl, _ := h.seatService.GetSeatLockTTL(flightID, seatID)

	c.JSON(http.StatusOK, gin.H{
		"is_locked": true,
		"locked_by": lockedBy,
		"ttl_seconds": ttl.Seconds(),
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
			"error": "Invalid request body",
		})
		return
	}

	locked, lockedSeats, err := h.seatService.LockMultipleSeats(req.FlightID, req.SeatIDs, req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to lock seats",
		})
		return
	}

	if !locked {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": "Some seats are already locked",
			"failed_seats": lockedSeats,
		})
		return
	}

	expiryTime := time.Now().Add(15 * time.Minute)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Seats locked successfully",
		"locked_seats": lockedSeats,
		"expiry_time": expiryTime,
		"ttl_seconds": 900,
	})
}
