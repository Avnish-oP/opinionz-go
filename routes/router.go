package routes

import (
	"fmt"

	"github.com/Avnish-oP/opinionz/controllers"
	"github.com/gorilla/mux"
)

// SetupRoutes exports the function to be accessible in main.go
func SetupRoutes() *mux.Router {
	fmt.Println("Setting up routes...")
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/register", controllers.Register).Methods("POST")
	r.HandleFunc("/api/v1/login", controllers.Login).Methods("POST")
	r.HandleFunc("/api/v1/logout", controllers.Logout).Methods("POST")
	return r
}
