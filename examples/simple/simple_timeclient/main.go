package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"time"

	"log"
	"gobible/muduo-go/muduo"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s host\n", os.Args[0])
		return
	}

	host := os.Args[1]
	conn, err := net.Dial("tcp", net.JoinHostPort(host, "2037"))
	muduo.PanicOnError(err)
	defer conn.Close()

	var unixtime int64
	err = binary.Read(conn, binary.BigEndian, &unixtime)
	muduo.PanicOnError(err)

	log.Println(unixtime)

	//自己解析的实际上收到的数据是上面的！！
	println(time.Unix((unixtime), 0).String())
}
