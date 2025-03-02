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

	// 创建TCP连接
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	// 读取用户输入
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Input lowercase sentence: ")
	sentence, _ := reader.ReadString('\n')

	// 发送数据
	_, err = conn.Write([]byte(sentence))
	if err != nil {
		fmt.Println("Error sending data:", err)
		return
	}

	// 接收数据
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error receiving data:", err)
		return
	}

	// 打印接收到的信息
	fmt.Printf("From Server: %s\n", string(buffer[:n]))
}
