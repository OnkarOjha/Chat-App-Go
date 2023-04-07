package Controllers

import (
	"encoding/json"
	"fmt"
	db "main/Database"
	models "main/Models"
	response "main/Response"
	validator "main/Validation"
	"net/http"
)

// Message Searching API
func MessageSearchController(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	EnableCors(&w)

	var mp = make(map[string]interface{})

	json.NewDecoder(r.Body).Decode(&mp)
	roomId := mp["roomId"]
	search := mp["search"]
	if roomId == nil {
		response.ShowResponse(
			"Failure",
			400,
			"Empty room Id",
			"",
			w,
		)
		return
	}
	var roomExists bool
	db.DB.Raw("SELECT EXISTS(SELECT * FROM rooms where room_id=?)", roomId).Scan(&roomExists)
	if !roomExists {
		response.ShowResponse(
			"Failure",
			400,
			"Room does not exist",
			"",
			w,
		)
		return
	}

	//validator
	err := validator.CheckValidation(search)
	if err != nil {
		response.ShowResponse(
			"Failure",
			400,
			"",
			err.Error(),
			w,
		)
		return
	}

	var messages []models.Message
	var messageTextExist bool
	query1 := "SELECT EXISTS(select * from messages where text LIKE '%" + search.(string) + "%')"

	db.DB.Raw(query1).Scan(&messageTextExist)
	fmt.Println("message test exists:", messageTextExist)
	if messageTextExist{
		query2 := "select * from messages where room_id='"+ roomId.(string) +"' and text LIKE '%" + search.(string) + "%'"
		db.DB.Raw(query2).Scan(&messages)
		response.ShowResponse(
			"Success",
			200,
			"Here are the List of Messages",
			messages,
			w,
		)
		return
	}else{
		response.ShowResponse(
			"Failure",
			400,
			"Message not found",
			"",
			w,
		)
		return
	}

}