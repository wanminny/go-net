package main

import (
	"gobible/muduo-go/examples/simple"
)

func main() {
	timeServer := simple.NewTimeServer(":2037")
	timeServer.Serve()
}
