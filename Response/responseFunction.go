package Response

import (
	"encoding/json"
	"fmt"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
)

func ShowResponse(status string, statusCode int64, message string, data interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(int(statusCode))
	response := Response{
		Status:  status,
		Code:    statusCode,
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(&response)
}



func SocketResponse(data interface{}, message string, s socketio.Conn) {
	socketResponse := Socket{
		Message: message,
		Data:    data,
	}
	s.Emit("reply", socketResponse, func() {
		fmt.Println("acknowledgement sent to client")
	})

}


