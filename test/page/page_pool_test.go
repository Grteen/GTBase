package page

import (
	"GtBase/src/page"
	"GtBase/utils"
	"testing"
)

func TestPagePoolSingleton(t *testing.T) {
	pool := page.GetPagePool()
	pool2 := page.GetPagePool()
	if pool != pool2 {
		t.Errorf("GetPagePool's address should be same but got %p and %p", pool, pool2)
	}
}

func TestReadFlush(t *testing.T) {
	data := []string{"abc", "Hello World"}
	for _, d := range data {
		pg, err := page.ReadPage(0)
		if err != nil {
			t.Errorf(err.Error())
		}
		pg.SetSrc([]byte(d))
		pg.Dirty()
		page.FlushPage(0)

		pg2, err := page.ReadPage(0)
		if err != nil {
			t.Errorf(err.Error())
		}

		if !utils.EqualByteSlice(pg2.Src(), []byte(d)) {
			t.Errorf("FlushPage should write %v but ReadPage read %s", d, pg2.Src())
		}
	}
}
