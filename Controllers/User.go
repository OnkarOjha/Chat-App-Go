package Controllers

import (
	"encoding/json"
	"fmt"
	db "main/Database"
	models "main/Models"
	response "main/Response"
	constants "main/Utils"
	validator "main/Validation"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/verify/v2"
)

var TwilioClient *twilio.RestClient

func TwilioInit(password string)  {
	TwilioClient = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username:constants.TWILIO_ACCOUNT_SID,
		Password: password,
	})


}

// send OTP to user
func SendOtpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	EnableCors(&w)

	var mp = make(map[string]interface{})

	json.NewDecoder(r.Body).Decode(&mp)

	// validator
	err := validator.CheckValidation(mp["phone"])
	if err != nil {
		response.ShowResponse(
			"Failure",
			400,
			"",
			err.Error(),
			w,
		)
		return
	}
	// check if number already exists and user is not acitve then he is logging in
	var userlogout bool
	db.DB.Raw("SELECT EXISTS(SELECT * FROM users WHERE phone=? and is_active=false)", mp["phone"]).Scan(&userlogout)
	fmt.Println("userlogout: ", userlogout)
	if userlogout {
		fmt.Println("user logging in")
		ok, sid := sendOtp("+91"+mp["phone"].(string), w)
		if ok {
			response.ShowResponse(
				"Success",
				200,
				"OTP sent successfully",
				sid,
				w,
			)
		}
		return
	}

	// Check if number already exists for new user sigining in
	var numberExists bool
	err = db.DB.Raw("SELECT EXISTS(SELECT 1 FROM users WHERE phone=? and is_active=true)", mp["phone"]).Scan(&numberExists).Error
	if err != nil {
		panic(err)
	}

	// Response
	if numberExists {
		response.ShowResponse(
			"Failure",
			409,
			"Number already exists",
			"",
			w,
		)
		return
	}

	ok, sid := sendOtp("+91"+mp["phone"].(string), w)
	if ok {
		response.ShowResponse(
			"Success",
			200,
			"OTP sent successfully",
			sid,
			w,
		)
	}
	return

}

// function to send OTP while user registration
func sendOtp(to string, w http.ResponseWriter) (bool, *string) {
	params := &openapi.CreateVerificationParams{}
	params.SetTo(to)

	params.SetChannel("sms")

	resp, err := TwilioClient.VerifyV2.CreateVerification(os.Getenv("VERIFY_SERVICE_SID"), params)

	if err != nil {
		response.ShowResponse(
			"API ERROR",
			401,
			"No credentials provided",
			"",
			w,
		)
		return false, nil
	} else {

		return true, resp.Sid
	}

}

// Check OTP status
func VerifyOTPHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	EnableCors(&w)

	var otp = make(map[string]interface{})
	json.NewDecoder(r.Body).Decode(&otp)

	_, ok1 := otp["phone"]
	_, ok2 := otp["otp"]

	if !ok1 || !ok2 {
		response.ShowResponse(
			"Failure",
			400,
			"",
			"Either OTP or Phone missing",
			w,
		)
		return

	}

	err := validator.CheckValidation(otp["phone"])
	if err != nil {
		response.ShowResponse(
			"Failure",
			400,
			"",
			err.Error(),
			w,
		)
		return
	}

	if CheckOtp("+91"+otp["phone"].(string), otp["otp"].(string), w) {
		u := UserSignupHandler(w, r, otp["phone"].(string))
		response.ShowResponse(
			"Success",
			200,
			"Phone Number Verified Successfully",
			u,
			w,
		)

	} else {
		// fmt.Println("Verifictaion failed")
		response.ShowResponse(
			"Not Found",
			404,
			"OTP ERROR",
			"",
			w,
		)
		return
	}
}

// OTP code verification
func CheckOtp(to string, code string, w http.ResponseWriter) bool {
	params := &openapi.CreateVerificationCheckParams{}
	params.SetTo(to)
	params.SetCode(code)
	resp, err := TwilioClient.VerifyV2.CreateVerificationCheck(os.Getenv("VERIFY_SERVICE_SID"), params)

	if err != nil {
		return false
	} else if *resp.Status == "approved" {
		return true
	} else {
		return false
	}
}

// User Registeration
func UserSignupHandler(w http.ResponseWriter, r *http.Request, phone string) models.User {
	w.Header().Set("Content-Type", "application/json")

	fmt.Println("We are making user entries....")
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	//check if we are creating token for a new user or that user already exists and he is logging in
	var userexists bool
	db.DB.Raw("SELECT EXISTS(SELECT * FROM users WHERE phone=? and is_active=false)", phone).Scan(&userexists)
	fmt.Println("token generation time user exists or not : ", userexists)
	if userexists {
			
		var user1 models.User
		db.DB.Raw("select * from users where phone=?",phone).Scan(&user1)
		
		// JWT token authentication
		TokenString := GenerateJwtToken(user1,phone,w)

		// Updation in DB
		db.DB.Model(&models.User{}).Where("phone=?",phone).Update("token" , TokenString)
		db.DB.Model(&models.User{}).Where("phone=?",phone).Update("is_active" , true)

		return user1

	}

	dateStr := time.Now().Truncate(time.Hour)
	user.Join_date = dateStr.Format("02 Jan 2006")
	user.Phone = phone
	user.Is_active = true
	user.Is_verified = true
	db.DB.Create(&user)

	// jwt authentication token

	TokenString := GenerateJwtToken(user , phone , w)

	db.DB.Model(&user).Where("user_id=?", user.User_Id).Updates(&models.User{Token: TokenString})

	return user

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

	err := validator.CheckValidationStruct(edituser)
	if err != nil {
		response.ShowResponse(
			"Failure",
			400,
			"",
			err.Error(),
			w,
		)
		return
	}

	var user models.User
	db.DB.Model(&models.User{}).Where("user_id=?", edituser.User_Id).Find(&user)

	result := db.DB.Model(&models.User{}).Where("user_id=?", edituser.User_Id).Updates(&edituser)
	var showUser models.User
	db.DB.Raw("SELECT * from users where user_id=?", edituser.User_Id).Scan(&showUser)

	if result.Error != nil {
		response.ShowResponse(
			"Failure",
			500,
			"DB error",
			"",
			w,
		)
		return
	} else if result.RowsAffected == 0 {
		db.DB.Create(&edituser)
		response.ShowResponse(
			"Success",
			200,
			"User added successfully",
			edituser,
			w,
		)
	} else {

		response.ShowResponse(
			"Success",
			200,
			"Old user updated successfully",
			showUser,
			w,
		)
	}

}

func GenerateJwtToken(user models.User , phone string , w http.ResponseWriter) string {
	//create user claims
	claims := models.Claims{
		User_id: user.User_Id,
		Phone:   phone,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(constants.TokenExpirationDuration),
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
