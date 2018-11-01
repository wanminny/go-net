package main

import (
	"gobible/muduo-go/examples/simple"

)

func main() {
	//CHARGEN是在TCP连接建立后，服务器不断传送任意的字符到客户端，直到客户端关闭连接
	// dos 攻击
	chargenServer := simple.NewChargenServer(":2019")
	chargenServer.ServeWithMeter()
}
