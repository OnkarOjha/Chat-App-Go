package Controllers

import (
	"context"
	"encoding/json"
	"fmt"
	db "main/Database"
	models "main/Models"
	response "main/Response"
	commonFunctions "main/Utils"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

func IsAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("inside middleware")
		var mp = make(map[string]interface{})
		json.NewDecoder(r.Body).Decode(&mp)
		fmt.Println("data in middlware after decoding: ", mp) 
		userId, ok := mp["userId"]
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
		// check if token is blacklisted
		if commonFunctions.BlacklistTokenHandler(r.Header["Token"][0]) {
			response.ShowResponse(
				"Failure",
				400,
				"Token is Blacklisted",
				"",
				w,
			)
			return
		}
		claims := &models.Claims{}

		if r.Header["Token"] != nil {

			token, err := jwt.ParseWithClaims(r.Header["Token"][0], claims, func(token *jwt.Token) (interface{}, error) {

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
					fmt.Println("claimsfwfw : ", claims)
					
					mpData := context.WithValue(r.Context(), "editUser" , mp)
					
					endpoint(w, r.WithContext(mpData))
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
