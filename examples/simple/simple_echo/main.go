package main

import (
	"gobible/muduo-go/examples/simple"

)


func main() {
	echoServer := simple.NewEchoServer(":2019")
	echoServer.Serve()
}
