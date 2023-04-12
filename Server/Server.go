package server

import (
	chat "main/ChatServer"
	server "main/Utils"
)

func Namespaces() {

	// connecting client to the server
	server.Server.OnConnect("/", chat.TokenCheck(chat.Connect))
	server.Server.OnEvent("/", "createroom", chat.RoomCreate)
	server.Server.OnEvent("/", "join", chat.RoomJoin)
	server.Server.OnEvent("/", "message", chat.Messages)
	server.Server.OnEvent("/", "leave", chat.RoomLeave)

	//socket server
	server.Mux.Handle("/socket.io/", server.Server)

}
