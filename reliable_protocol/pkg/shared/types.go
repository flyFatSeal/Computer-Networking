package shared

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"time"
)

const lossProb = 0.05       // 每 100 个丢失 5 个
const corruptionProb = 0.02 // 每 100 个损坏 2 个

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

	packetToSend := SendToMedium(packet, lossProb, corruptionProb)

	if packetToSend == nil {
		// 数据包丢失
		return fmt.Errorf("packet lost: SeqNum=%d", packet.SeqNum)
	}

	packetToSend.Checksum = CalculateChecksum(*packetToSend)

	data, err := json.Marshal(packetToSend)
	if err != nil {
		return fmt.Errorf("failed to marshal packet: %v", err)
	}

	// 打印发送的内容
	fmt.Printf("发送数据包: Seq=%d, Ack=%d, Data=%s, Checksum=%d\n",
		packetToSend.SeqNum, packetToSend.AckNum, packetToSend.Data, packetToSend.Checksum)
	_, _error := conn.WriteToUDP(data, addr)
	if _error != nil {
		return fmt.Errorf("failed to send UDP packet: %v", _error)
	}

	return nil
}

func SendUDPPacketConnected(conn *net.UDPConn, packet Packet) error {

	packetToSend := SendToMedium(packet, lossProb, corruptionProb)

	if packetToSend == nil {
		// 数据包丢失
		return fmt.Errorf("packet lost: SeqNum=%d", packet.SeqNum)
	}
	packetToSend.Checksum = CalculateChecksum(*packetToSend)

	// 将 Packet 序列化为 JSON
	data, err := json.Marshal(packetToSend)
	if err != nil {
		return fmt.Errorf("failed to marshal packet: %v", err)
	}

	// 打印发送的内容
	fmt.Printf("发送数据包: Seq=%d, Ack=%d, Data=%s, Checksum=%d\n",
		packetToSend.SeqNum, packetToSend.AckNum, packetToSend.Data, packetToSend.Checksum)
	// 使用连接式的 Write 方法发送数据
	_, errWrite := conn.Write(data)
	if errWrite != nil {
		return fmt.Errorf("failed to send UDP packet: %v", errWrite)
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
