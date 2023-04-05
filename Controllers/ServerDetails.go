package Controllers

import (
	"encoding/json"
	"fmt"
	db "main/Database"
	models "main/Models"
	"net/http"
	response "main/Response"
)

// fetch the room participant information by room_id
func ParticipantDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	EnableCors(&w)

	fmt.Println("we are fetching participant details from DB..")
	var mp = make(map[string]interface{})
	json.NewDecoder(r.Body).Decode(&mp)
	room_id , ok := mp["room_id"]
	if !ok{	
		response.ShowResponse(
			"Bad Request",
			400,
			"Unable to fetch room_id or invalid room_id",
			"",
			w,
		)
		return
	}

	var participants []models.Participant
	db.DB.Raw("select * from participants where room_id=?", room_id).Scan(&participants)

	response.ShowResponse(
		"OK",
		200,
		"User information for this room",
		participants,
		w,
	)
}

// fetch the all the room information 
func RoomDetails(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	EnableCors(&w)

	fmt.Println("We are fetching room details from DB...")

	var mp = make(map[string]interface{})
	json.NewDecoder(r.Body).Decode(&mp)
	room_id,ok := mp["room_id"]
	if !ok{
		response.ShowResponse(
			"Bad Request",
			400,
			"Unable to fetch room_id or invalid room_id",
			"",
			w,
		)
		return
	}

	var rooms []models.Room
	db.DB.Raw("Select * from rooms where room_id =?", room_id).Scan(&rooms)

	
	response.ShowResponse(
		"OK",
		200,
		"Room Details",
		rooms,
		w,
	)
}

//fetch all the message from DB with limit
func MessageDetails(w http.ResponseWriter , r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	EnableCors(&w)
	var mp = make(map[string]interface{})
	json.NewDecoder(r.Body).Decode(&mp)
	room_id,ok := mp["room_id"]
	if !ok{
		response.ShowResponse(
			"Bad Request",
			400,
			"Unable to fetch room_id or invalid room_id",
			"",
			w,
		)
		return
	}

	var messages []models.Message

	db.DB.Raw("SELECT * FROM messages where room_id=? ORDER BY created_at  DESC LIMIT 10",room_id).Scan(&messages)

	response.ShowResponse(
		"OK",
		200,
		"All the messages in this room have been shown below",
		messages,
		w,
	)
}

//Handler to fetch all the details of the user provided that in which room it is present
func UserRoomsDetails(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	EnableCors(&w)

	var mp = make(map[string]interface{})
	json.NewDecoder(r.Body).Decode(&mp)
	user_id , ok := mp["user_id"]
	if !ok{
		response.ShowResponse(
			"Bad Request",
			400,
			"Unable to fetch user_id or invalid user_id",
			"",
			w,
		)
		return
	}	

	var userRoomInformation []models.Participant	

	db.DB.Raw("select * from participants where user_id = ?",user_id).Scan(&userRoomInformation)

	response.ShowResponse(
		"OK",
		200,
		"User has access to following rooms",
		userRoomInformation,
		w,
	)
}