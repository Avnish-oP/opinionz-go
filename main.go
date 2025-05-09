package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Avnish-oP/opinionz/config"
	"github.com/Avnish-oP/opinionz/models"
	"github.com/Avnish-oP/opinionz/routes"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()
	config.DB.AutoMigrate(&models.User{}, &models.Post{})

	r := routes.SetupRoutes()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}
	fmt.Println("Server is running on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, r))

}
