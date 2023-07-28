package replic

import (
	"GtBase/pkg/constants"
	"GtBase/src/client"
	"GtBase/src/page"
	"GtBase/src/replic"
	"GtBase/src/server"
	"GtBase/utils"
	"fmt"
	"net"
	"syscall"
	"testing"
)

func TestHeart(t *testing.T) {
	listenSock, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		t.Fatalf("Failed to create socket: %v", err)
	}
	defer syscall.Close(listenSock)

	addr := syscall.SockaddrInet4{Port: 7663}
	copy(addr.Addr[:], net.ParseIP("127.0.0.1").To4())

	err = syscall.Bind(listenSock, &addr)
	if err != nil {
		t.Fatalf("Failed to bind to address: %v", err)
	}

	syscall.Listen(listenSock, 0)

	poller := &server.EPoller{}
	poller.Run(listenSock)

	ch := make(chan []byte)
	h := make(chan struct{})
	go func() {
		result := make([]byte, 0)
		conn, err := net.Dial("tcp", "127.0.0.1:7663")
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		defer conn.Close()
		for {
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				t.Errorf(err.Error())
			}
			if !utils.EqualByteSlice(buf[:n], []byte(constants.HeartCommand)) {
				result = append(result, buf[0:n]...)
				if utils.EqualByteSlice(result[len(result)-2:], []byte(constants.ReplicRedoLogEnd)) {
					result = result[:len(result)-2]
					ch <- result[len(constants.RedoCommand)+1 : len(constants.RedoCommand)+1+int(constants.SendRedoLogSeqSize)]
					result = result[len(constants.RedoCommand)+int(constants.SendRedoLogSeqSize)+2:]
					break
				}
			} else {
				h <- struct{}{}
			}
		}
		ch <- result
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

	s := replic.CreateSlave(0, 0, 1, client.CreateGtBaseClient(nfd, client.CreateAddress("127.0.0.1", 0)))
	s.SendHeartToSlave()
	<-h
	s.SetSyncStateLock(constants.SlaveSync)
	s.GetHeartRespFromSlave(0, 0, 2)

	seqbts := <-ch
	seq := utils.EncodeBytesSmallEndToint32(seqbts)
	if seq != 2 {
		t.Errorf("seq should be %v but got %v", 2, seq)
	}
	result := <-ch

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

	if !utils.EqualByteSlice(res, result) {
		t.Errorf("ReadRedoPage and SendRedoLog not same")
		fmt.Println(len(res), len(result))
	}
}
