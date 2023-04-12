package Controllers

import (
	"encoding/json"
	"fmt"
	db "main/Database"
	models "main/Models"
	response "main/Response"
	commonFunctions "main/Utils"
	validator "main/Validation"
	"net/http"
)

//	@Summary		Message Searching for Rooms
//	@Description	All the message matching the message text
//	@Tags			Chat-Room API
//	@Accept			json
//	@Produce		json
//  @Param          User body string true "userId of the user" SchemaExample({"roomId":"string" , "search" : "string"})
//  @Success        200 {object}    response.Response
//	@Failure		404	{string}	response.Response
//	@Failure		400	{string}	response.Response
//	@Failure		500	{string}	response.Response
//	@Router			/messageSearch [post]
func MessageSearchController(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	commonFunctions.EnableCors(&w)

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

	query1 := "SELECT EXISTS(select * from messages where LOWER(text) LIKE LOWER('%\\" + search.(string) + "%'))"

	db.DB.Raw(query1).Scan(&messageTextExist)
	fmt.Println("message test exists:", messageTextExist)
	if messageTextExist {
		query2 := "select * from messages where room_id='" + roomId.(string) + "' and LOWER(text) LIKE LOWER('%\\" + search.(string) + "%')"
		db.DB.Raw(query2).Scan(&messages)
		response.ShowResponse(
			"Success",
			200,
			"Here are the List of Messages",
			messages,
			w,
		)
		return
	} else {
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
