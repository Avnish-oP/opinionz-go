package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Avnish-oP/opinionz/config"
	"github.com/Avnish-oP/opinionz/models"
	"github.com/Avnish-oP/opinionz/utils"
	"github.com/google/uuid"
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

	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "cookie not found", http.StatusUnauthorized)
		return
	}
	tokenString := cookie.Value
	userID, err := utils.ValidateJWT(tokenString)
	if err != nil {
		http.Error(w, "can't get userid", http.StatusUnauthorized)
		return
	}

	fmt.Println("User ID from context:", userID)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var input models.Post
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	// Set the userID for the post
	input.UserID = userID
	input.PostID = uuid.New().String()
	input.CreatedAt = time.Now()

	// Save the post to the database
	if err := config.DB.Create(&input).Error; err != nil {
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
