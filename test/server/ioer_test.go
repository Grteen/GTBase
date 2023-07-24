package server

import (
	"GtBase/pkg/constants"
	"GtBase/src/server"
	"net"
	"syscall"
	"testing"
)

func TestEpoller(t *testing.T) {
	listenSock, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		t.Fatalf("Failed to create socket: %v", err)
	}
	defer syscall.Close(listenSock)

	addr := syscall.SockaddrInet4{Port: 1234}
	copy(addr.Addr[:], net.ParseIP("127.0.0.1").To4())

	err = syscall.Bind(listenSock, &addr)
	if err != nil {
		t.Fatalf("Failed to bind to address: %v", err)
	}

	syscall.Listen(listenSock, 0)

	poller := &server.EPoller{}
	poller.Run(listenSock)

	go func() {
		conn, err := net.Dial("tcp", "127.0.0.1:1234")
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		defer conn.Close()
	}()

	tasks, err := poller.Wait()
	if err != nil {
		t.Fatalf("Failed to wait for events: %v", err)
	}

	if len(tasks) != 1 {
		t.Fatalf("Expected 1 task, got %d", len(tasks))
	}

	if tasks[0].EventFd() != listenSock {
		t.Fatalf("Expected task Fd to be listenSock, got %d", tasks[0].EventFd())
	}

	if tasks[0].EventType() != constants.IoerAccept {
		t.Fatalf("Expected task EventType to be IoerAccept, got %d", tasks[0].EventType())
	}
}
