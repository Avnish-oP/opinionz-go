package models

import "time"

type Post struct {
	ID        string    `bson:"_id,omitempty" json:"post_id"` // MongoDB uses `_id` as the primary key
	UserID    string    `bson:"user_id" json:"user_id"`
	Content   string    `bson:"content" json:"content"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	Tags      []string  `bson:"tags" json:"tags"`
	Anonymous bool      `bson:"anonymous" json:"anonymous"`
	Doodle    string    `bson:"doodle" json:"doodle"`
	Images    []string  `bson:"images" json:"images"`
	Comments  []Comment `bson:"comments" json:"comments"`
	Upvotes   []string  `bson:"upvotes_user_ids" json:"upvotes_user_ids"`
	Downvotes []string  `bson:"downvotes_user_ids" json:"downvotes_user_ids"`
}
