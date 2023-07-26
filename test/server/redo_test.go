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
	go page.FlushDirtyList(ctx)

	dataw := []struct {
		cmn int32
		cmd []byte
	}{
		{1, []byte("Set key val")},
		{2, []byte("Set Hello World")},
		{3, []byte("Del Hello World")},
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
		{[]byte("Get key"), "val"},
		{[]byte("Get Hello"), constants.ServerGetNilReturn},
	}

	for _, d := range data {
		res := analyzer.CreateCommandAssign(d.cmd, -1).Assign().Analyze().Exec().ToString()
		if res != d.res {
			t.Errorf("Exec should get %v but got %v", d.res, res)
		}
	}

	cancel()
}
