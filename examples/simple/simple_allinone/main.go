package main

import (
	"gobible/muduo-go/examples/simple"
)

func main() {
	ch := make(chan bool)

	// long connection

	chargenServer := simple.NewChargenServer(":2019")
	go chargenServer.Serve()

	discardServer := simple.NewDiscardServer(":2009")
	go discardServer.Serve()

	echoServer := simple.NewEchoServer(":2007")
	go echoServer.Serve()

	// short connection
	// [用完后直接关闭了！]
	daytimeServer := simple.NewDaytimeServer(":2013")
	go daytimeServer.Serve()

	timeServer := simple.NewTimeServer(":2037")
	go timeServer.Serve()

	<-ch // wait forever
}
