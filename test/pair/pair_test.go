package pair

import (
	"GtBase/src/object"
	"GtBase/src/pair"
	"GtBase/utils"
	"encoding/binary"
	"fmt"
	"testing"
)

func TestPairToByte(t *testing.T) {
	data := createTestPairToByteData()

	for _, d := range data {
		p := pair.CreatePair(d.key, d.val, d.flag, pair.CreateOverFlow(d.overFlowIndex, d.overFlowOffset))
		if !utils.EqualByteSlice(p.ToByte(), d.res) {
			t.Errorf("Pair.ToByte should got %v but got %v", d.res, p.ToByte())
		}
	}
}

type PairByteTest struct {
	key            object.Object
	val            object.Object
	flag           int8
	overFlowIndex  int32
	overFlowOffset int32
	res            []byte
}

func createTestPairToByteData() []PairByteTest {
	data := []PairByteTest{
		{object.CreateGtString("Hello World"), object.CreateGtString("Good World"), 0, 0, 30, make([]byte, 0)},
		{object.CreateGtString("Key"), object.CreateGtString("Val"), 1, 0, 30, make([]byte, 0)},
	}
	data[0].res = append(data[0].res, []byte{71, 111, 111, 100, 32, 87, 111, 114, 108, 100}...)
	data[0].res = append(data[0].res, []byte{10, 0, 0, 0}...)
	data[0].res = append(data[0].res, []byte{72, 101, 108, 108, 111, 32, 87, 111, 114, 108, 100}...)
	data[0].res = append(data[0].res, []byte{11, 0, 0, 0}...)
	data[0].res = append(data[0].res, []byte{0, 0, 0, 0}...)
	data[0].res = append(data[0].res, []byte{0, 0, 0, 0}...)
	data[0].res = append(data[0].res, []byte{30, 0, 0, 0}...)

	data[1].res = append(data[1].res, []byte{86, 97, 108}...)
	data[1].res = append(data[1].res, []byte{3, 0, 0, 0}...)
	data[1].res = append(data[1].res, []byte{75, 101, 121}...)
	data[1].res = append(data[1].res, []byte{3, 0, 0, 0}...)
	data[1].res = append(data[1].res, []byte{1, 0, 0, 0}...)
	data[1].res = append(data[1].res, []byte{0, 0, 0, 0}...)
	data[1].res = append(data[1].res, []byte{30, 0, 0, 0}...)
	return data
}

func TestBinaryAppend(t *testing.T) {
	result := make([]byte, 0, 32)
	str := "Hello Worldasdsad"
	fmt.Println(len(str))
	result = binary.AppendVarint(result, int64(len(str)))

	fmt.Println(result)
}
