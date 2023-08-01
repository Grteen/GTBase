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
	"strconv"
	"syscall"
	"testing"
	"time"
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

	addr := syscall.SockaddrInet4{Port: 9855}
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
		parts := make([]byte, 0)
		listener, err := net.Listen("tcp", "127.0.0.1:4123")
		if err != nil {
			t.Errorf(err.Error())
		}

		conn, err := net.Dial("tcp", "127.0.0.1:9855")
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		defer conn.Close()

		con, err := listener.Accept()
		if err != nil {
			t.Errorf(err.Error())
		}
		for {
			buf := make([]byte, 1024)
			n, err := con.Read(buf)
			if err != nil {
				t.Errorf(err.Error())
			}

			parts = append(parts, buf[0:n]...)
			if utils.EqualByteSlice(parts[len(parts)-2:], []byte(constants.ReplicRedoLogEnd)) {
				fileds := utils.DecodeGtBasePacket(parts[:len(parts)-2])
				ch <- fileds
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
	c := client.CreateGtBaseClient(nfd, client.CreateAddress("127.0.0.1", 4123))

	cmd := make([][]byte, 0)
	cmd = append(cmd, utils.Encodeint32ToBytesSmallEnd(0))
	cmd = append(cmd, utils.Encodeint32ToBytesSmallEnd(0))
	cmd = append(cmd, utils.Encodeint32ToBytesSmallEnd(1))
	cmd = append(cmd, []byte("127.0.0.1"))
	cmd = append(cmd, utils.Encodeint32ToBytesSmallEnd(4123))

	args := analyzer.CreateCommandAssignArgs(c, rs, "127.0.0.1", 9855)

	analyzer.CreateSlaveAnalyzer(cmd, nil, 0, args).Analyze().Exec()

	fileds := <-ch
	seq := utils.EncodeBytesSmallEndToint32(fileds[1])
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

	if !utils.EqualByteSlice(res, fileds[2]) {
		t.Errorf("ReadRedoPage and SendRedoLog not same")
		fmt.Println(len(res), len(fileds[2]))
	}

	logIdx := int32(len(fileds[2])) / int32(constants.PageSize)
	logOff := int32(len(fileds[2])) % int32(constants.PageSize)

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

func TestHeart(t *testing.T) {
	listenSock, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		t.Fatalf("Failed to create socket: %v", err)
	}
	defer syscall.Close(listenSock)

	addr := syscall.SockaddrInet4{Port: 7841}
	copy(addr.Addr[:], net.ParseIP("127.0.0.1").To4())

	err = syscall.Bind(listenSock, &addr)
	if err != nil {
		t.Fatalf("Failed to bind to address: %v", err)
	}

	syscall.Listen(listenSock, 0)

	poller := &server.EPoller{}
	poller.Run(listenSock)

	ch := make(chan [][]byte)
	h := make(chan struct{})
	go func() {
		listener, err := net.Listen("tcp", "127.0.0.1:5146")
		if err != nil {
			t.Errorf(err.Error())
		}

		conn, err := net.Dial("tcp", "127.0.0.1:7841")
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		defer conn.Close()
		parts := make([]byte, 0)
		con, err := listener.Accept()
		if err != nil {
			t.Errorf(err.Error())
		}
		for {
			buf := make([]byte, 1024)
			n, err := con.Read(buf)
			if err != nil {
				t.Errorf(err.Error())
			}
			temp := buf[0:n]
			parts = append(parts, temp...)
			if utils.EqualByteSlice(parts[len(parts)-2:], []byte(constants.CommandSep)) {
				heartFileds := utils.DecodeGtBasePacket(parts[:len(parts)-2])

				if !utils.EqualByteSlice(heartFileds[0], []byte(constants.HeartCommand)) {
					if utils.EqualByteSlice(parts[len(parts)-2:], []byte(constants.ReplicRedoLogEnd)) {
						fileds := utils.DecodeGtBasePacket(parts[:len(parts)-2])
						ch <- fileds
						break
					}
				} else {
					h <- struct{}{}
				}
				parts = make([]byte, 0)
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

	cmd2 := make([][]byte, 0)
	cmd2 = append(cmd2, utils.Encodeint32ToBytesSmallEnd(0))
	cmd2 = append(cmd2, utils.Encodeint32ToBytesSmallEnd(0))
	cmd2 = append(cmd2, utils.Encodeint32ToBytesSmallEnd(0))
	cmd2 = append(cmd2, utils.Encodeint32ToBytesSmallEnd(2))

	c := client.CreateGtBaseClient(nfd, client.CreateAddress("127.0.0.1", 5146))
	s := replic.CreateSlave(0, 0, 1, c)
	s.InitClient("127.0.0.1", 5146)
	rs := replic.CreateReplicState()
	rs.AppendSlaveLock(s)

	s.SendHeartToSlave()
	<-h
	s.SetSyncStateLock(constants.SlaveSync)
	args := analyzer.CreateCommandAssignArgs(c, rs, "127.0.0.1", 7841)

	analyzer.CreateGetHeartAnalyzer(cmd2, nil, -1, args).Analyze().Exec()

	fileds := <-ch
	seq := utils.EncodeBytesSmallEndToint32(fileds[1])
	if seq != 2 {
		t.Errorf("seq should be %v but got %v", 2, seq)
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

	if !utils.EqualByteSlice(res, fileds[2]) {
		t.Errorf("ReadRedoPage and SendRedoLog not same")
		fmt.Println(len(res), len(fileds[2]))
	}
}

func TestHeartAnalyzer(t *testing.T) {
	port := 4244
	lfd, err := utils.BindAndListen(port)
	if err != nil {
		t.Errorf(err.Error())
	}

	go func() {
		fd, err := utils.LocalDial(port)
		if err != nil {
			t.Errorf(err.Error())
		}
		defer utils.CloseFd(fd)

		parts, errr := utils.ReadFd(fd)
		if errr != nil {
			t.Errorf(errr.Error())
		}

		fileds := make([][]byte, 0)
		fileds = append(fileds, []byte(constants.HeartCommand))
		fileds = append(fileds, utils.Encodeint32ToBytesSmallEnd(0))
		com := utils.EncodeFieldsToGtBasePacket(fileds)

		if !utils.EqualByteSlice(parts, com) {
			t.Errorf("should read %v but got %v", com, parts)
		}

		c := client.CreateGtBaseClient(fd, client.CreateAddress("127.0.0.1", port))
		m := replic.CreateMaster(10, 5000, -1, c)

		rs := replic.CreateReplicState()
		rs.SetMasterLock(m)

		args := analyzer.CreateCommandAssignArgs(c, rs, "127.0.0.1", 9855)

		analyzer.CreateHeartAnalyzer(nil, nil, -1, args).Analyze().Exec()
	}()

	nfd, erra := utils.Accept(int(lfd))
	if erra != nil {
		t.Errorf(erra.Error())
	}

	c := client.CreateGtBaseClient(nfd, client.CreateAddress("127.0.0.1", port))
	s := replic.CreateSlave(0, 0, 1, c)
	s.SetSyncStateLock(constants.SlaveSync)

	errh := s.SendHeartToSlave()
	if errh != nil {
		t.Errorf(errh.Error())
	}

	rs := replic.CreateReplicState()
	rs.AppendSlaveLock(s)

	parts := make([][]byte, 0)
	parts = append(parts, utils.Encodeint32ToBytesSmallEnd(0))
	parts = append(parts, utils.Encodeint32ToBytesSmallEnd(10))
	parts = append(parts, utils.Encodeint32ToBytesSmallEnd(5000))
	parts = append(parts, utils.Encodeint32ToBytesSmallEnd(1))

	analyzer.CreateGetHeartAnalyzer(parts, nil, 1, analyzer.CreateCommandAssignArgs(c, rs, "127.0.0.1", port)).Analyze().Exec()

	logIdx, logOff := s.GetLogInfo()
	if logIdx != 10 || logOff != 5000 {
		t.Errorf("should get %v idx %v off but got %v idx %v off", 10, 5000, logIdx, logOff)
	}
}

func TestRedoAnalyzer(t *testing.T) {
	port := 6544
	lfd, err := utils.BindAndListen(port)
	if err != nil {
		t.Errorf(err.Error())
	}

	redo, errr := page.ReadRedoPage(0)
	if errr != nil {
		t.Errorf(errr.Error())
	}

	slaveRedo := redo.Src()[:19]
	ch := make(chan struct{})
	go func() {
		fd, err := utils.LocalDial(port)
		if err != nil {
			t.Errorf(err.Error())
		}
		defer utils.CloseFd(fd)

		result := make([]byte, 0)
		fileds := make([][]byte, 0)
		for {
			res, errr := utils.ReadFd(fd)
			if errr != nil {
				t.Errorf(errr.Error())
			}
			result = append(result, res...)
			if utils.EqualByteSlice(result[len(result)-2:], []byte(constants.ReplicRedoLogEnd)) {
				fileds = utils.DecodeGtBasePacket(result[:len(result)-2])
				break
			}
		}
		seqBts := fileds[1]

		c := client.CreateGtBaseClient(fd, client.CreateAddress("127.0.0.1", port))
		m := replic.CreateMaster(0, int32(len(slaveRedo)), 1, c)
		rs := replic.CreateReplicState()
		rs.SetMasterLock(m)

		parts := make([][]byte, 0)
		parts = append(parts, seqBts)
		parts = append(parts, result)

		analyzer.CreateRedoAnalyzer(parts, nil, -1, analyzer.CreateCommandAssignArgs(c, rs, "127.0.0.1", port)).Analyze().Exec()
		ch <- struct{}{}
	}()

	nfd, erra := utils.Accept(int(lfd))
	if erra != nil {
		t.Errorf(erra.Error())
	}

	c := client.CreateGtBaseClient(nfd, client.CreateAddress("127.0.0.1", port))
	s := replic.CreateSlave(0, int32(len(slaveRedo)), 1, c)

	errh := s.SendRedoLogToSlave()
	if errh != nil {
		t.Errorf(errh.Error())
	}

	<-ch

	pg, err := page.ReadRedoPage(0)
	if err != nil {
		t.Errorf(err.Error())
	}

	if !utils.EqualByteSlice(pg.Src(), redo.Src()) {
		t.Errorf("not same")
	}
}

func TestBecomeSlaveAnalyzer(t *testing.T) {
	portm := 1695
	ports := 5212
	go func() {
		err := server.CreateGtBaseServer("127.0.0.1", portm).Run()
		if err != nil {
			t.Errorf(err.Error())
		}
	}()

	cmd := make([][]byte, 0)
	cmd = append(cmd, []byte("127.0.0.1"))
	cmd = append(cmd, utils.Encodeint32ToBytesSmallEnd(int32(portm)))
	c := client.CreateGtBaseClient(-1, client.CreateAddress("127.0.0.1", portm))
	rs := replic.CreateReplicState()

	time.Sleep(1 * time.Second)
	analyzer.CreateBecomeSlaveAnalyzer(cmd, nil, -1, analyzer.CreateCommandAssignArgs(c, rs, "127.0.0.1", ports)).Analyze().Exec()

	fileds := make([][]byte, 0)
	fileds = append(fileds, []byte(constants.HeartCommand))
	fileds = append(fileds, utils.Encodeint32ToBytesSmallEnd(0))
	result := utils.EncodeFieldsToGtBasePacket(fileds)

	listener, err := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(ports))
	if err != nil {
		t.Errorf(err.Error())
	}

	con, err := listener.Accept()
	if err != nil {
		t.Errorf(err.Error())
	}

	for i := 0; i < 2; i++ {
		res := make([]byte, 1024)
		n, err := con.Read(res)
		if err != nil {
			t.Errorf(err.Error())
		}

		if i >= 1 {
			if !utils.EqualByteSlice(res[:n], result) {
				t.Errorf("should get %v but got %v", result, res)
			}
		}
	}
}
