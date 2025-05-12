package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Avnish-oP/opinionz/config"
	"github.com/Avnish-oP/opinionz/middlewares"
	"github.com/Avnish-oP/opinionz/models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

// func GetUserProfile(w http.ResponseWriter, r *http.Request) {
// 	// Extract userID from the context
// 	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
// 	if !ok || userID == "" {
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	// Fetch user details
// 	userCollection := config.MongoDB.Collection("users")
// 	var user models.User
// 	err := userCollection.FindOne(r.Context(), bson.M{"_id": userID}).Decode(&user)
// 	if err != nil {
// 		http.Error(w, "Failed to fetch user details", http.StatusInternalServerError)
// 		return
// 	}

// 	// Fetch user's posts with comments and votes
// 	postCollection := config.MongoDB.Collection("posts")
// 	var posts []map[string]interface{}
// 	cursor, err := postCollection.Find(r.Context(), bson.M{"user_id": userID})
// 	if err != nil {
// 		http.Error(w, "Failed to fetch user posts", http.StatusInternalServerError)
// 		return
// 	}
// 	defer cursor.Close(r.Context())

// 	for cursor.Next(r.Context()) {
// 		var post models.Post
// 		if err := cursor.Decode(&post); err != nil {
// 			continue
// 		}

// 		// Fetch comments for the post
// 		commentCollection := config.MongoDB.Collection("comments")
// 		var comments []models.Comment
// 		commentCursor, err := commentCollection.Find(r.Context(), bson.M{"post_id": post.ID})
// 		if err == nil {
// 			defer commentCursor.Close(r.Context())
// 			for commentCursor.Next(r.Context()) {
// 				var comment models.Comment
// 				if err := commentCursor.Decode(&comment); err == nil {
// 					comments = append(comments, comment)
// 				}
// 			}
// 		}

// 		// Add post data with comments and votes
// 		postData := map[string]interface{}{
// 			"post":      post,
// 			"comments":  comments,
// 			"upvotes":   len(post.Upvotes),
// 			"downvotes": len(post.Downvotes),
// 		}
// 		posts = append(posts, postData)
// 	}

// 	if err := cursor.Err(); err != nil {
// 		http.Error(w, "Error iterating user posts", http.StatusInternalServerError)
// 		return
// 	}

// 	// Respond with user profile
// 	response := models.Response{
// 		Message: "User profile fetched successfully",
// 		Success: true,
// 		Data: map[string]interface{}{
// 			"user":  user,
// 			"posts": posts,
// 		},
// 	}
// 	json.NewEncoder(w).Encode(response)
// }

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

	// Respond with user profile, including images
	response := models.Response{
		Message: "User profile fetched successfully",
		Success: true,
		Data: map[string]interface{}{
			"user":  user,
			"posts": posts,
		},
	}

	// Set the content-type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Encode and send the response
	json.NewEncoder(w).Encode(response)
}

func UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	// Ensure method is PUT
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from context
	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// No need to convert to ObjectID since userID is a UUID string
	// Validate UUID format (optional)
	_, err := uuid.Parse(userID)
	if err != nil {
		log.Println("Invalid UUID format:", userID)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Parse request body to extract update data
	var updateData struct {
		Email     string   `json:"email"`
		Contact   string   `json:"contact"`
		Interests []string `json:"interests"` // Interests as an array of strings
		Doodle    string   `json:"doodle"`
	}

	// Decode the request body into updateData
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		log.Println("Failed to decode request body:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate interests (ensure it's an array of strings)
	if updateData.Interests == nil {
		updateData.Interests = []string{} // Ensure interests is never nil, default to empty array
	}

	// Build the update document for MongoDB
	update := bson.M{
		"$set": bson.M{
			"email":     updateData.Email,
			"contact":   updateData.Contact,
			"interests": updateData.Interests,
			"doodle":    updateData.Doodle,
		},
	}

	// Access user collection
	userCollection := config.MongoDB.Collection("users")
	filter := bson.M{"_id": userID} // Using UUID string directly as _id

	// Perform the update operation in MongoDB
	result, err := userCollection.UpdateOne(r.Context(), filter, update)
	if err != nil {
		log.Println("Error updating user:", err)
		http.Error(w, "Database error during update", http.StatusInternalServerError)
		return
	}

	// Check if the user was found and updated
	if result.MatchedCount == 0 {
		log.Println("User not found for ID:", userID)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Fetch the updated user data from the database
	var updatedUser models.User
	if err := userCollection.FindOne(r.Context(), filter).Decode(&updatedUser); err != nil {
		log.Println("Failed to fetch updated user:", err)
		http.Error(w, "Failed to fetch updated user", http.StatusInternalServerError)
		return
	}

	// Respond with the updated user data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.Response{
		Message: "User profile updated successfully",
		Success: true,
		Data:    updatedUser,
	})
}
