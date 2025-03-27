package main

import (
	"encoding/json"
	"fmt"
	"go-reliable/shared"
	"net"
	"testing"
	"time"
)

func BenchmarkServerHandleClient(b *testing.B) {
	server := NewServer()
	go server.Start("8080") // 启动服务器
	defer func() {
		server.mapLok.Lock()
		server.user = make(map[string]User)
		server.mapLok.Unlock()
	}()

	// 模拟多个客户端
	numClients := 10 // 模拟的客户端数量
	var clients []net.Conn

	// 创建多个客户端连接
	for i := 0; i < numClients; i++ {
		clientAddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:0")
		serverAddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:8080")
		conn, err := net.DialUDP("udp", clientAddr, serverAddr)
		if err != nil {
			b.Fatalf("客户端 %d 无法连接到服务器: %v", i, err)
		}
		defer conn.Close()
		clients = append(clients, conn)

		// 客户端发送初始化消息
		initPacket := shared.Packet{
			SeqNum:   0,
			Data:     fmt.Sprintf("Client %d: Init", i),
			Checksum: 0,
		}
		initPacket.Checksum = shared.CalculateChecksum(initPacket)

		data, _ := json.Marshal(initPacket)
		_, err = conn.Write(data)
		if err != nil {
			b.Fatalf("客户端 %d 初始化消息发送失败: %v", i, err)
		}
	}

	// 重置计时器
	b.ResetTimer()

	// 模拟服务端向客户端发送数据
	for i := 0; i < b.N; i++ {
		for clientID, conn := range clients {
			go func(clientID int, conn net.Conn) {
				packet := shared.Packet{
					SeqNum:   i,
					Data:     fmt.Sprintf("Server to Client %d: Message %d", clientID, i),
					Checksum: 0,
				}
				packet.Checksum = shared.CalculateChecksum(packet)

				data, _ := json.Marshal(packet)

				// 服务端发送数据包
				_, err := conn.Write(data)
				if err != nil {
					b.Fatalf("服务端发送数据包失败: %v", err)
				}

				// 客户端接收数据
				buf := make([]byte, 1024)
				conn.SetReadDeadline(time.Now().Add(2 * time.Second))
				_, _, err = conn.ReadFromUDP(buf)
				if err != nil {
					b.Logf("客户端 %d 接收数据超时: %v", clientID, err)
				}
			}(clientID, conn)
		}
	}
}
