package object

import (
	"GtBase/src/object"
	"GtBase/utils"
	"testing"
)

func TestGtStringLength(t *testing.T) {
	data := []struct {
		arg string
		res int
	}{
		{arg: "Hello World", res: 11},
		{arg: `abc`, res: 3},
	}

	for _, d := range data {
		testTarget := object.CreateGtString(d.arg)
		if testTarget.Length() != d.res {
			t.Errorf("GtString.Length() should return %d but got %d", d.res, testTarget.Length())
		}
	}
}

func TestGtStringToByte(t *testing.T) {
	data := []string{"Hello World", "abc"}
	for _, d := range data {
		testTarget := object.CreateGtString(d)
		bts := testTarget.ToByte()

		if !utils.EqualByteSlice(bts, testTarget.Value()) {
			t.Errorf("GtString.ToByte() should return %s but got %s", testTarget.Value(), bts)
		}
	}
}
