package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/Avnish-oP/opinionz/config"
	"github.com/Avnish-oP/opinionz/middlewares"
	"github.com/Avnish-oP/opinionz/models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type PostResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
	Data    any    `json:"data"`
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
		return
	}

	// Extract userID from the context
	userID := r.Context().Value(middlewares.UserIDKey).(string)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userIDStr := userID
	fmt.Println("User ID from post token:", userIDStr)

	var input models.Post
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	// Set the userID for the post
	input.UserID = userIDStr
	input.ID = uuid.New().String()
	input.CreatedAt = time.Now()

	// Save the post to MongoDB
	collection := config.MongoDB.Collection("posts")

	_, err := collection.InsertOne(r.Context(), input)
	if err != nil {
		response := PostResponse{
			Message: "Error creating post",
			Success: false,
			Data:    nil,
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Respond with success
	response := PostResponse{
		Message: "Post created successfully",
		Success: true,
		Data:    input,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func GetRecommendedPosts(w http.ResponseWriter, r *http.Request) {
	// Extract userID from the context
	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Fetch user details to get interests
	userCollection := config.MongoDB.Collection("users")
	var user models.User
	err := userCollection.FindOne(r.Context(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		http.Error(w, "Failed to fetch user details", http.StatusInternalServerError)
		return
	}

	// Fetch posts based on user interests
	postCollection := config.MongoDB.Collection("posts")
	var posts []models.Post
	filter := bson.M{"tags": bson.M{"$in": user.Interests}} // Match posts with tags in user's interests
	cursor, err := postCollection.Find(r.Context(), filter)
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())

	for cursor.Next(r.Context()) {
		var post models.Post
		if err := cursor.Decode(&post); err != nil {
			continue
		}
		posts = append(posts, post)
	}

	if err := cursor.Err(); err != nil {
		http.Error(w, "Error iterating posts", http.StatusInternalServerError)
		return
	}

	// Sort posts by a recommendation algorithm (e.g., based on upvotes and recency)
	sortPostsByRecommendation(&posts)

	// Respond with recommended posts
	response := models.Response{
		Message: "Recommended posts fetched successfully",
		Success: true,
		Data:    posts,
	}
	json.NewEncoder(w).Encode(response)
}

func sortPostsByRecommendation(posts *[]models.Post) {
	sort.Slice(*posts, func(i, j int) bool {
		// Example: Higher upvotes and more recent posts are prioritized
		return len((*posts)[i].Upvotes)*2+len((*posts)[i].Downvotes) > len((*posts)[j].Upvotes)*2+len((*posts)[j].Downvotes)
	})
}
