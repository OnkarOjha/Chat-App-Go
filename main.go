package main

import (
	"fmt"
	twilio "main/Controllers"
	route "main/Routes"
	"os"

	"github.com/joho/godotenv"
)

// @title Chat-Application
// @version 1.0.0
// @description API Documentation for Chat-Application
// @host localhost:8000
func main() {
	fmt.Println("Chat Application Backend")
	envErr := godotenv.Load(".env")
	twilio.TwilioInit(os.Getenv("TWILIO_AUTH_TOKEN"))
	if envErr != nil {

		fmt.Println("could not load environment")
	}
	// to call API
	route.Routes()

}
