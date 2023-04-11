package Controllers

import (
	"encoding/json"
	"fmt"
	db "main/Database"
	models "main/Models"
	"net/http"
	response "main/Response"
	commonFunctions "main/Utils"
)

// Give me RoomId , i will give you how many users are there in the room
func ParticipantDetails(w http.ResponseWriter, r *http.Request) {
	commonFunctions.SetHeader(w)
	commonFunctions.EnableCors(&w)

	fmt.Println("we are fetching participant details from DB..")
	var mp = make(map[string]interface{})
	json.NewDecoder(r.Body).Decode(&mp)
	roomId , ok := mp["roomId"]
	if !ok{	
		response.ShowResponse(
			"Failure",
			400,
			"Unable to fetch room_id or invalid room_id",
			"",
			w,
		)
		return
	}

	var exists bool
	db.DB.Raw("SELECT EXISTS(select * from rooms where room_id = ? and  user_count=0)",roomId).Scan(&exists)
	if exists{
		response.ShowResponse(
			"Failure",
			400,
			"This room does not have any participants",
			"",
			w,
		)
		return 
	}

	var participants []models.Participant
	db.DB.Raw("select * from participants where room_id=?", roomId).Scan(&participants)

	response.ShowResponse(
		"Success",
		200,
		"User information for this room",
		participants,
		w,
	)
}

// specifically room details dega ki admin kaun hai and user_count kya hai
func RoomDetails(w http.ResponseWriter, r *http.Request){
	commonFunctions.SetHeader(w)
	commonFunctions.EnableCors(&w)

	fmt.Println("We are fetching room details from DB...")

	var mp = make(map[string]interface{})
	json.NewDecoder(r.Body).Decode(&mp)
	roomId,ok := mp["roomId"]
	if roomId == ""{
		response.ShowResponse(
			"Failure",
			404,
			"Room Id not found",
			"",
			w,
		)
		return
	}
	if !ok{
		response.ShowResponse(
			"Failure",
			400,
			"Unable to fetch room_id or invalid room_id",
			"",
			w,
		)
		return
	}

	var rooms []models.Room
	fmt.Println("roomidd: ",roomId)
	var exists bool
	db.DB.Raw("SELECT EXISTS(SELECT * FROM rooms where room_id =?)",roomId).Scan(&exists)
	if !exists{
		response.ShowResponse(
			"Failure",
			500,
			"Record not found in DB",
			"",
			w,
		)
		return
	}
	err := db.DB.Raw("Select * from rooms where room_id =?", roomId).Scan(&rooms).Error
	if err!=nil{
		response.ShowResponse(
			"Failure",
			500,
			"Record not found in DB",
			err.Error(),
			w,
		)
		return
	}

	
	response.ShowResponse(
		"Success",
		200,
		"Room Details",
		rooms,
		w,
	)
}

//fetch all the message from DB with limit
func MessageDetails(w http.ResponseWriter , r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	commonFunctions.EnableCors(&w)
	var mp = make(map[string]interface{})
	json.NewDecoder(r.Body).Decode(&mp)
	roomId,ok := mp["roomId"]
	if roomId == ""{
		response.ShowResponse(
			"Failure",
			404,
			"Room Id not found",
			"",
			w,
		)
		return
	}
	if !ok{
		response.ShowResponse(
			"Failure",
			400,
			"Unable to fetch room_id or invalid room_id",
			"",
			w,
		)
		return
	}

	var exists bool
	db.DB.Raw("SELECT EXISTS(SELECT * FROM rooms where room_id =?)",roomId).Scan(&exists)
	if !exists{
		response.ShowResponse(
			"Failure",
			500,
			"Record not found in DB",
			"",
			w,
		)
		return
	}

	var messages []models.Message

	db.DB.Raw("SELECT * FROM messages where room_id=? ORDER BY created_at  DESC LIMIT 10",roomId).Scan(&messages)

	response.ShowResponse(
		"Success",
		200,
		"All the messages in this room have been shown below",
		messages,
		w,
	)
}

//Handler to fetch all the details of the user provided that in which room it is present
func UserRoomsDetails(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	commonFunctions.EnableCors(&w)

	var mp = make(map[string]interface{})
	json.NewDecoder(r.Body).Decode(&mp)
	userId , ok := mp["userId"]
	if userId == ""{
		response.ShowResponse(
			"Failure",
			404,
			"Room Id not found",
			"",
			w,
		)
		return
	}
	if !ok{
		response.ShowResponse(
			"Failure",
			400,
			"Unable to fetch user_id or invalid user_id",
			"",
			w,
		)
		return
	}	

	var userRoomInformation []models.Participant	

	var exists bool
	db.DB.Raw("SELECT EXISTS(select * from participants where user_id = ?)",userId).Scan(&exists)
	if !exists {
		response.ShowResponse(
			"Failure",
			400,
			"User is alredy logged out",
			"",
			w,


		)
		return
	}

	db.DB.Raw("select * from participants where user_id = ?",userId).Scan(&userRoomInformation)

	response.ShowResponse(
		"Success",
		200,
		"User has access to following rooms",
		userRoomInformation,
		w,
	)
}