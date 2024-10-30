package handlers

import (
	"GolangProject/db"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// ensureTableExists checks if the users table exists and creates it if it doesn't
func ensureTableExists() error {
	query := `
        CREATE TABLE IF NOT EXISTS users (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            name VARCHAR(100) NOT NULL,
            age INT NOT NULL
        )`

	_, err := db.DB.Exec(query)
	return err
}

// CreateUser handles creating a new user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	// First, ensure the table exists
	if err := ensureTableExists(); err != nil {
		http.Error(w, "Error creating table: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert the new user, letting PostgreSQL generate the UUID
	query := `INSERT INTO users (name, age) VALUES ($1, $2) RETURNING id, name, age`
	err := db.DB.QueryRow(query, user.Name, user.Age).Scan(&user.ID, &user.Name, &user.Age)
	if err != nil {
		log.Println("Error creating user:", err)
		http.Error(w, "Error creating user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// GetUserByID handles retrieving a user by their ID
func GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if _, err := uuid.Parse(id); err != nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}

	var user User
	query := `SELECT id, name, age FROM users WHERE id = $1`
	err := db.DB.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Age)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// UpdateUser handles updating an existing user
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if _, err := uuid.Parse(id); err != nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `UPDATE users SET name = $1, age = $2 WHERE id = $3 RETURNING id, name, age`
	err := db.DB.QueryRow(query, user.Name, user.Age, id).Scan(&user.ID, &user.Name, &user.Age)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// DeleteUser handles deleting a user by their ID
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if _, err := uuid.Parse(id); err != nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM users WHERE id = $1`
	result, err := db.DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetAllUsers handles retrieving all users
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	query := `SELECT id, name, age FROM users`
	rows, err := db.DB.Query(query)
	if err != nil {
		http.Error(w, "Error retrieving users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Age); err != nil {
			http.Error(w, "Error scanning users", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	json.NewEncoder(w).Encode(users)
}
