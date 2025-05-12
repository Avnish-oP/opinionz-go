package controllers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Avnish-oP/opinionz/config"
	"github.com/Avnish-oP/opinionz/middlewares"
	"github.com/Avnish-oP/opinionz/models"
	"github.com/google/uuid"
)
func CreateComment(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value(middlewares.UserIDKey)
	userID, ok := userIDRaw.(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	defer r.Body.Close()

	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(comment.Content) == "" || strings.TrimSpace(comment.PostID) == "" {
		http.Error(w, "Content and PostID are required", http.StatusBadRequest)
		return
	}

	comment.UserID = userID
	comment.CreatedAt = time.Now()
	comment.ID = uuid.New().String()
	comment.Upvotes = 0
	comment.Downvotes = 0

	collection := config.MongoDB.Collection("comments")
	_, err := collection.InsertOne(r.Context(), comment)
	if err != nil {
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := models.Response{
		Message: "Comment created successfully",
		Success: true,
		Data:    comment,
	}
	json.NewEncoder(w).Encode(response)
}
