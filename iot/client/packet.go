package main

import (
	"hash/crc32"
	"fmt"
)

//使用的协议与服务器端保持一致

//数据包的格式为|0xFF|0xFF|len(高)|len(低)|Data|CRC高16位|0xFF|0xFE
//其中len为data的长度，实际长度为len(高)*256+len(低)  // 其中256 = 2^8 次方
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
