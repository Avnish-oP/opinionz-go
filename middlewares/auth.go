package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Avnish-oP/opinionz/config"
	"github.com/Avnish-oP/opinionz/models"
	"github.com/Avnish-oP/opinionz/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type contextKey string

const UserIDKey contextKey = "userID"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Unauthorized from middleware", http.StatusUnauthorized)
			return
		}
		tokenString := cookie.Value

		userID, err := utils.ValidateJWT(tokenString)
		if err != nil {
			fmt.Println("Error validating token:", err)
			http.Error(w, "Unauthorized from middleware2", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		// Log the userID for debugging purposes
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(UserIDKey).(string)
		if !ok || userID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Fetch user details and check role
		userCollection := config.MongoDB.Collection("users")
		var user models.User
		err := userCollection.FindOne(r.Context(), bson.M{"_id": userID}).Decode(&user)
		if err != nil || user.Role != "admin" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
