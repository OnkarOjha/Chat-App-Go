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

//	@Summary		Participant Details 
//	@Description	Information about room which user is present in which room and the user had left the room or not
//	@Tags			Chat-Room API
//	@Accept			json
//	@Produce		json
//  @Param          User body string true "roomId of the user" SchemaExample({"roomId":"string"})
//  @Success        200 {object}    response.Response
//	@Failure		404	{string}	response.Response
//	@Failure		400	{string}	response.Response
//	@Failure		500	{string}	response.Response
//	@Router			/participants [post]
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

//	@Summary		Room Details 
//	@Description	Information about room (User Count , Deleted or not , Admin Id etc)
//	@Tags			Chat-Room API
//	@Accept			json
//	@Produce		json
//  @Param          User body string true "roomId of the user" SchemaExample({"roomId":"string"})
//  @Success        200 {object}    response.Response
//	@Failure		404	{string}	response.Response
//	@Failure		400	{string}	response.Response
//	@Failure		500	{string}	response.Response
//	@Router			/rooms [post]
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

//	@Summary		Message Details 
//	@Description	Information about All the messages sent in the room
//	@Tags			Chat-Room API
//	@Accept			json
//	@Produce		json
//  @Param          User body string true "roomId of the user" SchemaExample({"roomId":"string"})
//  @Success        200 {object}    response.Response
//	@Failure		404	{string}	response.Response
//	@Failure		400	{string}	response.Response
//	@Failure		500	{string}	response.Response
//	@Router			/messages [post]
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

//	@Summary		User Rooms 
//	@Description	Information about user room access
//	@Tags			user
//	@Accept			json
//	@Produce		json
//  @Param          User body string true "userId of the user" SchemaExample({"userId":"string"})
//  @Success        200 {object}    response.Response
//	@Failure		404	{string}	response.Response
//	@Failure		400	{string}	response.Response
//	@Failure		500	{string}	response.Response
//	@Router			/userRoomInfo [post]
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