package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Avnish-oP/opinionz/config"
	"github.com/Avnish-oP/opinionz/middlewares"
	"github.com/Avnish-oP/opinionz/models"
	"github.com/google/uuid"
)

func CreateComment(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(middlewares.UserIDKey).(string)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	comment.UserID = userID
	comment.CreatedAt = time.Now()
	comment.ID = uuid.New().String()

	collection := config.MongoDB.Collection("comments")
	_, err := collection.InsertOne(r.Context(), comment)
	if err != nil {
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	response := models.Response{
		Message: "Comment created successfully",
		Success: true,
		Data:    comment,
	}
	json.NewEncoder(w).Encode(response)

}
