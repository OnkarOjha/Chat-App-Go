package Controllers

import (
	"encoding/json"
	"fmt"
	db "main/Database"
	models "main/Models"
	response "main/Response"
	"net/http"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("logout handler")
	w.Header().Set("Content-Type", "application/json")
	EnableCors(&w)
	var mp = make(map[string]interface{})
	json.NewDecoder(r.Body).Decode(&mp)
	userId, ok := mp["userId"]
	fmt.Println("You are Logging out user :", userId)
	if mp["userId"] == nil {
		response.ShowResponse(
			"Failure",
			400,
			"Empty userId",
			"",
			w,
		)
		return
	}
	if !ok {
		response.ShowResponse(
			"Failure",
			400,
			"Error fetching userId",
			"",
			w,
		)
		return
	}

	// Token expiration
	var user models.User

	var exists bool
	db.DB.Raw("SELECT EXISTS(select * from users where user_id=?)", userId).Scan(&exists)
	if !exists {
		response.ShowResponse(
			"Failure",
			400,
			"User does not exist",
			"",
			w,
		)
		return
	}

	err := db.DB.Raw("SELECT * from users where user_id = ?", userId).Scan(&user).Error
	if err != nil {
		response.ShowResponse(
			"Failure",
			500,
			"User not found",
			err.Error(),
			w,
		)
		return
	}
	if user.Is_active == true {
		
		user.Is_active = false
		db.DB.Model(&models.User{}).Where("user_id=?", userId).Update("is_active", false)
		db.DB.Model(&models.User{}).Where("user_id=?", userId).Update("token" , nil)
		fmt.Println("user now:", user)
		response.ShowResponse(
			"Success",
			200,
			"User successfully logged out",
			user,
			w,
		)
	}
}
