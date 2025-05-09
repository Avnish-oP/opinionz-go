package models

import "time"

type Post struct {
	ID        string    `bson:"_id,omitempty" json:"post_id"` // MongoDB uses `_id` as the primary key
	UserID    string    `bson:"user_id" json:"user_id"`
	Content   string    `bson:"content" json:"content"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	Upvotes   int       `bson:"upvotes" json:"upvotes"`
	Downvotes int       `bson:"downvotes" json:"downvotes"`
	Tags      []string  `bson:"tags" json:"tags"`
	Anonymous bool      `bson:"anonymous" json:"anonymous"`
	Doodle    string    `bson:"doodle" json:"doodle"`
}
