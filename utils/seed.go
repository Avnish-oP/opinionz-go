package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/Avnish-oP/opinionz/config"
	"github.com/Avnish-oP/opinionz/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SeedDatabase() {
	ctx := context.Background()

	// Seed Users
	userCollection := config.MongoDB.Collection("users")
	var users []interface{}
	for i := 1; i <= 100; i++ {
		users = append(users, models.User{
			UserID:    primitive.NewObjectID().Hex(),
			Email:     fmt.Sprintf("user%d@example.com", i),
			Password:  fmt.Sprintf("password%d", i),
			Contact:   fmt.Sprintf("12345678%02d", i),
			Verified:  true,
			Role:      "user",
			CreatedAt: time.Now(),
			Interests: []string{"tech", "sports", "music", "travel"}[i%4:],
		})
	}
	userIDs := []string{}
	for _, user := range users {
		userIDs = append(userIDs, user.(models.User).UserID)
	}
	_, err := userCollection.InsertMany(ctx, users)
	if err != nil {
		fmt.Println("Error seeding users:", err)
	}

	// Seed Posts
	postCollection := config.MongoDB.Collection("posts")
	var posts []interface{}
	for i := 1; i <= 200; i++ {
		posts = append(posts, models.Post{
			ID:        primitive.NewObjectID().Hex(),
			UserID:    userIDs[i%len(userIDs)],
			Content:   fmt.Sprintf("This is post number %d", i),
			CreatedAt: time.Now(),
			Tags:      []string{"tech", "sports", "music", "travel"}[i%4:],
			Upvotes:   []string{},
			Downvotes: []string{},
		})
	}
	postIDs := []string{}
	for _, post := range posts {
		postIDs = append(postIDs, post.(models.Post).ID)
	}
	_, err = postCollection.InsertMany(ctx, posts)
	if err != nil {
		fmt.Println("Error seeding posts:", err)
	}

	// Seed Comments
	commentCollection := config.MongoDB.Collection("comments")
	var comments []interface{}
	for i := 1; i <= 500; i++ {
		comments = append(comments, models.Comment{
			ID:        primitive.NewObjectID().Hex(),
			UserID:    userIDs[i%len(userIDs)],
			Content:   fmt.Sprintf("This is comment number %d", i),
			CreatedAt: time.Now(),
			PostID:    postIDs[i%len(postIDs)],
			Upvotes:   0,
			Downvotes: 0,
		})
	}
	_, err = commentCollection.InsertMany(ctx, comments)
	if err != nil {
		fmt.Println("Error seeding comments:", err)
	}

	fmt.Println("Database seeded successfully!")
}
