package server

import (
	chat "main/ChatServer"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
	
)
var Server = socketio.NewServer(nil)
// all the socket configurations are here
func Namespaces() {
	
	// connecting client to the server
	Server.OnConnect("/", chat.Connect)
	Server.OnEvent("/","join", chat.RoomJoin)
	//socket server
	http.Handle("/socket.io/", Server)
	// go routines to operate the server
	
}
