package simple

import (
	"fmt"
	"log"
	"net"
	"time"

	"gobible/muduo-go/muduo"
	"encoding/binary"
)

type DaytimeServer struct {
	listener net.Listener
}

func NewDaytimeServer(listenAddr string) *DaytimeServer {
	server := new(DaytimeServer)
	server.listener = muduo.ListenTcpOrDie(listenAddr)
	return server
}

func (s *DaytimeServer) Serve() {

	defer s.listener.Close()
	for {
		conn, err := s.listener.Accept()
		if err == nil {
			printConn(conn, "daytime", "UP")
			str := fmt.Sprintf("%v\n", time.Now())

			log.Println(str,len(str))
			//var lenStr int32 = int32(len(str))

			// 发送 方式1
			//conn.Write([]byte(str))

			 //发送 方式 2
			//binary.Write(conn,binary.BigEndian,lenStr)
			binary.Write(conn,binary.BigEndian,str)

			printConn(conn, "daytime", "DOWN")
			conn.Close()
		} else {
			log.Println("daytime:", err.Error())
			// TODO: break if ! temporary
		}
	}
}
