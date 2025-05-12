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
		fmt.Println("Error decoding JSON:", err)
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	// Check if the user already exists in MongoDB
	collection := config.MongoDB.Collection("users")
	var existingUser models.User
	err := collection.FindOne(r.Context(), map[string]interface{}{"email": input.Email}).Decode(&existingUser)
	if err == nil {
		response.Message = "User already exists"
		response.Success = false
		response.Data = nil
		json.NewEncoder(w).Encode(response)
		return
	}

	// Hash the password and generate OTP
	hashedPassword, _ := utils.HashPassword(input.Password)
	otp := utils.GenerateOTP()
	fmt.Println("Generated OTP:", otp)

	user := models.User{
		UserID:    uuid.New().String(),
		Email:     input.Email,
		Password:  hashedPassword,
		Contact:   input.Contact,
		Role:      "user",
		OTP:       otp,
		OTPExpiry: time.Now().Add(15 * time.Minute),
		CreatedAt: time.Now(),
		Interests: input.Interests,
	}

	// Insert the user into MongoDB
	_, err = collection.InsertOne(r.Context(), user)
	if err != nil {
		response.Message = "Error creating user"
		response.Success = false
		response.Data = nil
		json.NewEncoder(w).Encode(response)
		return
	}

	// Send OTP via email
	go utils.SendEmail(user.Email, otp)

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

	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	fmt.Println("Login attempt for email:", input.Email)
	// Query MongoDB for the user
	collection := config.MongoDB.Collection("users")
	var user models.User
	err := collection.FindOne(r.Context(), map[string]interface{}{"email": input.Email}).Decode(&user)
	fmt.Println("User found:", user)
	fmt.Println("Error:", err)
	if err != nil {
		response.Message = "User not found"
		response.Success = false
		response.Data = nil
		json.NewEncoder(w).Encode(response)
		return
	}

	// Check if the user is verified
	if !user.Verified {
		response.Message = "User not verified"
		response.Success = false
		response.Data = nil
		json.NewEncoder(w).Encode(response)
		return
	}

	// Validate the password
	if !utils.CheckPasswordHash(input.Password, user.Password) {
		response.Message = "Invalid password"
		response.Success = false
		response.Data = nil
		json.NewEncoder(w).Encode(response)
		return
	}

	// Generate a JWT token
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

	// Clear the token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Unix(0, 0), // Expire immediately
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Message: "Logout successful",
		Success: true,
		Data:    nil,
	})
}

func Verify(w http.ResponseWriter, r *http.Request) {
	var response Response
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	collection := config.MongoDB.Collection("users")

	var user models.User
	err := collection.FindOne(r.Context(), map[string]any{"email": input.Email}).Decode(&user)
	if err != nil {
		response.Message = "User not found"
		response.Success = false
		response.Data = nil
		json.NewEncoder(w).Encode(response)
		return
	}

	// Check if OTP is valid
	if user.OTP != input.OTP || time.Now().After(user.OTPExpiry) {
		response.Message = "Invalid or expired OTP"
		response.Success = false
		response.Data = nil
		json.NewEncoder(w).Encode(response)
		return
	}

	// Mark user as verified
	user.Verified = true
	user.OTP = "" // Clear OTP after verification
	user.OTPExpiry = time.Time{}
	_, err = collection.UpdateOne(r.Context(), map[string]any{"email": input.Email}, map[string]any{"$set": user})
	if err != nil {
		response.Message = "Error updating user"
		response.Success = false
		response.Data = nil
		json.NewEncoder(w).Encode(response)
		return
	}

	response.Message = "User verified successfully"
	response.Success = true
	response.Data = nil
	json.NewEncoder(w).Encode(response)
}
