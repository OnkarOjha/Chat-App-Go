package Controller

import (
	"encoding/json"
	"fmt"
	db "main/Database"
	models "main/Models"
	"net/http"
)

// endpoint for setting room topic
func RoomTopicController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	fmt.Println("We are setting Room Topics...")

	var topic models.Topic
	json.NewDecoder(r.Body).Decode(&topic)

	db.DB.Create(&topic)
	json.NewEncoder(w).Encode(&topic)
}

// endpoint for getting room topics
func RoomTopicGetter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	fmt.Println("You are viewing Room Topics...")
	var topics []models.Topic
	query := "SELECT * FROM topics"
	db.DB.Raw(query).Scan(&topics)
	json.NewEncoder(w).Encode(&topics)
}
