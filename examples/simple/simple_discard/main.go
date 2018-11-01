package main

import (
	"github.com/chenshuo/muduo-examples-in-go/examples/simple"
	"log"
)

func main() {
	discardServer := simple.NewDiscardServer(":2009")
	log.Println("discard on :2009")
	discardServer.Serve()
}
