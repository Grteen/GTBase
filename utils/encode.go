package utils

import (
	"encoding/binary"
)

// return a 4 length []byte
func Encodeint32ToBytesSmallEnd(x int32) []byte {
	bytes := make([]byte, 4)
	bytes[0] = byte(x)
	bytes[1] = byte(x >> 8)
	bytes[2] = byte(x >> 16)
	bytes[3] = byte(x >> 24)
	return bytes
}

func EncodeBytesSmallEndToint32(x []byte) int32 {
	return int32(binary.LittleEndian.Uint32(x))
}

func EncodeBytesSmallEndToInt8(x []byte) int8 {
	return int8(x[0])
}

// size fields size files....
// size is 4 byte to declare length of fileds behind it
func EncodeFieldsToGtBasePacket(x [][]byte) []byte {
	result := make([]byte, 0)
	for _, p := range x {
		result = append(result, Encodeint32ToBytesSmallEnd(int32(len(p)))...)
		result = append(result, p...)
	}

	return result
}
