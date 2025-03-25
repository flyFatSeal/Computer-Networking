package shared

import (
	"fmt"
	"net"
)

// 数据包结构
type Packet struct {
	SeqNumber int    // 序列号
	Data      string // 数据内容
}

func SendUDPPacket(conn *net.UDPConn, addr *net.UDPAddr, packet Packet) error {
	fmt.Printf("发送数据包: Seq=%d, Data=%s\n", packet.SeqNumber, packet.Data)
	data := append([]byte{byte(packet.SeqNumber)}, []byte(packet.Data)...)
	_, err := conn.WriteToUDP(data, addr)
	if err != nil {
		return fmt.Errorf("failed to send UDP packet: %v", err)
	}

	return nil
}

func SendUDPPacketConnected(conn *net.UDPConn, packet Packet) error {
	fmt.Printf("发送数据包: Seq=%d, Data=%s\n", packet.SeqNumber, packet.Data)
	data := append([]byte{byte(packet.SeqNumber)}, []byte(packet.Data)...)

	// 使用连接式的 Write 方法发送数据
	_, err := conn.Write(data)
	if err != nil {
		return fmt.Errorf("failed to send UDP packet: %v", err)
	}

	return nil
}
