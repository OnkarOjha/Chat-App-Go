package Controllers

import (
	"encoding/json"
	"fmt"
	models "main/Models"
	"net/http"
	db "main/Database"
)

func RoomDelete(w http.ResponseWriter, r *http.Request) {
	// only that person who is the admin can delete the room
	// var room models.Room
	w.Header().Set("Content-Type", "application/json")
	EnableCors(&w)
	fmt.Println("We are deleting the room")
	var mp = make(map[string]string)
	json.NewDecoder(r.Body).Decode(&mp)
	room_id := mp["room_id"]
	admin_id := mp["admin_id"]
	fmt.Println("room_id: ", room_id)
	fmt.Println("admin_id: ", admin_id)

	//admin check 
	var isAdmin bool
	err := db.DB.Raw("SELECT EXISTS(SELECT 1 FROM rooms WHERE admin_id=?)", admin_id).Scan(&isAdmin).Error
	if err != nil {
		panic(err)
	}
	fmt.Println("isadmin:", isAdmin)
	// Check if the participant exists
	if !isAdmin {
		fmt.Println("This client is not an admin")
		return
	}else{
		db.DB.Where("room_id=?", room_id).Updates(&models.Room{Is_deleted: true})
	}
}
