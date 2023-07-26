package command

import (
	"GtBase/src/command"
	"GtBase/src/object"
	"GtBase/src/page"
	"testing"
)

func TestDel(t *testing.T) {
	page.DeleteBucketPageFile()
	page.DeletePageFile()
	page.InitBucketPageFile()
	page.InitPageFile()

	data := []struct {
		key string
		val string
		res bool
	}{
		{"Key", "Val", true},
		{"Hello", "World", false},
		{"Good", "Morning", true},
	}

	for _, d := range data {
		err := command.Set(object.CreateGtString(d.key), object.CreateGtString(d.val), -1)
		if err != nil {
			t.Errorf(err.Error())
		}
	}

	err := command.Del(object.CreateGtString(data[1].key), -1)
	if err != nil {
		t.Errorf(err.Error())
	}

	for _, d := range data {
		val, err := command.Get(object.CreateGtString(d.key))
		if err != nil {
			t.Errorf(err.Error())
		}

		if !d.res {
			if val != nil {
				t.Errorf("Get should get nil but got %v", val.ToString())
			}
			continue
		}

		if val.ToString() != d.val {
			t.Errorf("Get should get %v but got %v", d.val, val.ToString())
		}
	}
}
