package pair

import (
	"GtBase/src/nextwrite"
	"GtBase/src/object"
	"GtBase/src/page"
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
	data[0].res = append(data[0].res, []byte{0}...)
	data[0].res = append(data[0].res, []byte{0, 0, 0, 0}...)
	data[0].res = append(data[0].res, []byte{30, 0, 0, 0}...)

	data[1].res = append(data[1].res, []byte{86, 97, 108}...)
	data[1].res = append(data[1].res, []byte{3, 0, 0, 0}...)
	data[1].res = append(data[1].res, []byte{75, 101, 121}...)
	data[1].res = append(data[1].res, []byte{3, 0, 0, 0}...)
	data[1].res = append(data[1].res, []byte{1}...)
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

func TestWriteInPage(t *testing.T) {
	page.DeletePageFile()
	page.InitPageFile()
	data := createTestPairToByteData()

	for i := 0; i < len(data); i++ {
		d := data[i]

		p := pair.CreatePair(d.key, d.val, d.flag, pair.CreateOverFlow(d.overFlowIndex, d.overFlowOffset))
		bts := p.ToByte()
		nw, err := nextwrite.GetNextWriteAndIncreaseIt(int32(len(bts)))
		if err != nil {
			t.Errorf(err.Error())
		}

		idx, off := nw.NextWriteInfo()
		p.WriteInPage(idx, off)

		pg, errr := page.ReadPage(idx)
		if errr != nil {
			t.Errorf(err.Error())
		}

		got := make([]byte, 0)
		for j := 0; j <= i; j++ {
			got = append(got, data[j].res...)
		}

		if !utils.EqualByteSliceOnlyInMinLen(pg.Src(), got) {
			t.Errorf("ReadPage should got %v but got %v", pg.Src()[:len(p.ToByte())], got)
		}
	}
}

func TestIsDelete(t *testing.T) {
	data := []struct {
		arg int8
		res bool
	}{
		{0, false},
		{1, true},
		{2, false},
		{3, true},
	}

	for _, d := range data {
		if pair.IsDelete(d.arg) != d.res {
			t.Errorf("IsDelete %v should got %v but got %v", d.arg, d.res, pair.IsDelete(d.arg))
		}
	}
}
