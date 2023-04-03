package server

import (
	chat "main/ChatServer"
	"net/http"
	server "main/Utils"
	// socketio "github.com/googollee/go-socket.io"
	
)

func Namespaces() {
	
	// connecting client to the server
	server.Server.OnConnect("/", chat.Connect)
	server.Server.OnEvent("/","join", chat.RoomJoin)
	server.Server.OnEvent("/","message", chat.Messages)
	

	//socket server
	http.Handle("/socket.io/", server.Server)
	
	
}
