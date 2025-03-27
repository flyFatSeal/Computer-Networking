package main

import (
	"encoding/json"
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

func (User *User) ReceiveMessage(data []byte) {
	var packet shared.Packet
	if err := json.Unmarshal(data, &packet); err != nil {
		fmt.Printf("收到无效的 packet: %v\n", err)
		return
	}

	// 检查数据包是否损坏
	if shared.IsCorrupted(packet) {
		fmt.Printf("收到损坏的 packet: Seq=%d\n", packet.SeqNum)
		User.sendAck(User.SeqNumber - 1) // 重新发送上一个 ACK
		return
	}

	// 检查序列号是否匹配
	if packet.SeqNum == User.SeqNumber {
		// 正确的数据包，提取数据并更新期望的序列号
		fmt.Printf("收到数据包: Seq=%d, Data=%s\n", packet.SeqNum, packet.Data)
		User.sendAck(User.SeqNumber) // 发送 ACK
		User.SeqNumber++             // 更新期望的序列号
	} else {
		// 乱序数据包，丢弃并重新发送上一个 ACK
		fmt.Printf("收到乱序数据包: Seq=%d, 期望 Seq=%d\n", packet.SeqNum, User.SeqNumber)
		User.sendAck(User.SeqNumber - 1) // 重新发送上一个 ACK
	}
}

// 发送 ACK 包
func (User *User) sendAck(ackNum int) {
	ackPacket := shared.Packet{
		SeqNum: ackNum,
		AckNum: ackNum,
		Data:   "ack_",
	}

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
		Data:   "init-client",
	}
	shared.SendUDPPacketConnected(user.conn, initPacket)

	for {
		// 接收客户端请求
		buf := make([]byte, 1024)
		n, _, err := user.conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}

		user.ReceiveMessage(buf[:n])

	}

}
