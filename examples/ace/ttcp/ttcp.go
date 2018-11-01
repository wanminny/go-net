package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"

	"gobible/muduo-go/muduo"
	"log"
)

type options struct {
	port     int
	length   int
	number   int

	transmit bool
	receive  bool

	nodelay  bool
	host     string
}

var opt options

type SessionMessage struct {
	Number, Length int32
}

func init() {
	flag.IntVar(&opt.port, "p", 5001, "TCP port")

	flag.IntVar(&opt.length, "l", 65536, "Buffer length")
	flag.IntVar(&opt.number, "n", 8192, "Number of buffers")

	flag.BoolVar(&opt.receive, "r", false, "Receive") //接受 作为服务器端
	// 主机
	flag.StringVar(&opt.host, "host", "", "Transmit")

	muduo.Check(binary.Size(SessionMessage{}) == 8, "packed struct")
}

//分析数据报文组成：
// msgHeader + msgPayLoad
// ------------------  msgHeader -------------------
//   4   number     |              4  length
// -------------------- payload ---------------------
//					|
//      length (4)  |   []byte   (length 个)
//---------------------------------------------------
//					|
//      length (4)  |   []byte   (length 个)
//---------------------------------------------------
//					|
//      length (4)  |   []byte   (length 个)
//---------------------------------------------------
//        			......      (number 个)
//---------------------------------------------------

// 客户端发送
func transmit() {

	sessionMessage := SessionMessage{int32(opt.number), int32(opt.length)}
	fmt.Printf("buffer length = %d\nnumber of buffers = %d\n",
		sessionMessage.Length, sessionMessage.Number)
	total_mb := float64(sessionMessage.Number) * float64(sessionMessage.Length) / 1024.0 / 1024.0
	fmt.Printf("%.3f MiB in total\n", total_mb)

	//客户端 拨号
	conn, err := net.Dial("tcp", net.JoinHostPort(opt.host, strconv.Itoa(opt.port)))
	muduo.PanicOnError(err)

	t := conn.(*net.TCPConn)
	t.SetNoDelay(false)

	defer conn.Close()

	start := time.Now()

	//发送啥？？ [先发送 sessionMessage 结构体！]
	err = binary.Write(conn, binary.BigEndian, &sessionMessage)
	muduo.PanicOnError(err)

	// 消息体 ？
	total_len := 4 + opt.length // binary.Size(int32(0)) == 4
	// println(total_len)

	payload := make([]byte, total_len)
	binary.BigEndian.PutUint32(payload, uint32(opt.length))
	for i := 0; i < opt.length; i++ {
		payload[4+i] = "0123456789ABCDEF"[i%16]
	}
	log.Println(string(payload))

	//多少个number个数
	for i := 0; i < opt.number; i++ {
		var n int
		n, err = conn.Write(payload)
		muduo.PanicOnError(err)
		muduo.Check(n == len(payload), "write payload")

		var ack int32 //因为发送的时候是int32
		err = binary.Read(conn, binary.BigEndian, &ack)
		muduo.PanicOnError(err)
		muduo.Check(ack == int32(opt.length), "ack")
	}

	elapsed := time.Since(start).Seconds()
	fmt.Printf("%.3f seconds\n%.3f MiB/s\n", elapsed, total_mb/elapsed)
}

//服务器端 接受
func receive() {

	//服务器端 监听 !
	listener := muduo.ListenTcpOrDie(fmt.Sprintf(":%d", opt.port))

	//两种资源 一个是 监听器！
	defer listener.Close()
	println("Accepting", listener.Addr().String())

	//一个是连接器 （另外一种资源！）
	conn, err := listener.Accept()
	muduo.PanicOnError(err)
	defer conn.Close()

	// Read header
	var sessionMessage SessionMessage
	//读取完整的一个 SessionMessage 结构体字节数; 【SessionMessage 结构就可以接受住了】
	err = binary.Read(conn, binary.BigEndian, &sessionMessage)
	muduo.PanicOnError(err)

	fmt.Printf("receive buffer length = %d\n receive number of buffers = %d\n",
		sessionMessage.Length, sessionMessage.Number)

	total_mb := float64(sessionMessage.Number) * float64(sessionMessage.Length) / 1024.0 / 1024.0
	fmt.Printf("%.3f MiB in total\n", total_mb)

	// 消息体:
	payload := make([]byte, sessionMessage.Length)
	start := time.Now()
	for i := 0; i < int(sessionMessage.Number); i++ {

		var length int32
		// 先读长度 ；长度只需要四个字节就可以接受住了；
		err = binary.Read(conn, binary.BigEndian, &length)
		muduo.PanicOnError(err)
		muduo.Check(length == sessionMessage.Length, "read length")

		var n int
		n, err = io.ReadFull(conn, payload)
		muduo.PanicOnError(err)
		muduo.Check(n == len(payload), "read payload")

		// ack             //每次发送完length个以后就来一个ack ;
		err = binary.Write(conn, binary.BigEndian, &length)
		muduo.PanicOnError(err)
	}

	elapsed := time.Since(start).Seconds()
	//总的传输量/总的实际 = 效率
	fmt.Printf("%.3f seconds\n%.3f MiB/s\n", elapsed, total_mb/elapsed)
}

//client ./ttcp -host  127.0.0.1
//server ./ttcp -r

func main() {
	flag.Parse()
	opt.transmit = opt.host != ""

	if opt.transmit == opt.receive {
		println("Either -r or -host must be specified.")
		return
	}

	if opt.transmit {
		transmit() //发送
	} else if opt.receive {
		receive() //接受
	} else {
		panic("unknown")
	}
}
