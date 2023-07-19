package utils

import (
	"GtBase/utils"
	"testing"
)

func TestInt32(t *testing.T) {
	data := []int32{1, 2, 4, -256, 512, 2147483647}
	for _, d := range data {
		bts := utils.Encodeint32ToBytesSmallEnd(d)
		res := utils.EncodeBytesSmallEndToint32(bts)
		if res != d {
			t.Errorf("EncodeBytesSmallEndToint32 should got %d but got %d", d, res)
		}
	}
}
