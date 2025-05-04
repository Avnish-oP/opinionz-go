package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Avnish-oP/opinionz/config"
	"github.com/Avnish-oP/opinionz/models"
	"github.com/Avnish-oP/opinionz/utils"
	"github.com/google/uuid"
)

type Response struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
	Data    any    `json:"data"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	var response Response
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
		return
	}

	var input models.User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	hashedPassword, _ := utils.HashPassword(input.Password)
	user := models.User{
		UserID:    uuid.New().String(),
		Email:     input.Email,
		Password:  hashedPassword,
		Contact:   input.Contact,
		Role:      "user",
		CreatedAt: time.Now(),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		response.Message = "Error creating user"
		response.Success = false
		response.Data = nil
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusCreated)
	response.Message = "User created successfully"
	response.Success = true
	response.Data = user
	json.NewEncoder(w).Encode(response)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var response Response
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
		return
	}

	var input models.User
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		response.Message = "User not found"
		response.Success = false
		response.Data = nil
		json.NewEncoder(w).Encode(response)
		return
	}

	if !utils.CheckPasswordHash(input.Password, user.Password) {
		response.Message = "Invalid password"
		response.Success = false
		response.Data = nil
		json.NewEncoder(w).Encode(response)
		return
	}

	token, _ := utils.GenerateJWT(user.UserID)
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})
	response.Message = "Login successful"
	response.Success = true
	response.Data = user
	json.NewEncoder(w).Encode(response)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Unix(0, 0),
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Message: "Logout successful",
		Success: true,
		Data:    nil,
	})
}
