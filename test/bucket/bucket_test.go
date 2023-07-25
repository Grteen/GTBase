package bucket

import (
	"GtBase/src/bucket"
	"GtBase/src/nextwrite"
	"GtBase/src/object"
	"GtBase/src/page"
	"GtBase/src/pair"
	"GtBase/utils"
	"testing"
)

func TestBucketToByte(t *testing.T) {
	data := []struct {
		idx int32
		off int32
		res []byte
	}{
		{0, 0, []byte{0, 0, 0, 0, 0, 0, 0, 0}},
		{1, 26, []byte{1, 0, 0, 0, 26, 0, 0, 0}},
		{2, 257, []byte{2, 0, 0, 0, 1, 1, 0, 0}},
	}

	for _, d := range data {
		b := bucket.CreateBucket(nil, d.idx, d.off)
		if !utils.EqualByteSlice(d.res, b.ToByte()) {
			t.Errorf("ToByte should got %v but got %v", d.res, b.ToByte())
		}
	}
}

func TestWriteInPage(t *testing.T) {
	page.DeleteBucketPageFile()
	page.InitBucketPageFile()

	data := []struct {
		firstHashValue  int32
		secondHashValue int32
		idx             int32
		off             int32
		res             []byte
	}{
		{0, 0, 0, 0, []byte{0, 0, 0, 0, 0, 0, 0, 0}},
		{0, 1, 1, 26, []byte{1, 0, 0, 0, 26, 0, 0, 0}},
		{0, 2, 2, 257, []byte{2, 0, 0, 0, 1, 1, 0, 0}},
	}

	for i := 0; i < len(data); i++ {
		d := data[i]

		b := bucket.CreateBucket(bucket.CreateBucketHeader(d.firstHashValue, d.secondHashValue), d.idx, d.off)
		b.WriteInPage()

		pg, err := page.ReadBucketPage(b.BucketHeader().CalIndexOfBucketPage())
		if err != nil {
			t.Errorf(err.Error())
		}

		got := make([]byte, 0)
		for j := 0; j < i; j++ {
			got = append(got, data[j].res...)
		}

		if !utils.EqualByteSliceOnlyInMinLen(pg.Src(), got) {
			t.Errorf("ReadPage should got %v but got %v", got, pg.Src()[:len(got)])
		}
	}

	data = []struct {
		firstHashValue  int32
		secondHashValue int32
		idx             int32
		off             int32
		res             []byte
	}{
		{1, 0, 0, 0, []byte{0, 0, 0, 0, 0, 0, 0, 0}},
		{10, 0, 1, 26, []byte{1, 0, 0, 0, 26, 0, 0, 0}},
		{20, 0, 2, 257, []byte{2, 0, 0, 0, 1, 1, 0, 0}},
		{8, 0, 3, 30, []byte{3, 0, 0, 0, 30, 0, 0, 0}},
	}

	for i := 0; i < len(data); i++ {
		d := data[i]

		b := bucket.CreateBucket(bucket.CreateBucketHeader(d.firstHashValue, d.secondHashValue), d.idx, d.off)
		b.WriteInPage()

		pg, err := page.ReadBucketPage(b.BucketHeader().CalIndexOfBucketPage())
		if err != nil {
			t.Errorf(err.Error())
		}

		s := b.BucketHeader().CalOffsetOfBucketPage()
		e := s + int32(len(b.ToByte()))

		if !utils.EqualByteSlice(pg.Src()[s:e], b.ToByte()) {
			t.Errorf("ReadPage should got %v but got %v", b.ToByte(), pg.Src()[s:e])
		}
	}
}

func TestFindFirstRecord(t *testing.T) {
	page.DeleteBucketPageFile()
	page.DeletePageFile()
	page.InitBucketPageFile()
	page.InitPageFile()

	data := []struct {
		key string
		val string
	}{
		{"Key", "Value"},
		{"Hello", "World"},
	}

	for _, d := range data {
		p := pair.CreatePair(object.CreateGtString(d.key), object.CreateGtString(d.val), 0, pair.CreateOverFlow(0, 0))
		nw, err := nextwrite.GetNextWriteAndIncreaseIt(int32(len(p.ToByte())))
		if err != nil {
			t.Errorf(err.Error())
		}

		idx, off := nw.NextWriteInfo()

		b := bucket.CreateBucketByKey(p.Key(), idx, off)

		p.WriteInPage(idx, off)
		b.WriteInPage()

		idxf, offf, errf := bucket.FindFirstRecordRLock(p.Key())
		if errf != nil {
			t.Errorf(errf.Error())
		}

		pg, errr := page.ReadPairPage(idxf)
		if errr != nil {
			t.Errorf(errr.Error())
		}

		bts := pg.SrcSlice(offf, offf+int32(len(p.ToByte())))

		if !utils.EqualByteSlice(bts, p.ToByte()) {
			t.Errorf("SrcSlice should got %v but got %v", p.ToByte(), bts)
		}
	}

	zeroIdx, zeroOff, err := bucket.FindFirstRecordRLock(object.CreateGtString("Impossible"))
	if err != nil {
		t.Errorf(err.Error())
	}

	if zeroIdx != 0 || zeroOff != 0 {
		t.Errorf("FindFirstRecord should got 0 idx 0 off but got %v idx %v off", zeroIdx, zeroOff)
	}
}
