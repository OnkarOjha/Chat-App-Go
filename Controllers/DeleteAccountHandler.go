package Controllers

import (
	"fmt"
	db "main/Database"
	models "main/Models"
	response "main/Response"
	commonFunctions "main/Utils"
	"net/http"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

//	@Summary		User Delete Handler
//	@Description	Deleting user Account Details
//	@Tags			user
//	@Accept			json
//	@Produce		json
//  @Param          User body string true "userId of the user" SchemaExample({"userId":"string"})
//	@Success		200	{string}	response.Response
//	@Failure		400	{string}	response.Response
//	@Failure		409	{string}	response.Response
//	@Failure		500	{string}	response.Response
//	@Router			/deleteAccount [post]
func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete account Handler")
	dataFromContext := r.Context().Value("editUser")
	userDetails := dataFromContext.(map[string]interface{})


	err := validation.Validate(userDetails,
		validation.Map(
			validation.Key("userId", validation.Required),
		),
	)
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

	w.Header().Set("Content-Type", "application/json")
	commonFunctions.EnableCors(&w)

	cookie, err := r.Cookie("cookie")
	if err!=nil{
		response.ShowResponse(
			"Failure",
			403,
			"Error fetching cookie",
			err.Error(),
			w,
		)
		return
	}


	// Token expiration
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

	err = db.DB.Raw("SELECT * from users where user_id = ?", userDetails["userId"]).Scan(&user).Error
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
		blackListedToken.Token = cookie.Value
		db.DB.Create(&blackListedToken)
		//delete user from DB
		db.DB.Model(&models.User{}).Where("user_id=?", userDetails["userId"]).Delete(&user)

		//delete cookie when user logout
		commonFunctions.DeleteCookie(w,r,cookie)
		
		response.ShowResponse(
			"Success",
			200,
			"User Account successfully Deleted",
			user,
			w,
		)
		return
	} else {
		response.ShowResponse(
			"Failure",
			400,
			"User is not logged in",
			"",
			w,
		)
		return
	}
}
