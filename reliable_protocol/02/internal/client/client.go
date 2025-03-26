package main

import (
	"fmt"
	"go-reliable/shared"
	"net"
	"os"
	"time"
)

type User struct {
	conn      *net.UDPConn
	server    *net.UDPAddr
	SeqNumber int
}

func (User *User) ReceiveMessage(data []byte) {

	receivedSeqNumber := int(data[0])

	// 解析 Data
	receivedData := string(data[1:])

	// 构造 Packet
	packet := shared.Packet{
		SeqNum: receivedSeqNumber,
		Data:   receivedData,
	}

	ackPacket := shared.Packet{
		SeqNum: User.SeqNumber,
		Data:   "",
	}

	if receivedSeqNumber == User.SeqNumber {
		fmt.Printf("收到数据包: Seq=%d, Data=%s\n", packet.SeqNum, packet.Data)
		User.SeqNumber++
	}
	time.Sleep(2 * time.Second)
	shared.SendUDPPacketConnected(User.conn, ackPacket)

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
	serverAddress := "127.0.0.1:8080" // 服务器地址
	user, err := ConnectToServer(serverAddress)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}
	defer CloseConnection(user.conn)

	initPacket := shared.Packet{
		SeqNum: 0,
		Data:   "",
	}
	shared.SendUDPPacketConnected(user.conn, initPacket)

	for {
		// 接收客户端请求
		buf := make([]byte, 1024)
		_, _, err := user.conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}

		user.ReceiveMessage(buf)

	}

}
