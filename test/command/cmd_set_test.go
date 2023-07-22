package command

import (
	"GtBase/src/bucket"
	"GtBase/src/command"
	"GtBase/src/object"
	"GtBase/src/page"
	"GtBase/src/pair"
	"testing"
)

func TestFirstSet(t *testing.T) {
	data := []struct {
		key  string
		val  string
		flag int8
	}{
		{"Key", "Val", 0},
		{"Hello", "World", 0},
	}

	for _, d := range data {
		p := pair.CreatePair(object.CreateGtString(d.key), object.CreateGtString(d.val), d.flag, pair.CreateNullOverFlow())
		err := command.FirstSetInThisBucket(p)
		if err != nil {
			t.Errorf(err.Error())
		}

		firstIdx, firstOff, errf := bucket.FindFirstRecord(object.CreateGtString(d.key))
		if errf != nil {
			t.Errorf(errf.Error())
		}

		pg, errr := page.ReadPage(firstIdx)
		if errr != nil {
			t.Errorf(errr.Error())
		}

		pread := pair.ReadPair(pg, firstOff)
		if pread.Key().ToString() != p.Key().ToString() {
			t.Errorf("ReadPair should got Key %v but got %v", p.Key().ToString(), pread.Key().ToString())
		}
		if pread.Value().ToString() != p.Value().ToString() {
			t.Errorf("ReadPair should got Value %v but got %v", p.Value().ToString(), pread.Value().ToString())
		}
		if pread.Flag() != p.Flag() {
			t.Errorf("ReadPair should got Flag %v but got %v", p.Flag(), pread.Flag())
		}
		if overIdx, overOff := pread.OverFlow().OverFlowInfo(); overIdx != 0 || overOff != 0 {
			t.Errorf("ReadPair should got OverFlow %v idx %v off but got %v idx %v off", 0, 0, overIdx, overOff)
		}
	}
}
