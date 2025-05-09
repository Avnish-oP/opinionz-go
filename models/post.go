package models

import (
	"time"
)

type Post struct {
	PostID    string    `gorm:"primaryKey" json:"post_id"`
	UserID    string    `json:"user_id,omitempty"` // optional if anonymous
	Content   string    `gorm:"not null" json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Upvotes   int       `json:"upvotes"`
	Downvotes int       `json:"downvotes"`
	Tags      []string  `gorm:"type:text[]" json:"tags"` // Use PostgreSQL array type
	Anonymous bool      `json:"anonymous"`
	Doodle    string    `json:"doodle,omitempty"` // optional
}
