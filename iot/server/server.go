package main

import (
	"net"
	"gobible/muduo-go/iot"
	"encoding/json"
	"fmt"
)

//与服务器相关的资源都放在这里面
type TcpServer struct {
	listener   *net.TCPListener
	hawkServer *net.TCPAddr
}

//服务器往客户端的数据包很简单地以\n换行结束了，偷了一个懒:)，正常情况下也可根据自己的协议来封装好
//然后在客户端写一个状态来处理
func processRecvData(packet *iot.Packet, conn net.Conn) {
	switch packet.PacketType {

	//是客户端给服务器端发送心态包！！
	case iot.HEART_BEAT_PACKET:
		var beatPacket iot.HeartPacket
		json.Unmarshal(packet.PacketContent, &beatPacket)
		fmt.Printf("recieve heat beat from [%s] ,data is [%v]\n", conn.RemoteAddr().String(), beatPacket)
		conn.Write([]byte("heartBeat received \n"))
	//上报数据包
	case iot.REPORT_PACKET:
		var reportPacket iot.ReportPacket
		json.Unmarshal(packet.PacketContent, &reportPacket)
		fmt.Printf("recieve report data from [%s] ,data is [%v]\n", conn.RemoteAddr().String(), reportPacket)
		conn.Write([]byte("Report data has received \n"))
	}
}

