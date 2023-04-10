package Utils

import (
	db "main/Database"
	"net/http"
	"github.com/golang-jwt/jwt/v4"
	models "main/Models"
	"fmt"
	"os"
	
)

func BlacklistTokenHandler(token string) bool {
	var blacklisted bool
	db.DB.Raw("SELECT EXISTS(SELECT * FROM blacklisted_tokens where token=?)", token).Scan(&blacklisted)
	return blacklisted
}

func EnableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func SetHeader(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func GenerateJwtToken(user models.User , phone string , w http.ResponseWriter) string {
	//create user claims
	claims := models.Claims{
		User_id: user.User_Id,
		Phone:   phone,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(TokenExpirationDuration),
		},
	}
	fmt.Println("claims: ", claims)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Println("token: ", token)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWTKEY")))
	if err != nil {
		fmt.Println("error is :", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	fmt.Println("new token string :", tokenString)

	return tokenString
}