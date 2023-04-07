package Constants

import (
	"time"

	socketio "github.com/googollee/go-socket.io"
)

const TWILIO_ACCOUNT_SID string = "AC5869424c6ae2d0b27f66a5d5b9b90485"
const VERIFY_SERVICE_SID string = "VA6602f8535f8f1369b0ed68eed5d6af67"

// server instance
var Server = socketio.NewServer(nil)

// token timing
var TokenExpirationDuration = time.Now().Add( 365 * 24 * time.Hour)
var TokenRefreshInterval = 1 * time.Minute

