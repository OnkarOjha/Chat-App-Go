package server

import (
	chat "main/ChatServer"
	"net/http"
	server "main/Utils"
	
	
)

func Namespaces() {
	
	// connecting client to the server
	server.Server.OnConnect("/", chat.TokenCheck(chat.Connect))
	server.Server.OnEvent("/","join", chat.RoomJoin)
	server.Server.OnEvent("/","message", chat.Messages)
	

	//socket server
	http.Handle("/socket.io/", server.Server)
	
	
}
