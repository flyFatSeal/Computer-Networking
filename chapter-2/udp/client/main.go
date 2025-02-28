package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	serverName := "127.0.0.1"
	serverPort := "12000"
	serverAddress := net.JoinHostPort(serverName, serverPort)

	// 创建UDP地址
	udpAddr, err := net.ResolveUDPAddr("udp", serverAddress)
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	// 创建UDP连接
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("Error dialing UDP:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		// 读取用户输入
		fmt.Print("Input lowercase sentence: ")
		message, _ := reader.ReadString('\n')

		// 发送数据
		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error sending data:", err)
			return
		}

		// 接收数据
		buffer := make([]byte, 2048)
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			return
		}

		// 打印接收到的信息
		fmt.Println("Received message from", addr, ":", string(buffer[:n]))
	}
}
