package nextwrite

import (
	"GtBase/src/nextwrite"
	"GtBase/src/redo"
	"testing"
)

func TestRedoNextWrite(t *testing.T) {
	redo.DeleteRedoLog()
	redo.InitRedoLog()
	err := nextwrite.InitNextWrite()
	if err != nil {
		t.Errorf(err.Error())
	}

	data := []struct {
		off    int32
		residx int32
		resoff int32
	}{
		{10, 0, 0},
		{20, 0, 10},
		{10000, 0, 30},
		{70, 0, 10030},
		{10000, 1, 0},
	}

	for _, d := range data {
		nw, err := nextwrite.GetRedoNextWriteAndIncreaseIt(d.off)
		if err != nil {
			t.Errorf(err.Error())
		}
		idx, off := nw.NextWriteInfo()
		if idx != d.residx || off != d.resoff {
			t.Errorf("GetRedoNextWriteAndIncreaseIt should got %v idx %v off but got %v idx %v off", d.residx, d.resoff, idx, off)
		}
	}
}
