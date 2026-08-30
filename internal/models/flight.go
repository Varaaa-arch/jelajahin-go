package models

import (
	"time"
)

type Airport struct {
	ID       string `gorm:"primaryKey;type:uuid" json:"id"`
	Code     string `gorm:"uniqueIndex;type:varchar(3)" json:"code"`
	Name     string `gorm:"type:varchar(255)" json:"name"`
	City     string `gorm:"type:varchar(100)" json:"city"`
	Country  string `gorm:"type:varchar(100)" json:"country"`
	Timezone string `gorm:"type:varchar(50)" json:"timezone"`
	IsActive bool   `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Airline struct {
	ID             string `gorm:"primaryKey;type:uuid" json:"id"`
	Code           string `gorm:"uniqueIndex;type:varchar(3)" json:"code"`
	Name           string `gorm:"type:varchar(255)" json:"name"`
	LogoURL        string `gorm:"type:text" json:"logo_url"`
	HeadquartersCity string `gorm:"type:varchar(100)" json:"headquarters_city"`
	IsActive       bool   `gorm:"default:true" json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type AircraftType struct {
	ID           string `gorm:"primaryKey;type:uuid" json:"id"`
	Name         string `gorm:"uniqueIndex;type:varchar(100)" json:"name"`
	Manufacturer string `gorm:"type:varchar(100)" json:"manufacturer"`
	Model        string `gorm:"type:varchar(100)" json:"model"`
	TotalSeats   int    `json:"total_seats"`
	CreatedAt    time.Time `json:"created_at"`
}

type Aircraft struct {
	ID                  string `gorm:"primaryKey;type:uuid" json:"id"`
	AircraftTypeID      string `gorm:"type:uuid" json:"aircraft_type_id"`
	AirlineID           string `gorm:"type:uuid" json:"airline_id"`
	RegistrationNumber  string `gorm:"uniqueIndex;type:varchar(20)" json:"registration_number"`
	ManufactureYear     int    `json:"manufacture_year"`
	LastMaintenanceDate time.Time `json:"last_maintenance_date"`
	IsActive            bool   `gorm:"default:true" json:"is_active"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type Route struct {
	ID                   string `gorm:"primaryKey;type:uuid" json:"id"`
	AirlineID            string `gorm:"type:uuid" json:"airline_id"`
	OriginAirportID      string `gorm:"type:uuid" json:"origin_airport_id"`
	DestinationAirportID string `gorm:"type:uuid" json:"destination_airport_id"`
	FlightNumberPrefix   string `gorm:"type:varchar(3)" json:"flight_number_prefix"`
	DistanceKm           int    `json:"distance_km"`
	EstimatedDurationMinutes int `json:"estimated_duration_minutes"`
	IsActive             bool   `gorm:"default:true" json:"is_active"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type Flight struct {
	ID             string    `gorm:"primaryKey;type:uuid" json:"id"`
	RouteID        string    `gorm:"type:uuid" json:"route_id"`
	AircraftID     string    `gorm:"type:uuid" json:"aircraft_id"`
	FlightNumber   string    `gorm:"uniqueIndex:,composite:departure_date;type:varchar(10)" json:"flight_number"`
	DepartureDate  string    `gorm:"uniqueIndex:,composite:flight_number;type:date" json:"departure_date"`
	DepartureTime  string    `gorm:"type:time" json:"departure_time"`
	ArrivalTime    string    `gorm:"type:time" json:"arrival_time"`
	BasePrice      float64   `gorm:"type:decimal(12,2)" json:"base_price"`
	TaxSurcharge   float64   `gorm:"type:decimal(12,2);default:0" json:"tax_surcharge"`
	FuelSurcharge  float64   `gorm:"type:decimal(12,2);default:0" json:"fuel_surcharge"`
	Status         string    `gorm:"type:varchar(50);default:'scheduled'" json:"status"`
	SeatsAvailable int       `json:"seats_available"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type SeatClass struct {
	ID                   string `gorm:"primaryKey;type:uuid" json:"id"`
	Name                 string `gorm:"uniqueIndex;type:varchar(50)" json:"name"`
	DisplayName          string `gorm:"type:varchar(100)" json:"display_name"`
	BaggageAllowanceKg   int    `gorm:"default:20" json:"baggage_allowance_kg"`
	CarryOnAllowanceKg   int    `gorm:"default:7" json:"carry_on_allowance_kg"`
	CreatedAt            time.Time `json:"created_at"`
}

type AircraftSeat struct {
	ID               string `gorm:"primaryKey;type:uuid" json:"id"`
	AircraftTypeID   string `gorm:"type:uuid" json:"aircraft_type_id"`
	SeatClassID      string `gorm:"type:uuid" json:"seat_class_id"`
	SeatNumber       string `gorm:"type:varchar(10)" json:"seat_number"`
	RowNumber        int    `json:"row_number"`
	ColumnLetter     string `gorm:"type:varchar(1)" json:"column_letter"`
	IsExitRow        bool   `gorm:"default:false" json:"is_exit_row"`
	IsExtraLegroom   bool   `gorm:"default:false" json:"is_extra_legroom"`
	CreatedAt        time.Time `json:"created_at"`
}

type FlightSeat struct {
	ID             string `gorm:"primaryKey;type:uuid" json:"id"`
	FlightID       string `gorm:"type:uuid" json:"flight_id"`
	AircraftSeatID string `gorm:"type:uuid" json:"aircraft_seat_id"`
	CurrentPrice   float64 `gorm:"type:decimal(12,2)" json:"current_price"`
	IsAvailable    bool   `gorm:"default:true" json:"is_available"`
	BookingID      *string `gorm:"type:uuid" json:"booking_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
