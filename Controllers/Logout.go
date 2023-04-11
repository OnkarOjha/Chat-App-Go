package Controllers

import (
	"fmt"
	db "main/Database"
	models "main/Models"
	response "main/Response"
	commonFunctions "main/Utils"
	"net/http"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("logout handler")
	dataFromContext := r.Context().Value("editUser")
	userDetails := dataFromContext.(map[string]interface{})
	fmt.Printf("user_id from context: %v\n", userDetails["userId"])
	w.Header().Set("Content-Type", "application/json")
	commonFunctions.EnableCors(&w)

	token := r.Header["Token"]
	fmt.Println("token : ", token[0])
	fmt.Println("You are Logging out user :",userDetails["userId"])


	var user models.User

	var exists bool
	db.DB.Raw("SELECT EXISTS(select * from users where user_id=?)", userDetails["userId"]).Scan(&exists)
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

	err := db.DB.Raw("SELECT * from users where user_id = ?", userDetails["userId"]).Scan(&user).Error
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
		// store the token in blacklisted table
		var blackListedToken models.BlacklistedTokens
		blackListedToken.Token = token[0]
		db.DB.Create(&blackListedToken)
		user.Is_active = false
		db.DB.Model(&models.User{}).Where("user_id=?", userDetails["userId"]).Update("is_active", false)
		db.DB.Model(&models.User{}).Where("user_id=?", userDetails["userId"]).Update("token" , nil)
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
