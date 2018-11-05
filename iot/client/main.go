package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"math/rand"
	"net"
	"os"
	"time"
	"gobible/muduo-go/iot"
)

//客户端对象
type TcpClient struct {
	connection *net.TCPConn
	hawkServer *net.TCPAddr
	stopChan   chan struct{}
}


func main() {
	//拿到服务器地址信息
	hawkServer, err := net.ResolveTCPAddr("tcp", iot.Server)
	if err != nil {
		fmt.Printf("hawk server [%s] resolve error: [%s]", iot.Server, err.Error())
		os.Exit(1)
	}
	//连接服务器
	connection, err := net.DialTCP("tcp", nil, hawkServer)
	if err != nil {
		fmt.Printf("connect to hawk server error: [%s]", err.Error())
		os.Exit(1)
	}
	client := &TcpClient{
		connection: connection,
		hawkServer: hawkServer,
		stopChan:   make(chan struct{}),
	}
	//启动接收
	go client.receivePackets()

	//发送心跳的goroutine
	go func() {
		heartBeatTick := time.Tick(8 * time.Second)
		for {
			select {
			case <-heartBeatTick:
				client.sendHeartPacket()
			case <-client.stopChan:
				return
			}
		}
	}()

	//测试用的，开300个goroutine每秒发送一个包
	for i := 0; i < 1; i++ {
		go func() {
			sendTimer := time.After(14 * time.Second)
			for {
				select {
				case <-sendTimer:
					client.sendReportPacket()
				case <-client.stopChan:
					return
				}
			}
		}()
	}
	//等待退出
	<-client.stopChan
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


//使用的协议与服务器端保持一致

//数据包的格式为|0xFF|0xFF|len(高)|len(低)|Data|CRC高16位|0xFF|0xFE
//其中len为data的长度，实际长度为len(高)*256+len(低)
//CRC为32位CRC，取了最高16位共2Bytes
//0xFF|0xFF 和 0xFF|0xFE类似于前导码

func EnPackSendData(sendBytes []byte) []byte {

	// 8个字节分别是前4个字节  + 后4个字节
	packetLength := len(sendBytes) + 8
	result := make([]byte, packetLength)

	result[0] = 0xFF
	result[1] = 0xFF

	result[2] = byte(uint16(len(sendBytes)) >> 8)	//len（高）
	result[3] = byte(uint16(len(sendBytes)) & 0xFF) //len(低)

	//Data
	copy(result[4:], sendBytes)

	// 高16位
	sendCrc := crc32.ChecksumIEEE(sendBytes)
	result[packetLength-4] = byte(sendCrc >> 24)
	result[packetLength-3] = byte(sendCrc >> 16 & 0xFF)

	result[packetLength-2] = 0xFF
	result[packetLength-1] = 0xFE

	fmt.Println("EnPackSendData: ",string(result))
	return result
}