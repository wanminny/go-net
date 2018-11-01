package main

import (
	"gobible/muduo-go/examples/simple"
	"log"
)

func main() {
	daytimeServer := simple.NewDaytimeServer(":2013")
	log.Println("day time server :2013")
	daytimeServer.Serve()
}
