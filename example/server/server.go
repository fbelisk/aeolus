package main

import (
	"fmt"
	sophonn "github.com/fbelisk/aeolus"
	"github.com/fbelisk/aeolus/util"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	w := &util.WaitGroupWrapper{}
	w.Wrap(func() {
		log.Println(http.ListenAndServe("localhost:9999", nil))
	})
	e, _ := sophonn.CreatePoller(func(inframe []byte) []byte {
		str := "recive success" + string(inframe)
		return []byte(str)
	})
	listener, err := net.Listen("tcp", ":10101")
	tcpListenr := listener.(*net.TCPListener)
	if err != nil {
		fmt.Println("HOLY SHIT IT DOES NOT RUN SUCCESS")
	}
	if err := e.Run(tcpListenr); err != nil {
		fmt.Println(err.Error())
	}
	w.Wait()
}

