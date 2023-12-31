package nextwrite

import (
	"GtBase/pkg/constants"
	"GtBase/src/nextwrite"
	"GtBase/src/page"
	"GtBase/src/redo"
	"os"
	"testing"
)

func TestInitCMNFile(t *testing.T) {
	nextwrite.InitCMNFile()

	if _, err := os.Stat(constants.CMNPathToDo); os.IsNotExist(err) {
		t.Errorf("InitCMNFile() should create the %s but it didn't", constants.CMNPathToDo)
	}
}

func TestFactorySingleton(t *testing.T) {
	fa1 := nextwrite.GetNextWriteFactory()
	fa2 := nextwrite.GetNextWriteFactory()
	if fa1 != fa2 {
		t.Errorf("GetNextWriteFactory's address should be same but got %p and %p", fa1, fa2)
	}
}

func TestReadWriteCMN(t *testing.T) {
	nextwrite.DeleteCMNFile()
	nextwrite.InitCMNFile()
	data := []int32{1, 2, 3, 4, 5, 6, 7, 8, 9}
	for _, d := range data {
		result, err := nextwrite.GetCMN()
		if err != nil {
			t.Errorf(err.Error())
		}

		if result != d {
			t.Errorf("GetCMN should get %v but got %v", d, result)
		}
	}
}

func TestInitNextWrite(t *testing.T) {
	page.DeletePageFile()
	page.InitPageFile()
	redo.DeleteRedoLog()
	redo.InitRedoLog()

	err := nextwrite.InitNextWrite()
	if err != nil {
		t.Errorf(err.Error())
	}

	nw, err := nextwrite.GetNextWriteAndIncreaseIt(0)
	if err != nil {
		t.Errorf(err.Error())
	}
	idx, off := nw.NextWriteInfo()
	if idx != 0 || off != 0 {
		t.Errorf("GetNextWrite should get %v idx %v off but got %v idx %v off", 0, 0, idx, off)
	}
}

func TestIncreaseNextWrite(t *testing.T) {
	page.DeletePageFile()
	page.InitPageFile()
	err := nextwrite.InitNextWrite()
	if err != nil {
		t.Errorf(err.Error())
	}

	data := []struct {
		off int32
		res int32
	}{
		{10, 10},
		{20, 30},
		{105, 135},
	}

	for _, d := range data {
		err := nextwrite.IncreaseNextWrite(d.off)
		if err != nil {
			t.Errorf(err.Error())
		}

		nw, err := nextwrite.GetNextWriteAndIncreaseIt(0)
		if err != nil {
			t.Errorf(err.Error())
		}
		_, off := nw.NextWriteInfo()

		if off != d.res {
			t.Errorf("IncreaseNextWrite should increase offset to %v but got %v", d.res, off)
		}
	}

	err2 := nextwrite.IncreaseNextWrite(1000000)
	if err2 == nil {
		t.Errorf("should got error but none")
	}
}

func TestGetNextWrite(t *testing.T) {
	page.DeletePageFile()
	page.InitPageFile()
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
		nw, err := nextwrite.GetNextWriteAndIncreaseIt(d.off)
		if err != nil {
			t.Errorf(err.Error())
		}
		idx, off := nw.NextWriteInfo()
		if idx != d.residx || off != d.resoff {
			t.Errorf("GetNextWrite should got %v idx %v off but got %v idx %v off", d.residx, d.resoff, idx, off)
		}
	}
}
