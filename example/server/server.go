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
		log.Println(http.ListenAndServe("0.0.0.0:9999", nil))
	})
	w.Wrap(func() {
		listener, err := net.Listen("tcp", ":10101")
		tcpListenr := listener.(*net.TCPListener)
		if err != nil {
			fmt.Println("HOLY SHIT IT DOES NOT RUN SUCCESS")
		}
		e, _ := sophonn.CreatePoller(tcpListenr, func(inframe []byte) []byte {
			str := "recive success" + string(inframe)
			return []byte(str)
		})
		if err := e.Run(); err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("异常退出")
	})
	w.Wait()
}

