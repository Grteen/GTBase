package server

import (
	"net"
	"syscall"
)

type GtBaseServer struct {
	clients map[int]*GtBaseClient
	ioer    Ioer
}

func (s *GtBaseServer) addClient(client *GtBaseClient) {
	s.clients[client.GetFd()] = client
}

func (s *GtBaseServer) Run(port int) error {
	listenFd, err := listenAndGetFd(port)
	if err != nil {
		return err
	}

	s.ioer.Run(listenFd)

	return nil
}

func listenAndGetFd(port int) (int, error) {
	listenSock, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return -1, err
	}

	addr := syscall.SockaddrInet4{Port: port}
	copy(addr.Addr[:], net.ParseIP("127.0.0.1").To4())

	errb := syscall.Bind(listenSock, &addr)
	if err != nil {
		return -1, errb
	}

	syscall.Listen(listenSock, 0)

	return listenSock, nil
}

func (s *GtBaseServer) handleAccept(listenFd int) error {
	nfd, _, err := syscall.Accept(listenFd)
	if err != nil {
		return err
	}

	s.addClient(CreateGtBaseClient(nfd))
	return nil
}

// func (s *GtBaseClient) handleCommand(client *GtBaseClient) error {
// 	bts, err := client.Read()
// 	if err != nil {
// 		return err
// 	}

// }
