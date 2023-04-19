package chat

import (
	response "main/Response"

	socketio "github.com/googollee/go-socket.io"
)

type Call struct {
	To         string
	FromNumber string
	CallStatus string
	CallSid    string
}

func VoiceCall(s socketio.Conn, data map[string]interface{}) {
	to := data["to"].(string)
	fromNumber := data["from"].(string)
	callStatus := data["callStatus"].(string)
	callSid := data["callSid"].(string)

	call := &Call{
		To:         to,
		FromNumber: fromNumber,
		CallStatus: callStatus,
		CallSid:    callSid,
	}

	response.SocketResponse(call, "", s)

}
