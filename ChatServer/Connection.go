package ChatServer

import (
	"errors"
	"fmt"
	db "main/Database"
	models "main/Models"
	"net/url"
	"time"

	socketio "github.com/googollee/go-socket.io"
	"gorm.io/gorm"
)

// var roomName string = "room1"

func Connect(s socketio.Conn) error {

	// Get the user ID from the query params
	rawQuery := s.URL().RawQuery
	query, _ := url.ParseQuery(rawQuery)
	fmt.Printf("Client %s Connected \n", query["id"][0])
	return nil
}

func RoomJoin(s socketio.Conn, data map[string]interface{}) {
	// Get the user ID from the query params
	rawQuery := s.URL().RawQuery
	query, _ := url.ParseQuery(rawQuery)
	RoomName, ok1 := data["room_name"].(string)
	topicName, ok2 := data["topic_name"].(string)
	fmt.Println("Room name: ", RoomName)
	fmt.Println("Topic name: ", topicName)
	if !ok1 || !ok2 {
		fmt.Println("invalid data provided while joining room")
		return
	}
	var client models.User
	var room models.Room

	// check that the room already exists or not
	err := db.DB.Where("name=?", RoomName).First(&room).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.Join(RoomName)

		dateStr := time.Now().Truncate(time.Hour)
		room.Created_at = dateStr.Format("02 Jan 2006")
		db.DB.Where("user_id", query["id"][0]).First(&client)
		fmt.Println("Connection id:", s.ID())
		// TODO admin check
		room.Admin_id = query["id"][0]

		room.User_count++
		room.Name = RoomName
		room.Topic_name = topicName

		db.DB.Create(&room)

		// // check that if the user who is hitting the conn already exists in the participants table then don't update the count
		// if !CheckParticipants(client.User_Id, room.Room_id) {
		// 	var roomcount int
		// 	roomcount++
		// 	fmt.Println("room count: ", roomcount)
		// 	db.DB.Model(&room).Where("room_id=?", room.Room_id).Updates(&models.Room{User_count: roomcount})
		// }

		Paricipants(query["id"][0], room.Room_id, RoomName)

	} else {
		s.Join(RoomName)
		var roomParticipants models.Room
		db.DB.Where("name = ?", RoomName).First(&roomParticipants)
		roomParticipants.User_count++

		// participant table updation
		Paricipants(query["id"][0], roomParticipants.Room_id, RoomName)

		// // check that if the user who is hitting the conn already exists in the participants table then don't update the count
		// if !CheckParticipants(client.User_Id, room.Room_id) {
		// 	var roomcount int
		// 	roomcount++
		// 	db.DB.Model(&room).Where("room_id=?", room.Room_id).Updates(&models.Room{User_count: roomcount})
		// }

		db.DB.Where("name=?", RoomName).Updates(&room)
		fmt.Printf("room %s joined\n", RoomName)
	}

}

// Participant table updation as soon as new user joins the room
func Paricipants(user_id string, room_id string, RoomName string) {

	var participants models.Participant

	if !CheckParticipants(user_id, room_id) {
		return
	}

	participants.User_id = user_id
	participants.Room_id = room_id
	participants.Room_name = RoomName
	db.DB.Create(&participants)

}

// check that if user already exists don't create new participant entry
func CheckParticipants(user_id string, room_id string) bool {

	var exists bool
	db.DB.Raw("SELECT * FROM participants where user_id = ?  AND room_id = ? ", user_id, room_id).Row().Scan(&exists)
	fmt.Println("exists: ", exists)
	if exists {
		fmt.Println("not creating new participant")
		return true
	}

	return false
}
