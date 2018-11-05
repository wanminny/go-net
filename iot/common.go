package iot

import (
	"fmt"
	"os"
	"math/rand"
	"log"
)

//数据包类型
const (
	HEART_BEAT_PACKET = 0x00
	REPORT_PACKET     = 0x01
)

//默认的服务器地址
var (
	Server = "127.0.0.1:8080"
)

//数据包
type Packet struct {
	PacketType    byte
	PacketContent []byte
}

//心跳包
type HeartPacket struct {
	Version   string `json:"version"`
	Timestamp int64  `json:"timestamp"`
}

//数据包
type ReportPacket struct {
	Content   string `json:"content"`
	Rand      int    `json:"rand"`
	Timestamp int64  `json:"timestamp"`
}


//处理错误，根据实际情况选择这样处理，还是在函数调之后不同的地方不同处理
func CheckErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

//拿一串随机字符
func GetRandString() string {
	length := rand.Intn(50)
	strBytes := make([]byte, length)
	for i := 0; i < length; i++ {
		strBytes[i] = byte(rand.Intn(26) + 97)
	}
	log.Println("GetRandString :",string(strBytes))
	return string(strBytes)
}
