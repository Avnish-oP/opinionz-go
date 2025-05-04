package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID          uint
	UserID      string `gorm:"primaryKey" json:"user_id"`
	Email       string `gorm:"uniqueIndex;not null" json:"email"`
	Password    string `gorm:"not null" json:"password"`
	Contact     string `gorm:"not null" json:"contact"`
	OTP         string
	OTPExpiry   time.Time
	Role        string
	CreatedAt   time.Time
	ResetToken  string
	ResetExpiry time.Time
	Doodle      string
	Interests   []string `gorm:"type:text[]"`
}
