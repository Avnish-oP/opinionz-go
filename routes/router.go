package routes

import (
	"fmt"
	"net/http"

	"github.com/Avnish-oP/opinionz/controllers"
	"github.com/Avnish-oP/opinionz/middlewares"
	"github.com/gorilla/mux"
)

// SetupRoutes exports the function to be accessible in main.go
func SetupRoutes() *mux.Router {
	fmt.Println("Setting up routes...")
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/register", controllers.Register).Methods("POST")
	r.HandleFunc("/api/v1/login", controllers.Login).Methods("POST")
	r.HandleFunc("/api/v1/logout", controllers.Logout).Methods("POST")
	r.HandleFunc("/api/v1/verify-otp", controllers.Verify).Methods("POST")
	r.Handle("/api/v1/create-post", middlewares.AuthMiddleware(http.HandlerFunc(controllers.CreatePost))).Methods("POST")
	r.Handle("/api/v1/create-comment", middlewares.AuthMiddleware(http.HandlerFunc(controllers.CreateComment))).Methods("POST")
	r.Handle("/api/v1/vote", middlewares.AuthMiddleware(http.HandlerFunc(controllers.HandleVote))).Methods("POST")
	r.Handle("/api/v1/recommended-posts", middlewares.AuthMiddleware(http.HandlerFunc(controllers.GetRecommendedPosts))).Methods("GET")
	return r
}
