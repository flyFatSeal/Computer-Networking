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

func SendUDPPacket(conn *net.UDPConn, addr *net.UDPAddr, packet Packet) {
	fmt.Printf("发送数据包: Seq=%d, Data=%s\n", packet.SeqNumber, packet.Data)
	data := append([]byte{byte(packet.SeqNumber)}, []byte(packet.Data)...)
	conn.WriteToUDP(data, addr)
}
