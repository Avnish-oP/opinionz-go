package models

import "time"

type Comment struct {
	ID        string    `bson:"_id,omitempty" json:"comment_id"` // MongoDB uses `_id` as the primary key
	UserID    string    `bson:"user_id" json:"user_id"`
	Content   string    `bson:"content" json:"content"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	Upvotes   int       `bson:"upvotes" json:"upvotes"`
	Downvotes int       `bson:"downvotes" json:"downvotes"`
	PostID    string    `bson:"post_id" json:"post_id"`
}
