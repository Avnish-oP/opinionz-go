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
	r.HandleFunc("/api/v1/register", controllers.Register).Methods("POST")                                                //done
	r.HandleFunc("/api/v1/login", controllers.Login).Methods("POST")                                                      // done
	r.HandleFunc("/api/v1/logout", controllers.Logout).Methods("POST")                                                    // done
	r.HandleFunc("/api/v1/verify-otp", controllers.Verify).Methods("POST")                                                // done
	r.Handle("/api/v1/create-post", middlewares.AuthMiddleware(http.HandlerFunc(controllers.CreatePost))).Methods("POST") // done
	r.Handle("/api/v1/create-comment", middlewares.AuthMiddleware(http.HandlerFunc(controllers.CreateComment))).Methods("POST") // done
	r.Handle("/api/v1/vote", middlewares.AuthMiddleware(http.HandlerFunc(controllers.HandleVote))).Methods("POST")                      // done
	r.Handle("/api/v1/recommended-posts", middlewares.AuthMiddleware(http.HandlerFunc(controllers.GetRecommendedPosts))).Methods("GET") // done
	r.Handle("/api/v1/user-profile", middlewares.AuthMiddleware(http.HandlerFunc(controllers.GetUserProfile))).Methods("GET")           // done
	r.Handle("/api/v1/user-profile", middlewares.AuthMiddleware(http.HandlerFunc(controllers.UpdateUserProfile))).Methods("PUT")        // done
	r.Handle("/api/v1/view-post/{id}", middlewares.AuthMiddleware(http.HandlerFunc(controllers.ViewPost))).Methods("GET")
	// r.HandleFunc("/api/v1/view/{id}", controllers.ViewPost).Methods("GET")
	r.Handle("/api/v1/admin/users", middlewares.AdminMiddleware(http.HandlerFunc(controllers.GetAllUsers))).Methods("GET")
	r.Handle("/api/v1/admin/posts", middlewares.AdminMiddleware(http.HandlerFunc(controllers.GetAllPosts))).Methods("GET")
	r.Handle("/api/v1/admin/comments", middlewares.AdminMiddleware(http.HandlerFunc(controllers.GetAllComments))).Methods("GET")
	r.Handle("/api/v1/admin/delete-user", middlewares.AdminMiddleware(http.HandlerFunc(controllers.DeleteUser))).Methods("DELETE")
	r.Handle("/api/v1/admin/delete-post", middlewares.AdminMiddleware(http.HandlerFunc(controllers.DeletePost))).Methods("DELETE")
	r.Handle("/api/v1/admin/delete-comment", middlewares.AdminMiddleware(http.HandlerFunc(controllers.DeleteComment))).Methods("DELETE")

	return r
}
