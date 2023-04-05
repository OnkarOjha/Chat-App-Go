package main

import (
	"fmt"
	route "main/Routes"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Chat Application Backend")
	envErr:=godotenv.Load(".env")
	if envErr!=nil {

		fmt.Println("could not load environment")
	}
	// to call API
	route.Routes()

	

}
