package chat

import (
	"fmt"
	db "main/Database"
	model "main/Models"
	cons "main/Utils"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt/v4"
	socketio "github.com/googollee/go-socket.io"
)

func TokenCheck(next func(socketio.Conn) error) func(socketio.Conn) error {
	return func(conn socketio.Conn) error {
		// get token from header
		httpheader := conn.RemoteHeader()
		fmt.Println("httpheader: ", httpheader["Token"][0])

		// get user_id from params
		rawQuery := conn.URL().RawQuery
		query, _ := url.ParseQuery(rawQuery)
		
		// middleware logic
		parsedToken, err := jwt.ParseWithClaims(httpheader["Token"][0], &model.Claims{}, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("error")
			}
			return cons.JwtKey, nil
		})
		var user model.User
		if claims, ok := parsedToken.Claims.(*model.Claims); ok && parsedToken.Valid {
			// refresh token logic
			if time.Until(claims.ExpiresAt.Time) < cons.TokenRefreshInterval {
				// Generate a new token with a new expiration time
				newExpirationTime := cons.TokenExpirationDuration
				//create user claims
				claims := model.Claims{
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(newExpirationTime),
					},
				}	
				newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				newTokenString, err := newToken.SignedString(cons.JwtKey)
				if err != nil {
					return err
				}
				fmt.Println("new token string: ", newTokenString)
				db.DB.Model(&user).Where("user_id=?", user.User_Id).Updates(&model.User{Token: newTokenString})

			}

			// phone number utha k lao from database

			db.DB.Raw("Select * from users where user_id=?", query["id"][0]).Scan(&user)
			if claims.Phone == user.Phone {

				return next(conn)
			}
		} else {
			fmt.Println(err)
			// return err
		}

		// Now we call the actual function
		return next(conn)
		
	}
}
