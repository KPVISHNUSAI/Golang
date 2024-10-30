package routes

import (
	"GolangProject/handlers"
	"GolangProject/middleware"

	"github.com/gorilla/mux"
)

// RegisterUserRoutes sets up routes for user-related endpoints
func RegisterUserRoutes(router *mux.Router) {
	userRouter := router.PathPrefix("/api/users").Subrouter()
	userRouter.Use(middleware.JWTMiddleware) // Apply JWT middleware to user routes

	userRouter.HandleFunc("", handlers.GetAllUsers).Methods("GET")        // Get all users
	userRouter.HandleFunc("/{id}", handlers.GetUserByID).Methods("GET")   // Get user by ID
	userRouter.HandleFunc("", handlers.CreateUser).Methods("POST")        // Create a new user
	userRouter.HandleFunc("/{id}", handlers.UpdateUser).Methods("PUT")    // Update existing user
	userRouter.HandleFunc("/{id}", handlers.DeleteUser).Methods("DELETE") // Delete user by ID
}
