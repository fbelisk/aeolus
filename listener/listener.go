package listener

import (
	"github.com/fbelisk/aeolus/connect"
	"net"
)

type Listener struct {
	Fd      int
	ln      *net.Listener
	clients []*connect.Client
	connections map[int]net.Conn
}

//func NewListener(protocol protocol.Protocal, addr string) (*Listener, error) {
//	listener, err := net.Listen(string(protocol), addr)
//	if err != nil {
//		return nil, err
//	}
//	l, ok := listener.(*net.TCPListener)
//	if !ok {
//		//todo error组件
//	}
//
//	file, err := l.File()
//	if err != nil {
//		return nil, err
//	}
//
//	if err != nil {
//		return nil, err
//	}
//	fd := int(file.Fd())
//	if err = unix.SetNonblock(fd, true); err != nil {
//		return nil, err
//	}
//	return &Listener{
//		Fd:      fd,
//		ln:      &listener,
//	}, nil
//}
//
//func readLen(r io.Reader, tmp []byte) (int32, error) {
//	_, err := io.ReadFull(r, tmp)
//	if err != nil {
//		return 0, err
//	}
//	return int32(binary.BigEndian.Uint32(tmp)), nil
//}
//
//func MessageHandle() ([]byte, error){
//	//获取数据段长度
//	bodyLen, err := readLen(client.Reader, client.lenSlice)
//	if err != nil {
//		return nil, err
//	}
//
//	//数据长度异常处理
//	if bodyLen <= 0 {
//		return nil, err
//	}
//
//	if int64(bodyLen) > p.ctx.nsqd.getOpts().MaxMsgSize {
//		return nil, err
//	}
//
//	//数据段读取
//	messageBody := make([]byte, bodyLen)
//	_, err = io.ReadFull(client.Reader, messageBody)
//	if err != nil {
//		return nil, err
//	}
//}
