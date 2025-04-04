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
	TimeoutDuration = 10 * time.Second // 超时时间
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

	// 为每个用户维护一个消息队列
	messageQueues := make(map[string]chan []byte)

	for {
		// 接收客户端请求
		buf := make([]byte, 1024)
		n, addr, err := listener.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}

		clientAddr := addr.String()
		Server.mapLok.Lock()
		_, exists := Server.user[clientAddr]
		if !exists {
			// 新用户，创建消息队列并启动处理协程
			Server.user[clientAddr] = User{addr: clientAddr}
			messageQueues[clientAddr] = make(chan []byte, 10)
			go Server.HandleClient(listener, addr, messageQueues[clientAddr])
		}
		Server.mapLok.Unlock()

		// 将消息发送到对应用户的消息队列
		if exists {
			messageQueues[clientAddr] <- buf[:n]
		}
	}
}

func (Server *Server) HandleClient(conn *net.UDPConn, addr *net.UDPAddr, messageQueue chan []byte) {
	currentSeq := 0
	ackCh := make(chan bool)
	timeout := time.NewTimer(TimeoutDuration)
	data := []string{"1", "2", "3", "4", "5", "6", "7"}

	// 处理消息队列中的消息
	go func() {
		for msg := range messageQueue {
			if int(msg[0]) == currentSeq {
				fmt.Printf("收到来自 %s 的 ACK: Seq=%d\n", addr.String(), currentSeq)
				ackCh <- true
			}
		}
	}()

	for currentSeq < len(data) {
		packet := shared.Packet{
			SeqNum: currentSeq,
			Data:   data[currentSeq],
		}
		time.Sleep(2 * time.Second)
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
	Server := NewServer()
	Server.Start("8080")
}
