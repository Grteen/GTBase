package nextwrite

import (
	"GtBase/src/nextwrite"
	"GtBase/src/page"
	"os"
	"testing"
)

func TestInitCMNFile(t *testing.T) {
	nextwrite.InitCMNFile()

	if _, err := os.Stat(nextwrite.CMNPathToDo); os.IsNotExist(err) {
		t.Errorf("InitCMNFile() should create the %s but it didn't", nextwrite.CMNPathToDo)
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
	data := []int32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
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
	page.InitPageFile()
	err := nextwrite.InitNextWrite()
	if err != nil {
		t.Errorf(err.Error())
	}
	nw := nextwrite.GetNextWrite(0)
	idx, off := nw.NextWriteInfo()
	if idx != 3 || off != 0 {
		t.Errorf("TEMP ERROR")
	}
}

func TestIncreaseNextWrite(t *testing.T) {
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

		nw := nextwrite.GetNextWrite(0)
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
