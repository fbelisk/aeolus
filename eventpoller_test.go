package sophonn

import (
	"fmt"
	"github.com/fbelisk/aeolus/poll"
	"net"
	"testing"
)

func TestEventpoller_Run(t *testing.T) {
	type fields struct {
		p          *poll.Poller
		clients    map[int]*Conn
		React      BusinessHandler
		ReadBuffer []byte
		Compressor Compressor
		Codec      Codec
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, _ := CreatePoller(func(inframe []byte) []byte {
				str := "recive success" + string(inframe)
				return []byte(str)
			})
			listener, err := net.Listen("tcp", ":10101")
			tcpListenr := listener.(*net.TCPListener)
			if err != nil {
				fmt.Println("HOLY SHIT IT DOES NOT RUN SUCCESS")
			}
			if err := e.Run(tcpListenr); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}