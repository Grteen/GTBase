package utils

import (
	"GtBase/pkg/constants"
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

func TestEncodePacket(t *testing.T) {
	fileds := [][]byte{[]byte(constants.SetCommand), []byte("key"), []byte("val")}
	result := []byte{3, 0, 0, 0, 83, 101, 116, 3, 0, 0, 0, 107, 101, 121, 3, 0, 0, 0, 118, 97, 108, 13, 10}

	res := utils.EncodeFieldsToGtBasePacket(fileds)
	if !utils.EqualByteSlice(res, result) {
		t.Errorf("should get %v but got %v", result, res)
	}
}

func TestDecodePacket(t *testing.T) {
	packet := []byte{3, 0, 0, 0, 83, 101, 116, 3, 0, 0, 0, 107, 101, 121, 3, 0, 0, 0, 118, 97, 108}
	result := [][]byte{[]byte(constants.SetCommand), []byte("key"), []byte("val")}

	res := utils.DecodeGtBasePacket(packet)

	if len(res) != len(result) {
		t.Errorf("should get %v but got %v", result, res)
	}

	for i := 0; i < len(res); i++ {
		if !utils.EqualByteSlice(res[i], result[i]) {
			t.Errorf("should get %v but got %v", result, res)
		}
	}
}
