package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/Avnish-oP/opinionz/config"
	"github.com/Avnish-oP/opinionz/middlewares"
	"github.com/Avnish-oP/opinionz/models"
	"go.mongodb.org/mongo-driver/bson"
)

func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	// Extract userID from the context
	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Fetch user details
	userCollection := config.MongoDB.Collection("users")
	var user models.User
	err := userCollection.FindOne(r.Context(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		http.Error(w, "Failed to fetch user details", http.StatusInternalServerError)
		return
	}

	// Fetch user's posts with comments and votes
	postCollection := config.MongoDB.Collection("posts")
	var posts []map[string]interface{}
	cursor, err := postCollection.Find(r.Context(), bson.M{"user_id": userID})
	if err != nil {
		http.Error(w, "Failed to fetch user posts", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())

	for cursor.Next(r.Context()) {
		var post models.Post
		if err := cursor.Decode(&post); err != nil {
			continue
		}

		// Fetch comments for the post
		commentCollection := config.MongoDB.Collection("comments")
		var comments []models.Comment
		commentCursor, err := commentCollection.Find(r.Context(), bson.M{"post_id": post.ID})
		if err == nil {
			defer commentCursor.Close(r.Context())
			for commentCursor.Next(r.Context()) {
				var comment models.Comment
				if err := commentCursor.Decode(&comment); err == nil {
					comments = append(comments, comment)
				}
			}
		}

		// Add post data with comments and votes
		postData := map[string]interface{}{
			"post":      post,
			"comments":  comments,
			"upvotes":   len(post.Upvotes),
			"downvotes": len(post.Downvotes),
		}
		posts = append(posts, postData)
	}

	if err := cursor.Err(); err != nil {
		http.Error(w, "Error iterating user posts", http.StatusInternalServerError)
		return
	}

	// Respond with user profile
	response := models.Response{
		Message: "User profile fetched successfully",
		Success: true,
		Data: map[string]interface{}{
			"user":  user,
			"posts": posts,
		},
	}
	json.NewEncoder(w).Encode(response)
}
