package Routes

import (
	"fmt"
	"log"
	controller "main/Controllers"
	db "main/Database"
	namespace "main/Server"
	server "main/Utils"
	_ "main/docs"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

func Routes() {
	fmt.Println("Listening on port 8000")

	err := db.Connect()
	if err != nil {
		panic(err)
	}

	// cors middleware
	server.Mux.Use(controller.CORSMiddleware)

	//room-topic end-points
	server.Mux.HandleFunc("/roomTopic", controller.RoomTopicController)
	server.Mux.HandleFunc("/getroomTopic", controller.RoomTopicGetter)

	//user end-points
	server.Mux.HandleFunc("/sendOtp", controller.SendOtpHandler)
	server.Mux.HandleFunc("/verifyOtp", controller.VerifyOTPHandler)
	server.Mux.HandleFunc("/getUser", controller.UserGetterHandler)
	server.Mux.Handle("/editUser", controller.IsAuthorized(controller.UserEditHandler))
	server.Mux.Handle("/logout", controller.IsAuthorized(controller.LogoutHandler))
	server.Mux.Handle("/deleteAccount", controller.IsAuthorized(controller.DeleteAccount))
	server.Mux.HandleFunc("/userRoomInfo", controller.UserRoomsDetails)

	// chat-room functions
	server.Mux.HandleFunc("/participants", controller.ParticipantDetails)
	server.Mux.HandleFunc("/rooms", controller.RoomDetails)
	server.Mux.HandleFunc("/messages", controller.MessageDetails)
	server.Mux.HandleFunc("/roomDelete", controller.RoomDelete)
	server.Mux.HandleFunc("/messageSearch", controller.MessageSearchController)

	//swagger api
	server.Mux.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	//fileUpload
	server.Mux.HandleFunc("/upload", controller.FileUpload)

	// to call socket-io
	namespace.Namespaces()
	go server.Server.Serve()
	defer server.Server.Close()

	// Serve files from a directory
	fs := http.FileServer(http.Dir("/home/chicmic/Desktop/Chat-App-Go/File/uploads"))
	server.Mux.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", fs))


	log.Fatal(http.ListenAndServe(":8000", server.Mux))
}
