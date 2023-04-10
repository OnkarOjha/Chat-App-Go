package Controllers

import (
	"encoding/json"
	"fmt"
	models "main/Models"
	"net/http"
	db "main/Database"
	response "main/Response"
	commonFunctions "main/Utils"
)

func RoomDelete(w http.ResponseWriter, r *http.Request) {
	// only that person who is the admin can delete the room
	// var room models.Room
	w.Header().Set("Content-Type", "application/json")
	commonFunctions.EnableCors(&w)
	fmt.Println("We are deleting the room")
	var mp = make(map[string]string)
	json.NewDecoder(r.Body).Decode(&mp)
	roomId := mp["roomId"]
	adminId := mp["adminId"]
	if roomId == ""  {
		response.ShowResponse(
			"Failure",
			400,
			"RoomId missing",
			"",
			w,
		)
		return
	}else if adminId == ""{
		response.ShowResponse(
			"Failure",
			400,
			"AdminId missing",
			"",
			w,
		)
		return
	}
	fmt.Println("room_id: ", roomId)
	fmt.Println("admin_id: ", adminId)

	// check if room exists or not
	 var roomexists bool
	 db.DB.Raw("SELECT EXISTS (SELECT * from rooms where room_id =? and is_deleted=false)",roomId).Scan(&roomexists)
	 if !roomexists{
		response.ShowResponse(
			"Failure",
			400,
			"No such Room exists",
			"",
			w,
		)
		return
	 }

	//admin check 
	var isAdmin bool
	err := db.DB.Raw("SELECT EXISTS(SELECT 1 FROM rooms WHERE admin_id=?)", adminId).Scan(&isAdmin).Error
	if err != nil {
		panic(err)
	}
	fmt.Println("isadmin:", isAdmin)
	// Check if the participant exists
	if !isAdmin {
		response.ShowResponse(
			"Failure",
			400,
			"User is not an admin",
			"",
			w,
		)
		return
	}else{
		db.DB.Where("room_id=?", roomId).Updates(&models.Room{Is_deleted: true})
	}
}
