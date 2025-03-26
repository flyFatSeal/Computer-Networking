package shared

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

// 数据包结构
type Packet struct {
	SeqNum   int    // 序列号
	AckNum   int    // ACK号
	Data     string // 数据
	Checksum int    // 校验和
}

// CalculateChecksum 计算校验和
func CalculateChecksum(packet Packet) int {
	checksum := packet.SeqNum + packet.AckNum
	for _, char := range packet.Data {
		checksum += int(char)
	}
	return checksum
}

// IsCorrupted 检查数据包是否损坏
func IsCorrupted(packet Packet) bool {
	return CalculateChecksum(packet) != packet.Checksum
}

func SendUDPPacket(conn *net.UDPConn, addr *net.UDPAddr, packet Packet) error {
	fmt.Printf("发送数据包: Seq=%d, Data=%s\n", packet.SeqNum, packet.Data)
	data := append([]byte{byte(packet.SeqNum)}, []byte(packet.Data)...)
	_, err := conn.WriteToUDP(data, addr)
	if err != nil {
		return fmt.Errorf("failed to send UDP packet: %v", err)
	}

	return nil
}

func SendUDPPacketConnected(conn *net.UDPConn, packet Packet) error {
	fmt.Printf("发送数据包: Seq=%d, Data=%s\n", packet.SeqNum, packet.Data)
	data := append([]byte{byte(packet.SeqNum)}, []byte(packet.Data)...)

	// 使用连接式的 Write 方法发送数据
	_, err := conn.Write(data)
	if err != nil {
		return fmt.Errorf("failed to send UDP packet: %v", err)
	}

	return nil
}

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// SendToMedium 模拟发送数据包到网络
func SendToMedium(packet Packet, lossProb, corruptionProb float64) *Packet {
	// 模拟丢包
	if rng.Float64() < lossProb {
		fmt.Printf("Packet lost: SeqNum=%d\n", packet.SeqNum)
		return nil
	}

	// 模拟损坏
	if rng.Float64() < corruptionProb {
		packet.Checksum = -1 // 损坏校验和
		fmt.Printf("Packet corrupted: SeqNum=%d\n", packet.SeqNum)
	}

	return &packet
}
