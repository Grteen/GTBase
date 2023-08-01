package server

import (
	"GtBase/pkg/constants"
	"GtBase/src/analyzer"
	"GtBase/src/nextwrite"
	"GtBase/src/page"
	"GtBase/src/redo"
	"GtBase/src/server"
	"context"
	"testing"
	"time"
)

func TestRedoLog(t *testing.T) {
	page.DeletePageFile()
	page.DeleteBucketPageFile()
	page.InitBucketPageFile()
	page.InitPageFile()
	page.DeleteCheckPointFile()
	page.InitCheckPointFile()
	page.DeleteRedoLog()
	page.InitRedoLog()
	nextwrite.DeleteCMNFile()
	nextwrite.InitCMNFile()

	ctx, cancel := context.WithCancel(context.Background())
	go page.FlushRedoDirtyList(ctx)

	dataw := []struct {
		cmn int32
		cmd []byte
	}{
		{1, []byte{3, 0, 0, 0, 83, 101, 116, 3, 0, 0, 0, 107, 101, 121, 3, 0, 0, 0, 118, 97, 108}},
		{2, []byte{3, 0, 0, 0, 83, 101, 116, 5, 0, 0, 0, 72, 101, 108, 108, 111, 5, 0, 0, 0, 87, 111, 114, 108, 100}},
		{3, []byte{3, 0, 0, 0, 68, 101, 108, 5, 0, 0, 0, 72, 101, 108, 108, 111}},
	}

	for _, d := range dataw {
		redo.WriteRedoLog(d.cmn, d.cmd)
	}

	time.Sleep(1 * time.Second)
	err := server.RedoLog()
	if err != nil {
		t.Errorf(err.Error())
	}

	data := []struct {
		cmd []byte
		res string
	}{
		{[]byte{3, 0, 0, 0, 71, 101, 116, 3, 0, 0, 0, 107, 101, 121}, "val"},
		{[]byte{3, 0, 0, 0, 71, 101, 116, 5, 0, 0, 0, 72, 101, 108, 108, 111}, constants.ServerGetNilReturn},
	}

	for _, d := range data {
		res := analyzer.CreateCommandAssign(d.cmd, -1, nil).Assign().Analyze().Exec().ToString()
		if res != d.res {
			t.Errorf("Exec should get %v but got %v", d.res, res)
		}
	}

	cancel()
}
