package server

import (
	"GtBase/pkg/constants"
	"GtBase/src/analyzer"
	"GtBase/src/client"
	"GtBase/src/nextwrite"
	"GtBase/src/option"
	"GtBase/src/page"
	"GtBase/src/replic"
	"context"
	"fmt"
	"net"
	"sync"
	"syscall"
)

type GtBaseServer struct {
	clients  map[int]*client.GtBaseClient
	clock    sync.Mutex
	ioer     Ioer
	listenFd int

	host string
	port int

	rs *replic.ReplicState
}

func (s *GtBaseServer) addClient(client *client.GtBaseClient) {
	s.clock.Lock()
	defer s.clock.Unlock()
	s.clients[client.GetFd()] = client
}

func (s *GtBaseServer) getClient(fd int) *client.GtBaseClient {
	s.clock.Lock()
	defer s.clock.Unlock()
	result, ok := s.clients[fd]
	if !ok {
		return nil
	}

	return result
}

func (s *GtBaseServer) Run() error {
	initFile()
	if option.NeedRedo() {
		errr := RedoLog()
		if errr != nil {
			return errr
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go page.FlushRedoDirtyList(ctx)

	if option.IsCache() {
		go page.FlushDirtyList(ctx)
	}

	listenFd, err := listenAndGetFd(s.port)
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
	addr, errg := client.GetAddressByFd(nfd)
	if errg != nil {
		return errg
	}

	s.addClient(client.CreateGtBaseClient(nfd, addr))
	return nil
}

func (s *GtBaseServer) handleCommand(client *client.GtBaseClient) error {
	bts, err := client.Read()
	if err != nil {
		if err.Error() == constants.ClientExitError {
			errr := s.ioer.Remove(client.GetFd())
			if errr != nil {
				return errr
			}
			return nil
		}
		return err
	}
	if bts == nil {
		return nil
	}

	cmn, errg := nextwrite.GetCMN()
	if errg != nil {
		return errg
	}

	args := analyzer.CreateCommandAssignArgs(client, s.rs, s.host, s.port, s.rs.GetUUID())
	result, msg := analyzer.CreateCommandAssign(bts, cmn, args).Assign().Analyze().Exec()
	if result != nil {
		errw := client.WriteResp(result.ToByte())
		if errw != nil {
			return errw
		}
	}

	if msg != nil {
		err := DoMessage(msg)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *GtBaseServer) assignTask(tasks []*Task) {
	for _, t := range tasks {
		if t.EventType() == constants.IoerAccept {
			err := s.handleAccept(s.listenFd)
			if err != nil {
				fmt.Println(err)
			}
		} else if t.EventType() == constants.IoerRead {
			err := s.handleCommand(s.getClient(int(t.EventFd())))
			if err != nil {
				fmt.Println(err)
			}
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

func CreateGtBaseServer(host string, port int) *GtBaseServer {
	return &GtBaseServer{ioer: &EPoller{}, clients: make(map[int]*client.GtBaseClient), rs: replic.CreateReplicState(), host: host, port: port}
}
