package models

import "time"

type Booking struct {
	ID              string     `gorm:"primaryKey;type:uuid" json:"id"`
	PnrCode         string     `gorm:"uniqueIndex;type:varchar(10)" json:"pnr_code"`
	UserID          string     `gorm:"type:uuid" json:"user_id"`
	FlightID        string     `gorm:"type:uuid" json:"flight_id"`
	PromoID         *string    `gorm:"type:uuid" json:"promo_id"`
	BaseAmount      float64    `gorm:"type:decimal(12,2)" json:"base_amount"`
	DiscountAmount  float64    `gorm:"type:decimal(12,2);default:0" json:"discount_amount"`
	TaxAmount       float64    `gorm:"type:decimal(12,2);default:0" json:"tax_amount"`
	TotalPrice      float64    `gorm:"type:decimal(12,2)" json:"total_price"`
	PassengerCount  int        `json:"passenger_count"`
	Status          string     `gorm:"type:varchar(50);default:'pending'" json:"status"`
	BookingDate     time.Time  `json:"booking_date"`
	ExpirationDate  *time.Time `json:"expiration_date"`
	ConfirmedAt     *time.Time `json:"confirmed_at"`
	CancelledAt     *time.Time `json:"cancelled_at"`
	SpecialRequests string     `gorm:"type:text" json:"special_requests"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type Passenger struct {
	ID             string     `gorm:"primaryKey;type:uuid" json:"id"`
	BookingID      string     `gorm:"type:uuid" json:"booking_id"`
	Title          string     `gorm:"type:varchar(10)" json:"title"`
	FirstName      string     `gorm:"type:varchar(100)" json:"first_name"`
	LastName       string     `gorm:"type:varchar(100)" json:"last_name"`
	DateOfBirth    string     `gorm:"type:date" json:"date_of_birth"`
	Gender         string     `gorm:"type:varchar(10)" json:"gender"`
	IdentityType   string     `gorm:"type:varchar(50)" json:"identity_type"`
	IdentityNumber string     `gorm:"type:varchar(50)" json:"identity_number"`
	Nationality    string     `gorm:"type:varchar(100)" json:"nationality"`
	PassportNumber string     `gorm:"type:varchar(50)" json:"passport_number"`
	PassportExpiry *string    `gorm:"type:date" json:"passport_expiry"`
	FlightSeatID   *string    `gorm:"type:uuid" json:"flight_seat_id"`
	CheckInStatus  string     `gorm:"type:varchar(50);default:'not_checked_in'" json:"check_in_status"`
	CheckedInAt    *time.Time `json:"checked_in_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}
