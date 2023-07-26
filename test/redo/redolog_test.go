package redo

import (
	"GtBase/pkg/constants"
	"GtBase/src/analyzer"
	"GtBase/src/nextwrite"
	"GtBase/src/page"
	"GtBase/src/redo"
	"GtBase/utils"
	"testing"
)

func TestRedoLog(t *testing.T) {
	nextwrite.DeleteCMNFile()
	nextwrite.InitCMNFile()
	page.DeleteRedoLog()
	page.InitRedoLog()

	data := []struct {
		cmd string
		res []byte
	}{
		{"Set key val", []byte{1, 0, 0, 0, 11, 0, 0, 0, 83, 101, 116, 32, 107, 101, 121, 32, 118, 97, 108}},
		{"Del key", []byte{1, 0, 0, 0, 11, 0, 0, 0, 83, 101, 116, 32, 107, 101, 121, 32, 118, 97, 108,
			2, 0, 0, 0, 7, 0, 0, 0, 68, 101, 108, 32, 107, 101, 121}},
	}

	for _, d := range data {
		cmn, errg := nextwrite.GetCMN()
		if errg != nil {
			t.Errorf(errg.Error())
		}

		ok := analyzer.CreateCommandAssign([]byte(d.cmd), cmn).Assign().Analyze().Exec().ToString()
		if ok != constants.ServerOkReturn {
			t.Errorf("Exec should get %v but got %v", constants.ServerOkReturn, ok)
		}

		pg, err := page.ReadRedoPage(0)
		if err != nil {
			t.Errorf(err.Error())
		}

		if !utils.EqualByteSliceOnlyInMinLen(pg.Src(), d.res) {
			t.Errorf("ReadRedoPage should get %v but got %v", d.res, pg.SrcSliceLength(0, int32(len(d.res))))
		}
	}

}

func TestReadRedo(t *testing.T) {
	nextwrite.DeleteCMNFile()
	nextwrite.InitCMNFile()
	page.DeleteRedoLog()
	page.InitRedoLog()

	data := []struct {
		cmd string
		len int32
		res *redo.Redo
	}{
		{"Set key val", 0, redo.CreateRedo(1, 11, []byte{83, 101, 116, 32, 107, 101, 121, 32, 118, 97, 108})},
		{"Del key", 19, redo.CreateRedo(2, 7, []byte{68, 101, 108, 32, 107, 101, 121})},
	}

	for _, d := range data {
		cmn, errg := nextwrite.GetCMN()
		if errg != nil {
			t.Errorf(errg.Error())
		}

		ok := analyzer.CreateCommandAssign([]byte(d.cmd), cmn).Assign().Analyze().Exec().ToString()
		if ok != constants.ServerOkReturn {
			t.Errorf("Exec should get %v but got %v", constants.ServerOkReturn, ok)
		}

		pg, err := page.ReadRedoPage(0)
		if err != nil {
			t.Errorf(err.Error())
		}

		r := redo.ReadRedo(pg, d.len)
		if r.GetCMN() != d.res.GetCMN() || !utils.EqualByteSlice(r.GetCmd(), d.res.GetCmd()) {
			t.Errorf("ReadRedo should get %v cmn %v cmd but got %v cmn %v cmd", d.res.GetCMN(), d.res.GetCmd(), r.GetCMN(), r.GetCmd())
		}
	}
}
