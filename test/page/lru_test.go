package page

import (
	"GtBase/src/page"
	"testing"
)

func TestLRUList(t *testing.T) {
	var a int32 = -1
	var b int32 = 1
	var c int32 = 4
	data := []struct {
		key int32
		res *int32
	}{
		{1, nil}, {-1, nil}, {3, nil}, {1, nil}, {4, nil}, {2, nil}, {5, &a}, {3, nil}, {6, &b}, {7, &c},
	}
	l := page.CreateLRUList(5)

	for _, d := range data {
		res := l.Put(d.key)
		if res == nil {
			continue
		}
		if *res != *(d.res) {
			t.Errorf("Put should get %v but got %v", *d.res, *res)
		}
	}
}
