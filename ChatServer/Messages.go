package chat

import (
	"fmt"
	db "main/Database"
	models "main/Models"
	response "main/Response"
	server "main/Utils"
	"net/url"

	socketio "github.com/googollee/go-socket.io"
)

//Message sending Handler
func Messages(s socketio.Conn, data map[string]interface{}) {
	var message models.Message
	roomId, ok1 := data["roomId"].(string)
	text, ok2 := data["text"].(string)
	fmt.Println("Message in room: ", roomId)
	fmt.Println("Message Text is: ", text)
	if roomId == "" {
		response.SocketResponse(
			"Failure",
			"Room id not found",
			s,
		)
		return
	}
	if !ok1 || !ok2 {
		response.SocketResponse(
			"Failure",
			"Either RoomID or Message is missing",
			s,
		)
		return
	}

	response.SocketResponse(
		text,
		"Message sent in room"+roomId,
		s,
	)

	// pick the user_id from params
	rawQuery := s.URL().RawQuery
	query, _ := url.ParseQuery(rawQuery)
	message.User_id = query["id"][0]
	fmt.Println("Message by: ", message.User_id)

	// search for the room
	var room models.Room
	err := db.DB.Where("room_id=?", roomId).First(&room).Error
	if err != nil {
		response.SocketResponse(
			"Failure",
			err.Error(),
			s,
		)
		return
	}
	message.Room_id = room.Room_id
	fmt.Println("message in room with id: ", message.Room_id)

	// set the message text
	message.Message_type = "Text Message"
	message.Text = text
	fmt.Println("message content is: ", message.Text)

	// Broadcast "typing" event to other users in the room
	typingResponse := models.BroadcastResponse{
		MessageBy: message.User_id,
		MessageIn: message.Room_id,
	}
	server.Server.BroadcastToRoom("/", roomId, "typing", typingResponse)
	broadcast := server.Server.BroadcastToRoom("/", roomId, "reply", text)
	if broadcast {
		fmt.Println("message  broadcasted: ", text)
	}

	// // // Emit "stop typing" event when user stops typing
	// // s.Emit("stoptyping", "Message by"+ message.User_id , "Message in" + message.Room_id)
	// server.Server.BroadcastToRoom("/", roomId, "stoptyping", message.User_id)

	db.DB.Create(&message)

}
