package chat

import (
	"errors"
	"fmt"
	db "main/Database"
	models "main/Models"
	// response "main/Response"
	"net/url"
	"time"
	response "main/Response"
	socketio "github.com/googollee/go-socket.io"
	"gorm.io/gorm"
)

func RoomCreate(s socketio.Conn, data map[string]interface{}) {
	fmt.Println("inside room creation process...")
	// Get the user ID from the query params
	rawQuery := s.URL().RawQuery
	query, _ := url.ParseQuery(rawQuery)
	RoomName, ok1 := data["roomName"].(string)
	topicName, ok2 := data["topicName"].(string)
	fmt.Println("Room name: ", RoomName)
	fmt.Println("Topic name: ", topicName)
	if RoomName == ""{
		response.SocketResponse(
			"Failure",
			"Room Name must be provided",
			s,
		)
		return
	}
	if !ok1 && !ok2 {
		response.SocketResponse(
			"Failure",
			"Either Room Name or Topic Name is missing",
			s,
		)
		return
	}
	var client models.User
	var room models.Room

	// check that the room already exists or not
	err := db.DB.Where("name=?", RoomName).First(&room).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("roomnbame: ", RoomName)
		s.Join(RoomName)

		dateStr := time.Now().Truncate(time.Hour)
		room.Created_at = dateStr.Format("02 Jan 2006")
		db.DB.Where("user_id", query["id"][0]).First(&client)
		fmt.Println("Connection id:", s.ID())
		
		room.Admin_id = query["id"][0]

		room.Name = RoomName
		room.Topic_name = topicName

		db.DB.Create(&room)

		// check that if the user who is hitting the conn already exists in the participants table then don't update the count
		if !CheckParticipants(client.User_Id, room.Room_id ,s) {
			var roomcount int64
			roomcount++
			fmt.Println("room count: ", roomcount)
			db.DB.Model(&room).Where("room_id=?", room.Room_id).Updates(&models.Room{User_count: roomcount})
		}

		Paricipant(query["id"][0], room.Room_id, RoomName , s)
		CheckRoomClient(RoomName , room.Room_id)


	} else {
		s.Join(RoomName)

		var roomParticipants models.Room
		db.DB.Where("name = ?", RoomName).First(&roomParticipants)

		// participant table updation
		Paricipant(query["id"][0], roomParticipants.Room_id, RoomName , s)

		// check that if the user who is hitting the conn already exists in the participants table then don't update the count
		if !CheckParticipants(client.User_Id, room.Room_id ,s) {
			var roomcount int64
			roomcount++
			db.DB.Model(&room).Where("room_id=?", room.Room_id).Updates(&models.Room{User_count: roomcount})
		}

		db.DB.Where("name=?", RoomName).Updates(&roomParticipants)
		fmt.Printf("room %s joined\n", RoomName)
		CheckRoomClient(RoomName , roomParticipants.Room_id)
	}
	fmt.Println("Rooms are: ", s.Rooms())

}

// Participant table updation as soon as new user joins the room
func Paricipant(user_id string, room_id string, RoomName string ,  s socketio.Conn) {

	var participants models.Participant

	if CheckParticipants(user_id, room_id ,s) {
		return
	}

	// if exists is false then create new participant
	participants.User_id = user_id
	participants.Room_id = room_id
	participants.Room_name = RoomName
	db.DB.Create(&participants)

}

// check that if user already exists don't create new participant entry
func CheckParticipant(user_id string, room_id string) bool {

	var exists bool
	err := db.DB.Raw("SELECT EXISTS(SELECT 1 FROM participants WHERE user_id = ? AND room_id = ?)", user_id, room_id).Scan(&exists).Error
	if err != nil {
		panic(err)
	}

	// Check if the participant exists
	if exists {
		fmt.Println("not creating participant")
	} 

	return exists
}


//  function to append room user count according to participant table
func CheckRoomClient(roomName string , room_id string){
		var count int64
		db.DB.Raw("SELECT COUNT(room_name) from participants where room_name = ? group by room_name;",roomName).Scan(&count)
		fmt.Println("count:",count)
		db.DB.Where("room_id=?",room_id).Updates(&models.Room{User_count: count})
		
}