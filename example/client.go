package example

import (
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/fbelisk/aeolus/util"
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
	w := &util.WaitGroupWrapper{}
	fmt.Println("Connecting to " + *host + ":" + *port)
	done := make(chan string)
	w.Wrap(func() {
		handleWrite(conn, done)
	})
	w.Wrap(func() {
		handleRead(conn, done)
	})
	w.Wait()
}

func handleWrite(conn net.Conn, done chan string) {
	for i := 1; i <= 10000000; i++ {
		msg := []byte("hello " + strconv.Itoa(i) + "\r\n")
		len := len(msg)
		lenByte := make([]byte, 4)
		binary.BigEndian.PutUint32(lenByte, uint32(len))
		_, e := conn.Write(lenByte)
		if e != nil {
			fmt.Println("Error to send message because of ", e.Error())
			break
		}
		_, e = conn.Write(msg)
		if e != nil {
			fmt.Println("Error to send message because of ", e.Error())
			break
		}
	}
}
func handleRead(conn net.Conn, done chan string) {
	for {
		buf := make([]byte, 1024)
		reqLen, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error to read message because of ", err)
			return
		}
		fmt.Println(string(buf[:reqLen-1]))
	}
}