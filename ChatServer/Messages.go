package chat

import (
	"fmt"
	db "main/Database"
	models "main/Models"
	"net/url"
	socketio "github.com/googollee/go-socket.io"
	server "main/Utils"
)

func Messages(s socketio.Conn, data map[string]interface{}) {
	var message models.Message
	roomName, ok1 := data["room_name"].(string)
	text, ok2 := data["text"].(string)
	fmt.Println("Message in room: ", roomName)
	fmt.Println("Message Text is: ", text)
	if !ok1 || !ok2 {
		fmt.Println("invalid data provided while joining room")
		return
	}
	s.Emit("reply", "message is: "+ text)

	// pick the user_id from params
	rawQuery := s.URL().RawQuery
	query, _ := url.ParseQuery(rawQuery)
	message.User_id = query["id"][0]
	fmt.Println("Message by: ", message.User_id)

	// search for the room
	var room models.Room
	db.DB.Where("name=?", roomName).First(&room)
	message.Room_id = room.Room_id
	fmt.Println("message in room with id: ", message.Room_id)

	// set the message text
	message.Text = text
	fmt.Println("message content is: ", message.Text)
	

	broadcast := server.Server.BroadcastToRoom("/",roomName , "reply" , "message broadcasted is : "+ text)
	if broadcast{
		fmt.Println("message  broadcasted: ",text)
	}
	
	db.DB.Create(&message)

}
