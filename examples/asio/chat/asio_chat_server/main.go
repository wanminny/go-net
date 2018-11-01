package main

import (
	"log"

	"gobible/muduo-go/examples/asio/chat"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	server := chat.NewChatServer(":3399")
	log.Println("chat server at :3399")
	server.Run()
}
