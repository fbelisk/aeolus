package example

import (
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
)
var host = flag.String("host", "localhost", "host")
var port = flag.String("port", "10101", "port")

func Run() {
	flag.Parse()
	conn, err := net.Dial("tcp", *host+":"+*port)
	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Println("Connecting to " + *host + ":" + *port)
	done := make(chan string)
	go handleWrite(conn, done)
	//go handleRead(conn, done)
	fmt.Println(<-done)
	//fmt.Println(<-done)
}
func handleWrite(conn net.Conn, done chan string) {
	for i := 110000; i > 0; i-- {
		msg := []byte("hello " + strconv.Itoa(i) + "\r\n")
		len := len(msg)
		lenByte := make([]byte, 4)
		binary.BigEndian.PutUint32(lenByte, uint32(len))
		_, e := conn.Write(append(lenByte))
		if e != nil {
			fmt.Println("Error to send message because of ", e.Error())
			break
		}
		_, e = conn.Write(append(msg))
		if e != nil {
			fmt.Println("Error to send message because of ", e.Error())
			break
		}
	}
	done <- "Sent Done"
}
func handleRead(conn net.Conn, done chan string) {
	buf := make([]byte, 1024)
	reqLen, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error to read message because of ", err)
		return
	}
	fmt.Println(string(buf[:reqLen-1]))
	done <- "Read"
}