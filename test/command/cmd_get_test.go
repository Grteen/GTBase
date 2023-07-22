package command

import (
	"GtBase/src/command"
	"GtBase/src/object"
	"GtBase/src/pair"
	"testing"
)

func TestGet(t *testing.T) {
	data := []struct {
		key string
		val string
	}{
		{"Key", "Val"},
		{"Hello", "World"},
		{"Good", "Morning"},
	}

	for _, d := range data {
		p := pair.CreatePair(object.CreateGtString(d.key), object.CreateGtString(d.val), 0, pair.CreateNullOverFlow())
		err := command.Set(p)
		if err != nil {
			t.Errorf(err.Error())
		}
	}

	for _, d := range data {
		val, err := command.Get(object.CreateGtString(d.key))
		if err != nil {
			t.Errorf(err.Error())
		}

		if val.ToString() != d.val {
			t.Errorf("Get should get %v but got %v", d.val, val.ToString())
		}
	}
}
