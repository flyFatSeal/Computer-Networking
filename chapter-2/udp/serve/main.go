package main

import (
	"fmt"
	"net"
)

func main() {
	serverPort := "12000"
	serverAddress, err := net.ResolveUDPAddr("udp", ":"+serverPort)
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	serverSocket, err := net.ListenUDP("udp", serverAddress)
	if err != nil {
		fmt.Println("Error creating socket:", err)
		return
	}
	defer serverSocket.Close()

	fmt.Println("The server is ready to receive")

	buffer := make([]byte, 2048)

	for {
		n, clientAddress, err := serverSocket.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}

		message := string(buffer[:n])
		modifiedMessage := []byte(message)

		for i := range modifiedMessage {
			if modifiedMessage[i] >= 'a' && modifiedMessage[i] <= 'z' {
				modifiedMessage[i] -= 'a' - 'A'
			}
		}

		_, err = serverSocket.WriteToUDP(modifiedMessage, clientAddress)
		if err != nil {
			fmt.Println("Error sending to UDP:", err)
		}
		fmt.Println("sending to client", modifiedMessage)
	}
}
