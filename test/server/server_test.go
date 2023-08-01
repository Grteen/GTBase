package server

import (
	"GtBase/src/page"
	"GtBase/src/server"
	"GtBase/utils"
	"net"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	page.DeletePageFile()
	page.DeleteBucketPageFile()
	page.DeleteRedoLog()
	page.InitBucketPageFile()
	page.InitPageFile()
	page.InitRedoLog()
	data := []struct {
		command []byte
		result  []byte
	}{
		{[]byte{3, 0, 0, 0, 83, 101, 116, 3, 0, 0, 0, 75, 101, 121, 3, 0, 0, 0, 86, 97, 108, 13, 10}, []byte("Ok")},
		{[]byte{3, 0, 0, 0, 71, 101, 116, 3, 0, 0, 0, 75, 101, 121, 13, 10}, []byte("Val")},
	}

	go func() {
		s := server.CreateGtBaseServer("127.0.0.1", 2222)
		err := s.Run()
		if err != nil {
			t.Errorf(err.Error())
		}
	}()

	time.Sleep(1 * time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:2222")
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	for _, d := range data {
		_, err := conn.Write(d.command)
		if err != nil {
			t.Errorf(err.Error())
		}

		buf := make([]byte, 1024)
		_, err = conn.Read(buf)
		if err != nil {
			t.Errorf(err.Error())
		}

		if !utils.EqualByteSliceOnlyInMinLen(buf, d.result) {
			t.Errorf("Read should get %v but got %v", d.result, buf)
		}
	}

	time.Sleep(2 * time.Second)
}
