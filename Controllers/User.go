package Controllers

import (
	"encoding/json"
	"fmt"
	db "main/Database"
	models "main/Models"
	cons "main/Utils"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/verify/v2"
)

// twilio client interface
var client *twilio.RestClient = twilio.NewRestClientWithParams(twilio.ClientParams{
	Username: cons.TWILIO_ACCOUNT_SID,
	Password: cons.TWILIO_AUTH_TOKEN,
})

// send OTP to user
func SendOtpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	EnableCors(&w)

	// phone := r.FormValue("phone")

	// fmt.Println("phone: ", phone)
	var mp = make(map[string]string)
	json.NewDecoder(r.Body).Decode(&mp)
	fmt.Println("phone: ", mp["phone"])
	sendOtp("+91" + mp["phone"])

}

// function to send OTP while user registration
func sendOtp(to string) {
	params := &openapi.CreateVerificationParams{}
	params.SetTo(to)

	params.SetChannel("sms")

	resp, err := client.VerifyV2.CreateVerification(cons.VERIFY_SERVICE_SID, params)

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("Sent verification '%s'\n", *resp.Sid)
	}
}

// Check OTP status
func VerifyOTPHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	EnableCors(&w)

	var otp = make(map[string]string)
	json.NewDecoder(r.Body).Decode(&otp)
	fmt.Println("phone: ", otp["phone"])
	fmt.Println("otp is:", otp["otp"])
	if CheckOtp("+91"+otp["phone"], otp["otp"]) {
		fmt.Println("Phone Number verified sucessfully")
		UserSignupHandler(w, r, otp["phone"])
	} else {
		fmt.Println("Verifictaion failed")
	}
}

// OTP code verification
func CheckOtp(to string, code string) bool {
	params := &openapi.CreateVerificationCheckParams{}
	params.SetTo(to)
	params.SetCode(code)
	resp, err := client.VerifyV2.CreateVerificationCheck(cons.VERIFY_SERVICE_SID, params)

	if err != nil {
		fmt.Println("Error is :", err)
		return false
	} else if *resp.Status == "approved" {
		return true
	} else {
		return false
	}
}

// User Registeration
func UserSignupHandler(w http.ResponseWriter, r *http.Request, phone string) {
	w.Header().Set("Content-Type", "application/json")

	fmt.Println("We are making user entries....")
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	dateStr := time.Now().Truncate(time.Hour)
	user.Join_date = dateStr.Format("02 Jan 2006")
	user.Phone = phone
	user.Is_active = true

	db.DB.Create(&user)

	// jwt authentication token

	//create user claims
	claims := models.Claims{
		User_id: user.User_Id,
		Phone:   phone,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(cons.TokenExpirationDuration),
		},
	}
	fmt.Println("claims: ", claims)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Println("token: ", token)
	tokenString, err := token.SignedString(cons.JwtKey)
	if err != nil {
		fmt.Println("error is :", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println("token string :", tokenString)

	db.DB.Model(&user).Where("user_id=?", user.User_Id).Updates(&models.User{Token: tokenString})

	json.NewEncoder(w).Encode(&user)

}

func UserGetterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("You are viewing all the registered user")
	var users []models.User
	query := "SELECT * FROM users"
	db.DB.Raw(query).Scan(&users)
	json.NewEncoder(w).Encode(&users)
}

// user details edit
func UserEditHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("You are editing user information.....")

	var edituser models.User
	json.NewDecoder(r.Body).Decode(&edituser)
	
	// verify token
	var user models.User
	db.DB.Model(&models.User{}).Where("user_id=?",edituser.User_Id).Find(&user)
	claims := &models.Claims{}
	fmt.Println("user", user)
	fmt.Println("token: ",user.Token)
	parsedToken ,err := jwt.ParseWithClaims(user.Token ,claims, func(token *jwt.Token) (interface{}, error) {
		if _,ok:=token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil,fmt.Errorf("error")
		}
		return cons.JwtKey , nil
	})
	if err != nil || !parsedToken.Valid {
		http.Error(w, "Invalid or expired token", http.StatusBadRequest)
		return
	}
	

	result := db.DB.Model(&models.User{}).Where("user_id=?", edituser.User_Id).Updates(&edituser)
	
	if result.Error != nil {
		fmt.Println("DB error")
		return
	} else if result.RowsAffected == 0 {
		db.DB.Create(&edituser)
		fmt.Fprintf(w, "New user added")
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Old user updated successfully"))
	}

	json.NewEncoder(w).Encode(&edituser)

}