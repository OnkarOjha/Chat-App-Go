package Controllers

import (
	"context"
	"encoding/json"
	"fmt"
	db "main/Database"
	models "main/Models"
	response "main/Response"
	commonFunctions "main/Utils"
	constants "main/Utils"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

//Token authorization check middleware
func IsAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("inside middleware")
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
		fmt.Println("cookies token: ", cookie.Value)
		var mp = make(map[string]interface{})
		json.NewDecoder(r.Body).Decode(&mp)
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
		//bring the token from DB
		var user models.User
		db.DB.Raw("select * from users where user_id=?",userId).Scan(&user)
		
		// check if token is blacklisted
		if commonFunctions.BlacklistTokenHandler(cookie.Value) {
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

		if cookie.Value != "" {

			token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {

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

				// Token expiration check 
				if claims.ExpiresAt.Before(time.Now().Add(2 * time.Minute)) {
					fmt.Println("refresh handler called")
					var userToRefresh models.User
					db.DB.Raw("select * from users where user_id=?",claims.User_id).Scan(&userToRefresh)
					fmt.Println("BBBB:",userToRefresh)
					newTokenString :=commonFunctions.GenerateJwtToken(userToRefresh , claims.Phone , w)
					fmt.Println("AAAA",newTokenString)
					db.DB.Model(&models.User{}).Where("user_id=?", claims.User_id).Update("token", newTokenString)
					commonFunctions.SetCookie(w , r , newTokenString)
				}

				if claims.User_id == user.User_Id {
					fmt.Println("claimsfwfw : ", claims)

					mpData := context.WithValue(r.Context(), "editUser", mp)

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


//Token Expiration check middleware
func TokenExpirationCheck(originalHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("inside token expiration time check")
		TokenFromHeader := r.Header["Token"][0]
		fmt.Println("token header in expiration check", TokenFromHeader)
		claims , err := commonFunctions.DecodeToken(TokenFromHeader)
		if err!=nil{
			response.ShowResponse(
				"Failure",
				400,
				"Invalid token",
				"",
				w,
			)
			return
		}
		now := time.Now()
		if claims.RegisteredClaims.ExpiresAt.Before(now.Add(constants.TokenCheckTimer)) {
			// blacklist the token from header

			// new token generation
			userId := claims.User_id
			phone := claims.Phone
			claims := models.Claims{
				User_id: userId,
				Phone: phone,
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(constants.TokenExpirationDuration),
				},
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString([]byte(os.Getenv("JWTKEY")))
			if err != nil {
				response.ShowResponse(
					"Failure",
					401,
					"Unable to sign the JWT",
					"",
					w,
				)
				return
			}
			// update the DB with new token
			var user models.User
			db.DB.Model(&models.User{}).Where("phone=?",phone).Update("token" , tokenString)
			db.DB.Raw("select * from users where phone=?",phone).Scan(&user)
			response.ShowResponse(
				"Success",
				200,
				"New Token issues successfully",
				user,
				w,
			)
			originalHandler.ServeHTTP(w, r)
		}else{
			originalHandler.ServeHTTP(w, r)
		}

		
	})

}


