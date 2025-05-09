package models

import "time"

type User struct {
	UserID      string    `bson:"_id,omitempty" json:"user_id"` // MongoDB uses `_id` as the primary key
	Email       string    `bson:"email" json:"email"`
	Password    string    `bson:"password" json:"password"`
	Contact     string    `bson:"contact" json:"contact"`
	OTP         string    `bson:"otp" json:"otp"`
	OTPExpiry   time.Time `bson:"otp_expiry" json:"otp_expiry"`
	Verified    bool      `bson:"verified" json:"verified"`
	Role        string    `bson:"role" json:"role"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	ResetToken  string    `bson:"reset_token" json:"reset_token"`
	ResetExpiry time.Time `bson:"reset_expiry" json:"reset_expiry"`
	Doodle      string    `bson:"doodle" json:"doodle"`
	Interests   []string  `bson:"interests" json:"interests"`
}
