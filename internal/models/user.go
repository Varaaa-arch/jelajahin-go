package models

import "time"

type User struct {
	ID               string    `gorm:"primaryKey;type:uuid" json:"id"`
	Email            string    `gorm:"uniqueIndex;type:varchar(255)" json:"email"`
	PasswordHash     string    `gorm:"type:varchar(255)" json:"-"`
	FirstName        string    `gorm:"type:varchar(100)" json:"first_name"`
	LastName         string    `gorm:"type:varchar(100)" json:"last_name"`
	PhoneNumber      string    `gorm:"type:varchar(20)" json:"phone_number"`
	DateOfBirth      *string   `gorm:"type:date" json:"date_of_birth"`
	Gender           string    `gorm:"type:varchar(10)" json:"gender"`
	IdentityType     string    `gorm:"type:varchar(50)" json:"identity_type"`
	IdentityNumber   string    `gorm:"uniqueIndex;type:varchar(50)" json:"identity_number"`
	City             string    `gorm:"type:varchar(100)" json:"city"`
	Address          string    `gorm:"type:text" json:"address"`
	Role             string    `gorm:"type:varchar(50);default:'customer'" json:"role"`
	IsEmailVerified  bool      `gorm:"default:false" json:"is_email_verified"`
	IsActive         bool      `gorm:"default:true" json:"is_active"`
	LastLogin        *time.Time `json:"last_login"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
