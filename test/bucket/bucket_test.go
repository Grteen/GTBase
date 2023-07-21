package bucket

import (
	"GtBase/src/bucket"
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
		b := bucket.CreateBucket(d.idx, d.off)
		if !utils.EqualByteSlice(d.res, b.ToByte()) {
			t.Errorf("ToByte should got %v but got %v", d.res, b.ToByte())
		}
	}
}
