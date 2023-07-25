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

		firstIdx, firstOff, errf := bucket.FindFirstRecordRLock(object.CreateGtString(d.key))
		if errf != nil {
			t.Errorf(errf.Error())
		}

		pg, errr := page.ReadPairPage(firstIdx)
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

	firstIdx, firstOff, err := bucket.FindFirstRecordRLock(p.Key())
	if err != nil {
		t.Errorf(err.Error())
	}

	p, loc, errf := command.FindFinalRecord(firstIdx, firstOff)
	if errf != nil {
		t.Errorf(errf.Error())
	}

	if loc.GetOff() != p.CalMidOffset(int32(9*len(p.ToByte()))) {
		t.Errorf("FindFinalRecord should got %v off but got %v off", p.CalMidOffset(int32(10*len(p.ToByte()))), loc.GetOff())
	}

	pg, errr := page.ReadPairPage(0)
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

func TestSet(t *testing.T) {
	page.DeleteBucketPageFile()
	page.DeletePageFile()
	page.InitBucketPageFile()
	page.InitPageFile()

	data := []struct {
		key string
		val string
	}{
		{"Key", "Val"},
		{"Hello", "World"},
		{"Good", "Morning"},
	}

	for _, d := range data {
		err := command.Set(object.CreateGtString(d.key), object.CreateGtString(d.val), -1)
		if err != nil {
			t.Errorf(err.Error())
		}
	}

	for _, d := range data {
		firstIdx, firstOff, err := bucket.FindFirstRecordRLock(object.CreateGtString(d.key))
		if err != nil {
			t.Errorf(err.Error())
		}

		prevp, _, errf := command.FindFinalRecord(firstIdx, firstOff)
		if errf != nil {
			t.Errorf(errf.Error())
		}

		if prevp.Value().ToString() != d.val || prevp.Key().ToString() != d.key {
			t.Errorf("FindFinalRecord should got %v %v but got %v %v", d.key, d.val, prevp.Key().ToString(), prevp.Value().ToString())
		}
	}
}

func TestSameKeySet(t *testing.T) {
	page.DeleteBucketPageFile()
	page.DeletePageFile()
	page.InitBucketPageFile()
	page.InitPageFile()

	data := []struct {
		key string
		val string
		res string
	}{
		{"Key", "Val", "Morning"},
		{"Hello", "World", "World"},
		{"Key", "Morning", "Morning"},
	}

	for _, d := range data {
		err := command.Set(object.CreateGtString(d.key), object.CreateGtString(d.val), -1)
		if err != nil {
			t.Errorf(err.Error())
		}
	}

	for _, d := range data {
		val, err := command.Get(object.CreateGtString(d.key))
		if err != nil {
			t.Errorf(err.Error())
		}

		if val.ToString() != d.res {
			t.Errorf("Get should get %v but got %v", d.res, val.ToString())
		}
	}
}
