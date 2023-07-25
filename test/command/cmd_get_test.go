package command

import (
	"GtBase/src/command"
	"GtBase/src/object"
	"GtBase/src/page"
	"testing"
)

func TestGet(t *testing.T) {
	page.DeleteBucketPageFile()
	page.DeletePageFile()
	page.InitBucketPageFile()
	page.InitPageFile()

	data := []struct {
		key string
		val string
	}{
		{"Key", "Val"},
		{"Hello", "World"},
		{"Good", "Morning"},
	}

	for _, d := range data {
		err := command.Set(object.CreateGtString(d.key), object.CreateGtString(d.val), -1)
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
