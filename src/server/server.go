package server

import (
	"GtBase/src/object"
	"net"
)

type Command interface {
	Exec() object.Object
}

type GtBaseServer struct {
	clients  map[int]*GtBaseClient
	listener *net.Listener
}
