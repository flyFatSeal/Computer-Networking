package main

import (
	"fmt"
	"go-reliable/shared"
	"net"
	"os"
)

type User struct {
	conn      *net.UDPConn
	server    *net.UDPAddr
	SeqNumber int
}

func (User *User) ReceiveMessage() {
	buf := make([]byte, 1024)
	n, _, err := User.conn.ReadFromUDP(buf)
	if err != nil {
		fmt.Println(err)
	}

	receivedSeqNumber := int(buf[0])
	if receivedSeqNumber == User.SeqNumber {
		fmt.Printf("收到 SeqNumber: Seq=%d\n", receivedSeqNumber)
		User.SeqNumber++

	}

	packet := shared.Packet{
		SeqNumber: receivedSeqNumber,
		Data:      "",
	}

	shared.SendUDPPacket(User.conn, User.server, packet)

}

func ConnectToServer(address string) (*User, error) {
	// 解析服务器地址
	serverAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve server address: %v", err)
	}

	// 创建 UDP 连接
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %v", err)
	}

	return &User{
		conn:      conn,
		server:    serverAddr,
		SeqNumber: 0,
	}, nil
}

func CloseConnection(conn net.Conn) {
	conn.Close()
}

func main() {
	serverAddress := "192.168.1.168:8080" // 服务器地址
	user, err := ConnectToServer(serverAddress)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}
	defer CloseConnection(user.conn)
	initPacket := shared.Packet{
		SeqNumber: 0,
		Data:      "",
	}
	shared.SendUDPPacket(user.conn, user.server, initPacket)

	for {
		// 接收客户端请求
		buf := make([]byte, 1024)
		_, addr, err := user.conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}

		user.ReceiveMessage()

	}

}
