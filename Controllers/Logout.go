package Controllers

import (
	"encoding/json"
	"fmt"
	db "main/Database"
	models "main/Models"
	response "main/Response"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	EnableCors(&w)
	var mp = make(map[string]interface{})
	json.NewDecoder(r.Body).Decode(&mp)
	user_id, ok := mp["user_id"]
	fmt.Println("You are Logging out user :", user_id)
	if !ok {
		response.ShowResponse(
			"Bad Request",
			400,
			"Error fetching user_id",
			"",
			w,
		)
		return
	}

	// // db mei jake us user ko dundo
	// db.DB.Where("user_id = ?", user_id).Updates(&models.User{Is_active: false})

	// End his token expiration time to time.Now()
	var user models.User

	db.DB.Raw("SELECT * from users where user_id = ?", user_id).Scan(&user)
	fmt.Println("user: ", user)
	tokenstring := user.Token
	fmt.Println("token: ", tokenstring)
	//TODO
	token, err := jwt.ParseWithClaims(tokenstring, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error")
		}
		return []byte(os.Getenv("JWTKEY")), nil
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("token:", token)
}
