package command

import (
	"GtBase/pkg/constants"
	"GtBase/src/client"
	"GtBase/src/command"
	"GtBase/src/page"
	"GtBase/src/replic"
	"GtBase/src/server"
	"GtBase/utils"
	"fmt"
	"net"
	"syscall"
	"testing"
)

func TestSlaveCommnd(t *testing.T) {
	listenSock, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		t.Fatalf("Failed to create socket: %v", err)
	}
	defer syscall.Close(listenSock)

	addr := syscall.SockaddrInet4{Port: 8888}
	copy(addr.Addr[:], net.ParseIP("127.0.0.1").To4())

	err = syscall.Bind(listenSock, &addr)
	if err != nil {
		t.Fatalf("Failed to bind to address: %v", err)
	}

	syscall.Listen(listenSock, 0)

	poller := &server.EPoller{}
	poller.Run(listenSock)

	ch := make(chan [][]byte)
	go func() {
		listner, err := net.Listen("tcp", "127.0.0.1:9677")
		if err != nil {
			t.Errorf(err.Error())
		}
		conn, err := net.Dial("tcp", "127.0.0.1:8888")
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		con, err := listner.Accept()
		if err != nil {
			t.Errorf(err.Error())
		}
		defer conn.Close()
		result := make([]byte, 0)
		for {
			buf := make([]byte, 1024)
			n, err := con.Read(buf)
			if err != nil {
				t.Errorf(err.Error())
			}

			result = append(result, buf[0:n]...)
			if utils.EqualByteSlice(result[len(result)-2:], []byte(constants.ReplicRedoLogEnd)) {
				fields := utils.DecodeGtBasePacket(result)
				ch <- fields
				break
			}
		}
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

	nfd, _, err := syscall.Accept(tasks[0].EventFd())
	if err != nil {
		t.Errorf(err.Error())
	}

	rs := replic.CreateReplicState()
	c := client.CreateGtBaseClient(nfd, client.CreateAddress("127.0.0.1", 9677))
	command.Slave(0, 0, 1, "127.0.0.1", 9677, "0", c, rs)

	fields := <-ch
	seq := utils.EncodeBytesSmallEndToint32(fields[1])
	if seq != 1 {
		t.Errorf("seq should be %v but got %v", 1, seq)
	}

	res := make([]byte, 0)
	for i := 0; i < 10; i++ {
		pg, err := page.ReadRedoPage(int32(i))
		if err != nil {
			t.Errorf(err.Error())
		}

		if pg.Src()[0] == byte(0) {
			break
		}

		res = append(res, pg.Src()...)
	}

	result := fields[2]

	if !utils.EqualByteSlice(res, result) {
		t.Errorf("ReadRedoPage and SendRedoLog not same")
		fmt.Println(len(res), len(result))
	}

	logIdx := int32(len(result)) / int32(constants.PageSize)
	logOff := int32(len(result)) % int32(constants.PageSize)
	command.GetRedo(logIdx, logOff, 2, "0", c, rs)

	s, ok := rs.GetSlave("0")
	if !ok {
		t.Errorf("Should be Ok")
	}

	if idx, off := s.GetLogInfo(); idx != logIdx || off != logOff {
		t.Errorf("GetLogInfo should get %v idx %v off but got %v idx %v off", logIdx, logOff, idx, off)
	}
}
