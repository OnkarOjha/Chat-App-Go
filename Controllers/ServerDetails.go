package Controllers

import (
	"encoding/json"
	"fmt"
	db "main/Database"
	models "main/Models"
	"net/http"
)

// fetch the room participant information by room_id
func ParticipantDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	EnableCors(&w)

	fmt.Println("we are fetching participant details from DB..")
	var mp = make(map[string]string)
	json.NewDecoder(r.Body).Decode(&mp)
	room_id := mp["room_id"]
	fmt.Println("room_id: ", room_id)

	var participants []models.Participant
	db.DB.Raw("select user_id, room_name from participants where room_id=?", room_id).Scan(&participants)

	json.NewEncoder(w).Encode(&participants)
}

// fetch the all the room information 
func RoomDetails(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	EnableCors(&w)

	fmt.Println("We are fetching room details from DB...")

	var mp = make(map[string]string)
	json.NewDecoder(r.Body).Decode(&mp)
	room_id := mp["room_id"]
	fmt.Println("room_id: ", room_id)

	var rooms []models.Room
	db.DB.Raw("Select * from rooms where room_id =?", room_id).Scan(&rooms)

	json.NewEncoder(w).Encode(&rooms)
}

//fetch all the message from DB with pagination
func MessageDetails(w http.ResponseWriter , r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	EnableCors(&w)
	var mp = make(map[string]string)
	json.NewDecoder(r.Body).Decode(&mp)
	room_id := mp["room_id"]

	fmt.Printf("We are fetching the messages of room %s" , room_id)

	var messages []models.Message

	db.DB.Raw("SELECT * FROM messages where room_id=? ORDER BY created_at  DESC LIMIT 10",room_id).Scan(&messages)

	json.NewEncoder(w).Encode(&messages)
}
