package chat

import (
	"fmt"
	db "main/Database"
	models "main/Models"

	socketio "github.com/googollee/go-socket.io"
)

func RoomLeave(s socketio.Conn, data map[string]interface{}) {

	room_id, ok1 := data["room_id"].(string)
	user_id, ok2 := data["user_id"].(string)
	fmt.Println("room_id: ", room_id)
	fmt.Println("user_id: ", user_id)
	if !ok1 || !ok2 {
		fmt.Println("invalid data provided while joining room")
		return

	}
	// participants table action
	var exists bool
	err := db.DB.Raw("SELECT EXISTS(SELECT 1 FROM participants WHERE user_id = ? AND room_id = ?)", user_id, room_id).Scan(&exists).Error
	if err != nil {
		panic(err)
	}
	// Check if the participant exists
	if !exists {
		fmt.Println("Participant does not exists")
		return
	} else {
		//soft delete functionality
		// db.DB.Where("user_id=? AND room_id=?", user_id, room_id).Delete(&models.Participant{})
		db.DB.Where("user_id=?", user_id).Updates(&models.Participant{Is_deleted: true})
		// reduce the user-count from rooms table
		var room models.Room
		db.DB.Raw("SELECT * FROM rooms where room_id=?", room_id).Scan(&room)
		if room.User_count >= 1 {
			db.DB.Model(&room).Where("room_id=?", room_id).Updates(&models.Room{User_count: room.User_count - 1})
		} else {
			fmt.Println("This room does not contain any clients")
		}
	}

}
