package command

import (
	"GtBase/src/command"
	"GtBase/src/object"
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
		err := command.Set(object.CreateGtString(d.key), object.CreateGtString(d.val))
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
