package command

import (
	"GtBase/src/bucket"
	"GtBase/src/command"
	"GtBase/src/nextwrite"
	"GtBase/src/object"
	"GtBase/src/page"
	"GtBase/src/pair"
	"testing"
)

func TestFirstSet(t *testing.T) {
	page.DeleteBucketPageFile()
	page.DeletePageFile()
	page.InitBucketPageFile()
	page.InitPageFile()

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

func TestFindFinalRecord(t *testing.T) {
	page.DeleteBucketPageFile()
	page.DeletePageFile()
	page.InitBucketPageFile()
	page.InitPageFile()

	p := pair.CreatePair(object.CreateGtString("First"), object.CreateGtString("Second"), 0, pair.CreateNullOverFlow())
	p.SetOverFlow(pair.CreateOverFlow(0, p.CalMidOffset(int32(len(p.ToByte())))))

	for i := 1; i <= 10; i++ {
		if i == 1 {
			errc := command.FirstSetInThisBucket(p)
			if errc != nil {
				t.Errorf(errc.Error())
			}
			continue
		}
		if i != 10 {
			p.SetOverFlow(pair.CreateOverFlow(0, p.CalMidOffset(int32(i*len(p.ToByte())))))
		} else {
			p.SetOverFlow(pair.CreateNullOverFlow())
		}
		nw, err := nextwrite.GetNextWriteAndIncreaseIt(int32(len(p.ToByte())))
		if err != nil {
			t.Errorf(err.Error())
		}

		idx, off := nw.NextWriteInfo()
		p.WriteInPage(idx, off)
	}

	firstIdx, firstOff, err := bucket.FindFirstRecord(p.Key())
	if err != nil {
		t.Errorf(err.Error())
	}

	p, errf := command.FindFinalRecord(firstIdx, firstOff)
	if errf != nil {
		t.Errorf(errf.Error())
	}

	pg, errr := page.ReadPage(0)
	if errr != nil {
		t.Errorf(errr.Error())
	}

	for i := 0; i < 10; i++ {
		p := pair.ReadPair(pg, p.CalMidOffset(int32(i*len(p.ToByte()))))
		_, off := p.OverFlow().OverFlowInfo()
		nextOff := p.CalMidOffset(int32((i + 1) * len(p.ToByte())))
		if i == 9 {
			nextOff = 0
		}
		if off != nextOff {
			t.Errorf("ReadPair should got %v but got %v", nextOff, off)
		}
	}
}
