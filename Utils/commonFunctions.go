package Utils

import (
	"fmt"
	db "main/Database"
	models "main/Models"
	response "main/Response"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

//Blacklist a token on cookie expiration , logout and delete handlers
func BlacklistTokenHandler(token string) bool {
	var blacklisted bool
	db.DB.Raw("SELECT EXISTS(SELECT * FROM blacklisted_tokens where token=?)", token).Scan(&blacklisted)
	return blacklisted
}

// Function to enable CORS
func EnableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

// Function to Allow all headers
func SetHeader(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

//Generate Token Function
func GenerateJwtToken(user models.User, phone string, w http.ResponseWriter) string {
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

//Decode Token Function
func DecodeToken(tokenString string) (models.Claims, error) {
	claims := &models.Claims{}

	parsedToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error")
		}
		return []byte(os.Getenv("JWTKEY")), nil
	})

	if err != nil || !parsedToken.Valid {
		return *claims, fmt.Errorf("invalid or expired token")
	}

	return *claims, nil
}

//Set cookie handler
func SetCookie(w http.ResponseWriter, r *http.Request, tokenString string) {
	cookie := http.Cookie{
		Name:    "cookie",
		Value:   tokenString,
		Expires: CookieExpirationTime,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	response.ShowResponse(
		"Success",
		200,
		"Cookies saved successfully",
		&cookie,
		w,
	)
}

//Delete cookie handler
func DeleteCookie(w http.ResponseWriter, r *http.Request, cookie *http.Cookie) {
	c := http.Cookie{
		Name:    "cookie",
		Expires: time.Now(),
	}
	http.SetCookie(w, &c)
	response.ShowResponse(
		"Success",
		200,
		"Cookies Deleted successfully",
		&cookie,
		w,
	)
}
