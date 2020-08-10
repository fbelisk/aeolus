package connect

import "net"

//todo 是否拆分协议版本
type Client struct {
	Id int
	conn net.Listener
}