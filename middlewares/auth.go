package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Avnish-oP/opinionz/utils"
)

type contextKey string

const userIDKey contextKey = "userID"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Unauthorized from middleware", http.StatusUnauthorized)
			return
		}
		tokenString := cookie.Value
		fmt.Println("Token from cookie:", tokenString)
		userID, err := utils.ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, "Unauthorized from middleware2", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		// Log the userID for debugging purposes
		fmt.Println("User ID from token:", ctx.Value(userIDKey))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
