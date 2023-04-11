package chat

import (
	"fmt"
	db "main/Database"
	model "main/Models"
	"net/url"
	"os"
	commonFunctions "main/Utils"
	"github.com/golang-jwt/jwt/v4"
	socketio "github.com/googollee/go-socket.io"
)

// Middleware for token verification
func TokenCheck(next func(socketio.Conn) error) func(socketio.Conn) error {
	return func(conn socketio.Conn) error {
		//TODO get token from header
		
		httpheader := conn.RemoteHeader()
		fmt.Println("httpheader: ", httpheader["Token"][0])

		// get user_id from params
		rawQuery := conn.URL().RawQuery
		query, _ := url.ParseQuery(rawQuery)

		if commonFunctions.BlacklistTokenHandler(httpheader["Token"][0]){
			fmt.Println("token is blacklisted")
			conn.Close()
		}

		// middleware logic
		parsedToken, err := jwt.ParseWithClaims(httpheader["Token"][0], &model.Claims{}, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("error")
			}
			return []byte(os.Getenv("JWTKEY")), nil
		})
		if err != nil {
			panic(err)
		}
		fmt.Println("dkjsnd: ",parsedToken)
		var user model.User
		if claims, ok := parsedToken.Claims.(*model.Claims); ok && parsedToken.Valid {
			fmt.Println("claims user_id: ", claims)
			// User details from database

			db.DB.Raw("Select * from users where user_id=?", query["id"][0]).Scan(&user)
			fmt.Println("user_id: ", user.User_Id)

			if claims.User_id != user.User_Id {
				conn.Close()
			}

			return next(conn)

		}
		return nil
	}

}
