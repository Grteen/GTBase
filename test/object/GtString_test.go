package object

import (
	"GtBase/src/object"
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
			t.Errorf("Length() should return %d but got %d", d.res, testTarget.Length())
		}
	}
}
