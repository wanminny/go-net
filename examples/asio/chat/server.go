package chat

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
	"runtime"
	"time"

	"gobible/muduo-go/muduo"
)

type ChatServer struct {
	listener   net.Listener  //每个server都有一个监听listener

	conns      map[*connection]bool //所有连接
	register   chan *connection  // 注册
	unregister chan *connection  //注销

	broadcast  chan []byte  //要广播的消息
}

// 新建 server 【全部初始化 map chan etc.】
func NewChatServer(listenAddr string) *ChatServer {
	return &ChatServer{
		listener:   muduo.ListenTcpOrDie(listenAddr),
		conns:      make(map[*connection]bool),
		broadcast:  make(chan []byte),      // size?
		register:   make(chan *connection), // size?
		unregister: make(chan *connection), // size?
	}
}

type connection struct {
	conn net.Conn
	// FIXME: use bufio to save syscall
	send chan []byte
}

func (c *connection) input(broadcast chan []byte) {
	for {
		message, err := c.readMessage()
		if err != nil {
			log.Println(err)
			break
		}
		broadcast <- message
	}
}

func (c *connection) output() {
	defer c.close()
	for m := range c.send {

		//先发送长度 【1】
		err := binary.Write(c.conn, binary.BigEndian, int32(len(m)))
		if err != nil {
			log.Println(err)
			break
		}
		var n int
		//后发生内容 【2】
		//TODO 【1】【2】区别？

		n, err = c.conn.Write(m)
		if err != nil {
			log.Println(err)
			break
		}
		// 判断发生的内容是否正确
		if n != len(m) {
			log.Println("short write")
			break
		}
	}
}

func (c *connection) close() {
	log.Println("close connection")
	c.conn.Close()
}

func (c *connection) readMessage() (message []byte, err error) {
	var length int32
	err = binary.Read(c.conn, binary.BigEndian, &length)
	if err != nil {
		return nil, err
	}
	if length > 65536 || length < 0 {
		return nil, errors.New("invalid length")
	}
	message = make([]byte, int(length))
	if length > 0 {
		var n int
		n, err = io.ReadFull(c.conn, message)
		if err != nil {
			return nil, err
		}
		if n != len(message) {
			return message, errors.New("short read")
		}
	}
	return message, nil
}

// 连接 协程 【区别于主协程】
func (s *ChatServer) ServeConn(conn net.Conn) {
	c := &connection{conn: conn, send: make(chan []byte, 1024)}
	s.register <- c
	defer func() { s.unregister <- c }()

	go c.output()
	c.input(s.broadcast)
}

func (s *ChatServer) Run() {
	ticks := time.Tick(time.Second * 1)
	go muduo.ServeTcp(s.listener, s, "chat")
	for {
		select {

		case c := <-s.register: //注册
			//所有连接
			s.conns[c] = true

		case c := <-s.unregister:  // 注销
			//
			delete(s.conns, c)
			close(c.send)

		case m := <-s.broadcast:  //广播
			for c := range s.conns {
				select {
				case c.send <- m:
				default:
					delete(s.conns, c)
					close(c.send)
					log.Println("kick slow connection")
				}
			}
		case _ = <-ticks:
			log.Println(len(s.conns), runtime.NumGoroutine())
		}
	}
}
