package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

type Contact struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

var users []Contact
var nextID int



func loadUsers() {
    data, err := os.ReadFile("contact.json")
	if err != nil {
        if os.IsNotExist(err) {
			fmt.Printf("File is fresh.", err)
			nextID = 1
            return
		} 
	}
	json.Unmarshal(data, &users)

	for _, user := range users {
		id, _ := strconv.Atoi(user.ID)
		if id >= nextID {
			nextID = id + 1
		}
	}

}

func saveUsers() {
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		fmt.Printf("Failed to marshal users")
	}

	os.WriteFile("contact.json", data, 0644)
}

func getContacts(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func getContact(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
    id := params["id"]

	for _, user := range users {
		if user.ID == id {
			json.NewEncoder(w).Encode(user)
			return
		}
	}
}

func deleteContacts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]

	for i, u := range users {
		if u.ID == id {
			users = append(users[:i], users[i+1:]...)
			saveUsers()
			json.NewEncoder(w).Encode(map[string]string{"message": "User deleted"})
            return
		}
	}

}

func updateContacts(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)
	id := params["id"]

	for i, u := range users {
		if u.ID == id {
			var updateUser Contact
			json.NewDecoder(r.Body).Decode(&updateUser)

			updateUser.ID = id
			users[i]= updateUser
			saveUsers()
			json.NewEncoder(w).Encode(updateUser)
			return
		}
	}
}

func createContacts(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
	var user Contact
	json.NewDecoder(r.Body).Decode(&user)

	for _, u := range users {
		if u.Name == user.Name && u.Email == user.Email {
			http.Error(w, "Username or Email already exists", http.StatusBadRequest)
			return
		}
	}

	user.ID = strconv.Itoa(nextID)
	nextID++
	users = append(users, user)
	saveUsers()
	json.NewEncoder(w).Encode(user)
}




func main() {
	loadUsers()

	r := mux.NewRouter()

	r.HandleFunc("/api/contacts", getContacts).Methods("GET")
	r.HandleFunc("/api/contacts/{id}", getContact).Methods("GET")
	r.HandleFunc("/api/contacts/{id}", deleteContacts).Methods("DELETE")
	r.HandleFunc("/api/contacts/{id}", updateContacts).Methods("PUT")
	r.HandleFunc("/api/contacts", createContacts).Methods("POST")

	http.ListenAndServe(":8080", r)
}