// main.go
package main

import (
	"GolangProject/db"
	"GolangProject/handlers"
	"GolangProject/middleware"
	"GolangProject/routes"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize database connection
	if err := db.Initialize(); err != nil {
		log.Fatal("Could not initialize database:", err)
	}
	defer db.Close()

	// Create a new router
	router := mux.NewRouter()

	// Register user routes with JWT middleware
	routes.RegisterUserRoutes(router)

	// Unprotected routes (for login, etc.)
	router.HandleFunc("/api/login", handlers.Login).Methods("POST") // Add login route

	// Protected user routes
	router.Use(middleware.JWTMiddleware) // Use JWT middleware

	// Start the server
	log.Println("Starting server on :9000")
	log.Fatal(http.ListenAndServe(":9000", router))
}
