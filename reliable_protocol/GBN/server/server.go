package main

import (
	"encoding/json"
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
			messageQueues[clientAddr] = make(chan []byte, 1024)
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
	ackCh := make(chan int)
	timeout := time.NewTimer(TimeoutDuration)
	window_size := 8
	base := 0
	nextSeq := 0

	data := make([]string, 1000)
	for i := 1; i <= 1000; i++ {
		data[i-1] = fmt.Sprintf("%d", i)
	}

	// 处理消息队列中的消息
	go func() {
		for msg := range messageQueue {
			fmt.Printf("从消息队列中取出数据: %s\n", string(msg[0:]))
			var packet shared.Packet
			err := json.Unmarshal(msg, &packet) // 将字节数据解析为 Packet
			if err != nil {
				fmt.Printf("收到无效的 packet: Seq=%d\n", packet.SeqNum)
				continue
			}
			if shared.IsCorrupted(packet) {
				fmt.Printf("收到损坏的 packet: Seq=%d\n", packet.SeqNum)
				continue
			}
			ackCh <- packet.AckNum
		}
	}()

	for base < len(data)-1 {
		timeout.Reset(TimeoutDuration)
		for nextSeq < base+window_size {
			if nextSeq > len(data)-1 {
				break
			}
			packet := shared.Packet{
				SeqNum: nextSeq,
				Data:   data[nextSeq],
			}
			nextSeq++
			shared.SendUDPPacket(conn, addr, packet)
		}
		select {
		case ackSeq := <-ackCh:
			// 检查 ACK 是否在窗口范围内
			if ackSeq >= base && ackSeq < nextSeq {
				fmt.Printf("收到 ACK: Seq=%d\n", ackSeq)
				base = ackSeq + 1 // 更新窗口起点
			} else {
				fmt.Printf("收到无效的 ACK: Seq=%d\n", ackSeq)
			}
		case <-timeout.C:
			nextSeq = base
			fmt.Println("超时，重传数据包")
		}
	}
	timeout.Stop()

}

func main() {
	Server := NewServer()
	Server.Start("8080")
}
