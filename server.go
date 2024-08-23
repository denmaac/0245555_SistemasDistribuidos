package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

// User struct with json tags
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Global variables
var ID_counter int // Counter for user IDs
var users []User   // Slice of users
var mux sync.Mutex // Mutex for synchronization

// Handler functions
func addUser(w http.ResponseWriter, r *http.Request) {
	mux.Lock()
	defer mux.Unlock()

	// Decode the request body into a User struct
	var user_temp User
	err := json.NewDecoder(r.Body).Decode(&user_temp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Assign an ID to the user and increment the counter
	user_temp.ID = ID_counter
	ID_counter++
	users = append(users, user_temp)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user_temp)
}

// Handler function
func getUsers(w http.ResponseWriter, r *http.Request) {
	mux.Lock()
	defer mux.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func main() {
	ID_counter = 0
	users = []User{}

	http.HandleFunc("/addUser", addUser)
	http.HandleFunc("/getUsers", getUsers)

	fmt.Println("Server running on port 8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
	//http.ListenAndServe(":8080", nil)
}
