package main

import (
	"fmt"
	"go-reliable/shared"
	"net"
	"os"
	"sync"
	"time"
)

const (
	TimeoutDuration = 2 * time.Second // 超时时间
)

type User struct {
	addr string
}
type Server struct {
	user   map[string]User
	mapLok sync.RWMutex
}

func NewServer() *Server {
	server := &Server{
		user: make(map[string]User),
	}

	return server
}

func (Server *Server) Start(port string) {
	serverAddr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf(":%s", port))
	listener, err := net.ListenUDP("udp", serverAddr)
	defer listener.Close()

	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}

	for {
		// 接收客户端请求
		buf := make([]byte, 1024)
		_, addr, err := listener.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}

		clientAddr := addr.String()
		Server.mapLok.Lock()
		_, exists := Server.user[clientAddr]
		if !exists {
			Server.user[clientAddr] = User{addr: clientAddr}
			go Server.SendMessage(listener, addr)
		}
		Server.mapLok.Unlock()
	}

}

func (Server *Server) SendMessage(conn *net.UDPConn, addr *net.UDPAddr) {
	currentSeq := 0
	ackCh := make(chan bool)
	timeout := time.NewTimer(TimeoutDuration)
	data := []string{"1", "2", "3", "4", "5", "6", "7"}

	go func() {
		buf := make([]byte, 1024)
		for {
			_, remoteAddr, err := conn.ReadFromUDP(buf)
			if err != nil {
				continue
			}
			// 确保消息来自当前客户端
			if remoteAddr.String() != addr.String() {
				continue // 忽略其他客户端的消息
			}
			receivedSeq := int(buf[0])
			if receivedSeq == currentSeq {
				fmt.Printf("收到 ACK: Seq=%d\n", receivedSeq)
				ackCh <- true
			}
		}
	}()

	for currentSeq < len(data) {
		packet := shared.Packet{
			SeqNumber: currentSeq,
			Data:      data[currentSeq],
		}

		shared.SendUDPPacket(conn, addr, packet)
		timeout.Reset(TimeoutDuration)

		select {
		case <-ackCh:
			currentSeq++ // 收到正确的 ACK，切换到下一个序列号
		case <-timeout.C:
			fmt.Println("超时，重传数据包")
		}
	}
}

func main() {
	// Set up the server to listen on a specific port
	serverAddr, _ := net.ResolveUDPAddr("udp", ":8080")

	listener, err := net.ListenUDP("udp", serverAddr)
	defer listener.Close()

	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}

	fmt.Println("Server is listening on port 8080...")

	Server := NewServer()

	Server.Start("8080")

}
