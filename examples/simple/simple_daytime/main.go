package main

import (
	"github.com/chenshuo/muduo-examples-in-go/examples/simple"
	"log"
)

func main() {
	daytimeServer := simple.NewDaytimeServer(":2013")
	log.Println("day time server :2013")
	daytimeServer.Serve()
}
