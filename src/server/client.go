package server

import "net"

type GtBaseClient struct {
	conn *net.Conn
}
