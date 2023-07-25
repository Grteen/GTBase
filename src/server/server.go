package server

import (
	"GtBase/pkg/constants"
	"GtBase/src/analyzer"
	"GtBase/src/nextwrite"
	"net"
	"sync"
	"syscall"
)

type GtBaseServer struct {
	clients  map[int]*GtBaseClient
	clock    sync.Mutex
	ioer     Ioer
	listenFd int
}

func (s *GtBaseServer) addClient(client *GtBaseClient) {
	s.clock.Lock()
	defer s.clock.Unlock()
	s.clients[client.GetFd()] = client
}

func (s *GtBaseServer) getClient(fd int) *GtBaseClient {
	s.clock.Lock()
	defer s.clock.Unlock()
	result, ok := s.clients[fd]
	if !ok {
		return nil
	}

	return result
}

func (s *GtBaseServer) Run(port int) error {
	listenFd, err := listenAndGetFd(port)
	if err != nil {
		return err
	}

	s.listenFd = listenFd
	s.ioer.Run(listenFd)

	for {
		tasks, err := s.ioer.Wait()
		if err != nil {
			return err
		}
		s.assignTask(tasks)
	}
}

func (s *GtBaseServer) handleAccept(listenFd int) error {
	nfd, _, err := syscall.Accept(listenFd)
	if err != nil {
		return err
	}

	erra := s.ioer.AddRead(nfd)
	if erra != nil {
		return erra
	}
	s.addClient(CreateGtBaseClient(nfd))
	return nil
}

func (s *GtBaseServer) handleCommand(client *GtBaseClient) error {
	bts, err := client.Read()
	if err != nil {
		return err
	}

	cmn, errg := nextwrite.GetCMN()
	if errg != nil {
		return errg
	}

	result := analyzer.CreateCommandAssign(bts, cmn).Assign().Analyze().Exec().ToByte()

	errw := client.Write(result)
	if errw != nil {
		return errw
	}

	return nil
}

func (s *GtBaseServer) assignTask(tasks []*Task) {
	for _, t := range tasks {
		if t.EventType() == constants.IoerAccept {
			s.handleAccept(s.listenFd)
		} else if t.EventType() == constants.IoerRead {
			s.handleCommand(s.getClient(int(t.EventFd())))
		}
	}
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

	errl := syscall.Listen(listenSock, 0)
	if errl != nil {
		return -1, errl
	}

	return listenSock, nil
}

func CreateGtBaseServer() *GtBaseServer {
	return &GtBaseServer{ioer: &EPoller{}, clients: make(map[int]*GtBaseClient)}
}
