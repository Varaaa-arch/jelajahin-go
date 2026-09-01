package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/Varaaa-arch/jelajahin-go/internal/services"
)

func TestLockSeatHandler(t *testing.T) {
	// Setup
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	seatService := services.NewSeatService(client)
	handler := NewSeatHandler(seatService)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/v1/seats/lock", handler.LockSeat)

	// Test request
	reqBody := LockSeatRequest{
		FlightID: "test-flight-1",
		SeatID:   "12A",
		UserID:   "test-user-1",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/v1/seats/lock", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp LockSeatResponse
	json.Unmarshal(w.Body.Bytes(), &resp)

	if !resp.Success {
		t.Errorf("Expected success to be true")
	}

	if resp.LockedBy != "test-user-1" {
		t.Errorf("Expected locked_by test-user-1, got %s", resp.LockedBy)
	}

	// Cleanup
	seatService.UnlockSeat("test-flight-1", "12A")
}

func TestLockSeatAlreadyLockedHandler(t *testing.T) {
	// Setup
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	seatService := services.NewSeatService(client)
	handler := NewSeatHandler(seatService)

	// Pre-lock seat
	seatService.LockSeat("test-flight-2", "12B", "user-1")

	// Setup router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/v1/seats/lock", handler.LockSeat)

	// Test request (different user)
	reqBody := LockSeatRequest{
		FlightID: "test-flight-2",
		SeatID:   "12B",
		UserID:   "user-2",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/v1/seats/lock", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should get 409 Conflict
	if w.Code != http.StatusConflict {
		t.Errorf("Expected status 409, got %d", w.Code)
	}

	var resp LockSeatResponse
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp.Success {
		t.Errorf("Expected success to be false")
	}

	if resp.LockedBy != "user-1" {
		t.Errorf("Expected locked_by user-1, got %s", resp.LockedBy)
	}

	// Cleanup
	seatService.UnlockSeat("test-flight-2", "12B")
}