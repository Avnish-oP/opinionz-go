package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Avnish-oP/opinionz/config"
	"github.com/Avnish-oP/opinionz/routes"
)

func main() {
	config.LoadEnv()
	config.ConnectMongoDB() // Connect to MongoDB

	r := routes.SetupRoutes()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}
	fmt.Println("Server is running on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
