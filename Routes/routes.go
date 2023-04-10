package Routes

import (
	"fmt"
	"log"
	controller "main/Controllers"
	db "main/Database"
	"net/http"
	server "main/Utils"
	namespace "main/Server"
)

func Routes() {
	fmt.Println("Listening on port 8000")

	err := db.Connect()
	if err != nil {
		panic(err)
	}

	//room-topic end-points
	http.HandleFunc("/roomTopic", controller.RoomTopicController)
	http.HandleFunc("/getroomTopic", controller.RoomTopicGetter)

	//user end-points
	http.HandleFunc("/sendOtp", controller.SendOtpHandler)
	http.HandleFunc("/verifyOtp", controller.VerifyOTPHandler)
	// http.HandleFunc("/userSignup", controller.UserSignupHandler)
	http.HandleFunc("/getUser", controller.UserGetterHandler)
	http.Handle("/editUser",controller.IsAuthorized(controller.UserEditHandler))
	http.Handle("/logout", controller.IsAuthorized(controller.LogoutHandler))
	http.Handle("/deleteAccount", controller.IsAuthorized(controller.DeleteAccount))
	http.HandleFunc("/userRoomInfo", controller.UserRoomsDetails)



	

	// chat-room functions
	http.HandleFunc("/participants",controller.ParticipantDetails)
	http.HandleFunc("/rooms",controller.RoomDetails)
	http.HandleFunc("/messages",controller.MessageDetails)
	http.HandleFunc("/roomDelete",controller.RoomDelete)
	http.HandleFunc("/messageSearch",controller.MessageSearchController)






	// to call socket-io
	namespace.Namespaces()
	go server.Server.Serve()
	defer server.Server.Close()

	http.HandleFunc("/", controller.HomeHandler)


	log.Fatal(http.ListenAndServe(":8000", nil))
}
