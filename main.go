package main

import (
	"applinh/gogrpcudemy/calculator/calculator_client"
	"applinh/gogrpcudemy/calculator/calculator_server"
	"applinh/gogrpcudemy/greet/greet_client"
	"applinh/gogrpcudemy/greet/greet_server"
	"fmt"
	"os"
)

func main() {

	a := os.Args
	if len(a) == 2 {
		switch a[1] {
		case "greet_server":
			greet_server.StartServer()
		case "greet_client":
			greet_client.StartClient()
		case "calculator_server":
			calculator_server.StartServer()
		case "calculator_client":
			calculator_client.StartClient()
		default:
			fmt.Println("Service not found.")
		}
	} else {
		fmt.Println("Please provide an argument.")
	}

}
