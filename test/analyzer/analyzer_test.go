package analyzer

import (
	"GtBase/pkg/constants"
	"GtBase/src/analyzer"
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

func TestAnalyzer(t *testing.T) {
	page.DeleteBucketPageFile()
	page.InitBucketPageFile()
	data := []struct {
		key string
		val string
	}{
		{"Key", "Val"},
		{"Hello", "World"},
	}

	for _, d := range data {
		cmd := make([][]byte, 0)
		cmd = append(cmd, []byte(d.key))
		cmd = append(cmd, []byte(d.val))

		a := analyzer.CreateSetAnalyzer(cmd, []byte(""), -1, nil)
		res := a.Analyze().Exec().ToString()
		if res != constants.ServerOkReturn {
			t.Errorf("Exec should get %v but got %v", constants.ServerOkReturn, res)
		}
	}

	for _, d := range data {
		cmd := make([][]byte, 0)
		cmd = append(cmd, []byte(d.key))

		a := analyzer.CreateGetAnalyzer(cmd, []byte(""), -1, nil)
		res := a.Analyze().Exec().ToString()
		if res != d.val {
			t.Errorf("Exec should get %v but got %v", d.val, res)
		}
	}

	cmd := make([][]byte, 0)
	cmd = append(cmd, []byte(data[1].key))

	a := analyzer.CreateDelAnalyzer(cmd, []byte(""), -1, nil)
	res := a.Analyze().Exec().ToString()
	if res != constants.ServerOkReturn {
		t.Errorf("Exec should get %v but got %v", constants.ServerOkReturn, res)
	}

	a2 := analyzer.CreateGetAnalyzer(cmd, []byte(""), -1, nil)
	res = a2.Analyze().Exec().ToString()
	if res != constants.ServerGetNilReturn {
		t.Errorf("Exec should get %v but got %v", constants.ServerGetNilReturn, res)
	}
}

func TestSlave(t *testing.T) {
	listenSock, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		t.Fatalf("Failed to create socket: %v", err)
	}
	defer syscall.Close(listenSock)

	addr := syscall.SockaddrInet4{Port: 7866}
	copy(addr.Addr[:], net.ParseIP("127.0.0.1").To4())

	err = syscall.Bind(listenSock, &addr)
	if err != nil {
		t.Fatalf("Failed to bind to address: %v", err)
	}

	syscall.Listen(listenSock, 0)

	poller := &server.EPoller{}
	poller.Run(listenSock)

	result := make([]byte, 0)

	ch := make(chan []byte)
	go func(result []byte) {
		conn, err := net.Dial("tcp", "127.0.0.1:7866")
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

			result = append(result, buf[0:n]...)
			if utils.EqualByteSlice(result[len(result)-2:], []byte(constants.ReplicRedoLogEnd)) {
				result = result[:len(result)-2]
				ch <- result[len(constants.RedoCommand) : len(constants.RedoCommand)+int(constants.SendRedoLogSeqSize)]
				result = result[len(constants.RedoCommand)+int(constants.SendRedoLogSeqSize):]
				break
			}
		}
		ch <- result
	}(result)

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
	c := client.CreateGtBaseClient(nfd, client.CreateAddress("127.0.0.1", 0))

	cmd := make([][]byte, 0)
	cmd = append(cmd, utils.Encodeint32ToBytesSmallEnd(0))
	cmd = append(cmd, utils.Encodeint32ToBytesSmallEnd(0))
	cmd = append(cmd, utils.Encodeint32ToBytesSmallEnd(1))

	args := analyzer.CreateCommandAssignArgs(c, rs)

	analyzer.CreateSlaveAnalyzer(cmd, nil, 0, args).Analyze().Exec()

	seqbts := <-ch
	seq := utils.EncodeBytesSmallEndToint32(seqbts)
	if seq != 1 {
		t.Errorf("seq should be %v but got %v", 1, seq)
	}
	result = <-ch

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

	logIdx := int32(len(result)) / int32(constants.PageSize)
	logOff := int32(len(result)) % int32(constants.PageSize)

	cmd2 := make([][]byte, 0)
	cmd2 = append(cmd2, utils.Encodeint32ToBytesSmallEnd(logIdx))
	cmd2 = append(cmd2, utils.Encodeint32ToBytesSmallEnd(logOff))
	cmd2 = append(cmd2, utils.Encodeint32ToBytesSmallEnd(2))

	analyzer.CreateGetRedoAnalyzer(cmd2, nil, 0, args).Analyze().Exec()
	s, ok := rs.GetSlave(c.GenerateKey())
	if !ok {
		t.Errorf("Should be Ok")
	}

	if idx, off := s.GetLogInfo(); idx != logIdx || off != logOff {
		t.Errorf("GetLogInfo should get %v idx %v off but got %v idx %v off", logIdx, logOff, idx, off)
	}
}
