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
	page.DeletePageFile()
	page.InitPageFile()
	data := []string{"abc", "Hello World"}
	for _, d := range data {
		pg, err := page.ReadPage(0)
		if err != nil {
			t.Errorf(err.Error())
		}
		pg.SetSrc([]byte(d))
		pg.DirtyPageLock()
		errf := pg.FlushPage()
		if errf != nil {
			t.Errorf(errf.Error())
		}

		pg2, err := page.ReadPage(0)
		if err != nil {
			t.Errorf(err.Error())
		}

		if !utils.EqualByteSlice(pg2.Src(), []byte(d)) {
			t.Errorf("FlushPage should write %v but ReadPage read %s", d, pg2.Src()[:len([]byte(d))])
		}
	}
}

func TestDirtyList(t *testing.T) {
	page.DeletePageFile()
	page.InitPageFile()
	data := []*page.Page{page.CreatePage(0, nil, ""), page.CreatePage(1, nil, ""), page.CreatePage(2, nil, ""), page.CreatePage(-1, nil, "")}
	for _, d := range data {
		page.GetPagePool().DirtyListPush(d, -1)
	}

	for i := 0; i < len(data); i++ {
		res, err := page.GetPagePool().DirtyListGet()
		if err != nil {
			t.Errorf(err.Error())
		}

		if res.GetPage() != data[i] {
			t.Errorf("DirtyListGet not same")
		}
	}
}
