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
		{[]byte("Set Key Val\r\n"), []byte("Ok")},
		{[]byte("Get Key\r\n"), []byte("Val")},
		{[]byte("Del Key\r\n"), []byte("Ok")},
		{[]byte("Get Key\r\n"), []byte("Nil")},
	}

	go func() {
		s := server.CreateGtBaseServer()
		err := s.Run(1235)
		if err != nil {
			t.Errorf(err.Error())
		}
	}()

	time.Sleep(1 * time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:1235")
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
