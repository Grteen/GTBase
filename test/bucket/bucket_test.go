package bucket

import (
	"GtBase/src/bucket"
	"GtBase/src/page"
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

		pg, err := page.ReadPage(-(d.firstHashValue + 1))
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
}
