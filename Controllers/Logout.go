package Controllers

import (
	"encoding/json"
	"fmt"
	db "main/Database"
	models "main/Models"
	response "main/Response"
	"net/http"
	"os"
	"time"

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

	// Token expiration
	var user models.User

	db.DB.Raw("SELECT * from users where user_id = ?", user_id).Scan(&user)
	tokenstring := user.Token
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenstring, claims, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error")
		}
		return []byte(os.Getenv("JWTKEY")), nil
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("token:", token)
	fmt.Println("expiration time before: ", claims.RegisteredClaims.ExpiresAt)
	claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now())
	fmt.Println("expiration time now: ", claims.RegisteredClaims.ExpiresAt)
	fmt.Println("user_id:", user_id)
	user.Is_active = false
	db.DB.Model(&models.User{}).Where("user_id=?", user_id).Update("is_active", false)
	fmt.Println("user now:", user)
	response.ShowResponse(
		"OK",
		200,
		"User successfully logged out",
		user,
		w,
	)
}
