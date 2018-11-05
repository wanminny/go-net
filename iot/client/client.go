package main

import (
	"net"
	"bufio"
	"fmt"
	"gobible/muduo-go/iot"
	"time"
	"encoding/json"
	"math/rand"
)

//客户端对象
type TcpClient struct {
	connection *net.TCPConn
	hawkServer *net.TCPAddr
	stopChan   chan struct{}
}

// 接收数据包
func (client *TcpClient) receivePackets() {
	reader := bufio.NewReader(client.connection)
	for {
		//承接上面说的服务器端的偷懒，我这里读也只是以\n为界限来读区分包
		msg, err := reader.ReadString('\n')
		if err != nil {
			//在这里也请处理如果服务器关闭时的异常
			close(client.stopChan)
			break
		}
		fmt.Print("recv: ",string(msg))
	}
}


//发送数据包
//仔细看代码其实这里做了两次json的序列化，有一次其实是不需要的
func (client *TcpClient) sendReportPacket() {
	reportPacket := iot.ReportPacket{
		Content:   iot.GetRandString(),
		Timestamp: time.Now().Unix(),
		Rand:      rand.Int(),
	}
	packetBytes, err := json.Marshal(reportPacket)
	if err != nil {
		fmt.Println(err.Error())
	}
	//这一次其实可以不需要，在封包的地方把类型和数据传进去即可
	packet := iot.Packet{
		PacketType:    iot.REPORT_PACKET,
		PacketContent: packetBytes,
	}
	sendBytes, err := json.Marshal(packet)
	if err != nil {
		fmt.Println(err.Error())
	}
	//发送
	client.connection.Write(EnPackSendData(sendBytes))
	fmt.Println("Send metric data success!")
}

//发送心跳包，与发送数据包一样
func (client *TcpClient) sendHeartPacket() {
	heartPacket := iot.HeartPacket{
		Version:   "1.0",
		Timestamp: time.Now().Unix(),
	}
	packetBytes, err := json.Marshal(heartPacket)
	if err != nil {
		fmt.Println(err.Error())
	}
	packet := iot.Packet{
		PacketType:    iot.HEART_BEAT_PACKET,
		PacketContent: packetBytes,
	}
	sendBytes, err := json.Marshal(packet)
	if err != nil {
		fmt.Println(err.Error())
	}
	client.connection.Write(EnPackSendData(sendBytes))
	fmt.Println("Send heartbeat data success!")
}
