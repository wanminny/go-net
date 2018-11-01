package main

import (
	"log"
	"gobible/muduo-go/examples/simple"

)

func main() {
	discardServer := simple.NewDiscardServer(":2009")
	log.Println("discard on :2009")
	discardServer.Serve()
}
