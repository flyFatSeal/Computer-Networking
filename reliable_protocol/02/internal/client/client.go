package client

import (
	"bufio"
	"fmt"
	"net"
)

func ConnectToServer(address string) (net.Conn, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func SendMessage(conn net.Conn, message string) error {
	_, err := fmt.Fprintf(conn, message+"\n")
	return err
}

func ReceiveMessage(conn net.Conn) (string, error) {
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return "", err
	}
	return message, nil
}

func CloseConnection(conn net.Conn) {
	conn.Close()
}

func main() {
	address := "localhost:8080" // Example address
	conn, err := ConnectToServer(address)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer CloseConnection(conn)

	// Example usage
	err = SendMessage(conn, "Hello, Server!")
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}

	response, err := ReceiveMessage(conn)
	if err != nil {
		fmt.Println("Error receiving message:", err)
		return
	}
	fmt.Println("Received from server:", response)
}
