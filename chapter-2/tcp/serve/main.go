package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	serverPort := "12000"
	serverAddress, err := net.ResolveTCPAddr("tcp", ":"+serverPort)
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	listener, err := net.ListenTCP("tcp", serverAddress)
	if err != nil {
		fmt.Println("Error creating listener:", err)
		return
	}
	defer listener.Close()

	fmt.Println("The server is ready to receive")

	for {
		connectionSocket, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleClient(connectionSocket)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from connection:", err)
		return
	}

	sentence := string(buffer[:n])
	capitalizedSentence := strings.ToUpper(sentence)
	_, err = conn.Write([]byte(capitalizedSentence))
	if err != nil {
		fmt.Println("Error writing to connection:", err)
	}
}
