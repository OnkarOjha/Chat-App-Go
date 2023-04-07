package chat

import (
	"fmt"
	db "main/Database"
	models "main/Models"
	response "main/Response"
	socketio "github.com/googollee/go-socket.io"
)

func RoomLeave(s socketio.Conn, data map[string]interface{}) {

	roomId, ok1 := data["roomId"].(string)
	userId, ok2 := data["userId"].(string)
	fmt.Println("room_id: ", roomId)
	fmt.Println("user_id: ", userId)
	if roomId == "" || userId == "" {
		response.SocketResponse(
			"Failure",
			"Either roomId or userId is empty",
			s,
		)
		return
	}

	if !ok1 || !ok2 {
		response.SocketResponse(
			"Failure",
			"Invalid data provided while joining room",
			s,
		)
		return

	}
	// participants table action
	var exists bool
	err := db.DB.Raw("SELECT EXISTS(SELECT 1 FROM participants WHERE user_id = ? AND room_id = ?)", userId, roomId).Scan(&exists).Error
	if err != nil {
		response.SocketResponse(
			"Failure",
			"Error while searching participant",
			s,
		)
		return
	}
	// Check if the participant does not exists and if he is already deleted
	var deleted bool
	err = db.DB.Raw("SELECT EXISTS(select * from participants where user_id = ? and is_deleted = ?)", userId , true).Scan(&deleted).Error
	if err != nil {
		response.SocketResponse(
			"Failure",
			"Error while searching participant",
			s,
		)
		return
	}
	if !exists || deleted {
		response.SocketResponse(
			"Failure",
			"This User_id does not exists in this room",
			s,
		)
		return
	} else {
		//soft delete functionality
		// db.DB.Where("user_id=? AND room_id=?", user_id, room_id).Delete(&models.Participant{})
		db.DB.Where("user_id=? and room_id=?", userId ,roomId).Updates(&models.Participant{Has_left: true})
		// reduce the user-count from rooms table
		var room models.Room
		db.DB.Raw("SELECT * FROM rooms where room_id=?", roomId).Scan(&room)
		if room.User_count >0 {
			fmt.Println("decreasing room count")
			fmt.Println("room.user_count = ", room.User_count)
			db.DB.Model(&models.Room{}).Where("room_id=?", roomId).Update("user_count" , room.User_count - 1)
			fmt.Println("user count: " , room.User_count)
		} else {
			response.SocketResponse(
				"Failure",
				"This room doesn't contain any more clients",
				s,
			)
			return
		}
	}

}
