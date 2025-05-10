package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Avnish-oP/opinionz/utils"
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
