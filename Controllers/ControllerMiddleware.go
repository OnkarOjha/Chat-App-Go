package Controllers

import (
	"fmt"
	db "main/Database"
	models "main/Models"
	response "main/Response"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

func IsAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("inside middleware")
		id := r.URL.Query().Get("id")
		userId := id
		if r.Header["Token"] != nil {

			token, err := jwt.ParseWithClaims(r.Header["Token"][0], &models.Claims{}, func(token *jwt.Token) (interface{}, error) {

				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("error")
				}
				return []byte(os.Getenv("JWTKEY")), nil
			})
			if err != nil {
				panic(err)
			}

			if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
				var user models.User
				db.DB.Raw("Select * from users where user_id=?", userId).Scan(&user)

				if claims.User_id == user.User_Id {

					endpoint(w, r)
				} else {
					response.ShowResponse(
						"Failure",
						401,
						"Check header token or user-id provided",
						"",
						w,
					)
					return
				}
			} else {
				// fmt.Println("Bad Request")
				response.ShowResponse(
					"Failure",
					400,
					"Check user claims",
					"",
					w,
				)
				return

			}

		}

	})
}
