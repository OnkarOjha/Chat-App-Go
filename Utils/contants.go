package Constants

import (
	"time"

	socketio "github.com/googollee/go-socket.io"
)

const Host = "localhost"
const Port = "5432"
const User = "postgres"
const Password = "Onkar17082001"
const Dbname = "chat_app"

// // // twilio credentials
// const TWILIO_ACCOUNT_SID string = "AC5869424c6ae2d0b27f66a5d5b9b90485"
// const TWILIO_AUTH_TOKEN string = "80c7112f1d2297548a6caec62c8bab81"
// const VERIFY_SERVICE_SID string = "VA6602f8535f8f1369b0ed68eed5d6af67"

// JWT secret key
var JwtKey = []byte("My_Key")

// server instance
var Server = socketio.NewServer(nil)

// token timing
var TokenExpirationDuration = time.Now().Add( 365 * 24 * time.Hour)
var TokenRefreshInterval = 1 * time.Minute
