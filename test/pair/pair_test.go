package pair

import (
	"GtBase/src/object"
	"GtBase/src/pair"
	"encoding/binary"
	"fmt"
	"testing"
)

func TestPairToByte(t *testing.T) {
	data := []struct {
		key            string
		val            string
		flag           int8
		overFlowIndex  int32
		overFlowOffset int32
		res            string
	}{
		{"1a", "1a", 0, 0, 30, "4 "},
	}

	p := pair.CreatePair(object.CreateGtString(data[0].key), object.CreateGtString(data[0].val), data[0].flag, pair.CreateOverFlow(data[0].overFlowIndex, data[0].overFlowOffset))
	fmt.Println((p.ToByte()))
}

func TestBinaryAppend(t *testing.T) {
	result := make([]byte, 0, 32)
	str := "Hello Worldasdsad"
	fmt.Println(len(str))
	result = binary.AppendVarint(result, int64(len(str)))

	fmt.Println(result)
}
