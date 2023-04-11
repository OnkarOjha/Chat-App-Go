package chat

import (
	"fmt"
	db "main/Database"
	models "main/Models"
	response "main/Response"
	validation "main/Validation"
	"net/url"

	socketio "github.com/googollee/go-socket.io"
)

// Socket Instance Connection
func Connect(s socketio.Conn) error {

	rawQuery := s.URL().RawQuery
	query, _ := url.ParseQuery(rawQuery)
	fmt.Printf("Client %s Connected \n", query["id"][0])
	return nil
}

// Room Joining Handler Function
func RoomJoin(s socketio.Conn, data map[string]interface{}) {
	fmt.Println("inside room join")
	rawQuery := s.URL().RawQuery
	query, _ := url.ParseQuery(rawQuery)
	roomId, _ := data["roomId"].(string)
	if roomId == "" {
		response.SocketResponse(
			"Failure",
			"Room id not found",
			s,
		)
		return
	}

	err := validation.CheckValidation(data)
	if err != nil {
		response.SocketResponse(
			"Failure",
			"Room id not found",
			s,
		)
		return
	}

	fmt.Println("Room id: ", roomId)

	//check that room is deleted
	var roomDeleted bool
	db.DB.Raw("SELECT EXISTS(select * from rooms where room_id=? and is_deleted=true)", roomId).Scan(&roomDeleted)
	if roomDeleted {
		response.SocketResponse(
			"Failure",
			"Room is already deleted",
			s,
		)
		return
	}

	var client models.User
	var room models.Room

	s.Join(roomId)

	response.SocketResponse(
		roomId,
		"User Successfully joined Room "+roomId,
		s,
	)
	var roomParticipants models.Room
	err = db.DB.Where("room_id = ?", roomId).First(&roomParticipants).Error
	if err != nil {
		response.SocketResponse(
			"",
			err.Error(),
			s,
		)
		return
	}

	// participant table updation
	Paricipants(query["id"][0], roomParticipants.Room_id, s)

	// check that if the user who is hitting the conn already exists in the participants table then don't update the count
	if !CheckParticipants(client.User_Id, room.Room_id, s) {
		var roomCounting models.Room
		db.DB.Raw("SELECT * FROM rooms WHERE room_id =?", roomId).Scan(&roomCounting)
		roomCounting.User_count += 1
		roomcount := roomCounting.User_count
		fmt.Println("roomcount: ", roomcount)
		db.DB.Model(&models.Room{}).Where("room_id=?", roomId).Update("user_count", roomcount)
	}

	db.DB.Where("room_id=?", roomId).Updates(&roomParticipants)
	fmt.Printf("room %s joined\n", roomParticipants.Name)
	CheckRoomClients(roomId)

	fmt.Println("Rooms are: ", s.Rooms())

}

// Participant table updation as soon as new user joins the room
func Paricipants(user_id string, roomId string, s socketio.Conn) {

	var participants models.Participant

	if CheckParticipants(user_id, roomId, s) {
		return
	}

	// if exists is false then create new participant
	participants.User_id = user_id
	participants.Room_id = roomId
	// participants.Room_name = RoomName
	db.DB.Create(&participants)

}

// check that if user already exists don't create new participant entry
func CheckParticipants(user_id string, room_id string, s socketio.Conn) bool {

	var exists bool
	err := db.DB.Raw("SELECT EXISTS(SELECT 1 FROM participants WHERE user_id = ? AND room_id = ?)", user_id, room_id).Scan(&exists).Error
	if err != nil {
		response.SocketResponse(
			"Participant does not exist",
			"Failure",
			s,
		)
		return false
	}

	// Check if the participant exists
	if exists {
		fmt.Println("not creating participant")
	}

	return exists
}

//  function to append room user count according to participant table
func CheckRoomClients(roomId string) {
	var count int64
	db.DB.Raw("SELECT COUNT(room_name) from participants where room_id = ? group by room_name;", roomId).Scan(&count)
	count += 1
	fmt.Println("count:", count)
	db.DB.Where("room_id=?", roomId).Updates(&models.Room{User_count: count})

}
