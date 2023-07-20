package page

import (
	"GtBase/src/page"
	"os"
	"testing"

	"GtBase/utils"
)

func TestInitPageFile(t *testing.T) {
	page.InitPageFile()

	if _, err := os.Stat(page.PageFilePathToDo); os.IsNotExist(err) {
		t.Errorf("InitPageFile() should create the %s but it didn't", page.PageFilePathToDo)
	}
}

func TestReadWritePage(t *testing.T) {
	testReadWritePageInSingleIndex(t, 0)
	testReadWritePageInSingleIndex(t, 1)
	testReadWritePageInSingleIndex(t, 2)
}

func readWritePageCreateData() [][]byte {
	result := make([][]byte, 0)
	data := []string{"Hello World", "abc"}

	for _, d := range data {
		t := make([]byte, page.PageSize)
		for i := 0; i < len(d); i++ {
			t[i] = d[i]
		}
		result = append(result, t)
	}

	return result
}

func testReadWritePageInSingleIndex(t *testing.T, idx int) {
	ph := page.CreatePageHeader(int32(idx))
	pg, err := page.ReadPage(ph.PageIndex())
	if err != nil {
		t.Errorf(err.Error())
	}

	data := readWritePageCreateData()

	for _, d := range data {
		pg.SetSrc(d)
		pg.FlushPage()

		spg, err := page.ReadPage(ph.PageIndex())
		if err != nil {
			t.Errorf(err.Error())
		}

		if !utils.EqualByteSlice(spg.Src(), d) {
			t.Errorf("WritePage should write %s but ReadPage reads %s", d, spg.Src())
		}
	}
}

func TestWriteBytes(t *testing.T) {
	page.DeletePageFile()
	page.InitPageFile()
	data := []struct {
		write []byte
		res   []byte
	}{
		{[]byte(""), []byte("")},
		{[]byte("First Write "), []byte("First Write ")},
		{[]byte("Second Write "), []byte("First Write Second Write ")},
		{[]byte("Hello World"), []byte("First Write Second Write Hello World")},
	}

	for i := 1; i < len(data); i++ {
		pg, err := page.ReadPage(1)
		if err != nil {
			t.Errorf(err.Error())
		}

		p, ok := page.GetPagePool().GetPage(1)
		if !ok {
			t.Errorf("GetPagePool should get index %v but not", 1)
		}

		if pg != p {
			t.Errorf("GetPagePool().GetPage() should be same as page.ReadBucketPage but not")
		}

		pg.WriteBytes(int32(len(data[i-1].res)), data[i].write)
		if pg.Dirty() != true {
			t.Errorf("page should be dirtied by WriteBytes but not")
		}

		if !utils.EqualByteSliceOnlyInMinLen(data[i].res, pg.Src()) {
			t.Errorf("page should be %v but it got %v", data[i].res, pg.Src()[:len(data[i].res)])
		}
	}

	pg, err := page.ReadPage(0)
	if err != nil {
		t.Errorf(err.Error())
	}

	pg.FlushPage()
	if pg.Dirty() != false {
		t.Errorf("page should be cleaned by FlushPage but not")
	}
}

func TestBucketWriteBytes(t *testing.T) {
	page.DeleteBucketPageFile()
	page.InitBucketPageFile()
	data := []struct {
		write []byte
		res   []byte
	}{
		{[]byte(""), []byte("")},
		{[]byte("First Write "), []byte("First Write ")},
		{[]byte("Second Write "), []byte("First Write Second Write ")},
		{[]byte("Hello World"), []byte("First Write Second Write Hello World")},
	}

	for i := 1; i < len(data); i++ {
		pg, err := page.ReadBucketPage(1)
		p, ok := page.GetPagePool().GetPage(-1)
		if !ok {
			t.Errorf("GetPagePool should get index %v but not", -1)
		}

		if pg != p {
			t.Errorf("GetPagePool().GetPage() should be same as page.ReadBucketPage but not")
		}

		if err != nil {
			t.Errorf(err.Error())
		}

		pg.WriteBytes(int32(len(data[i-1].res)), data[i].write)
		if pg.Dirty() != true {
			t.Errorf("page should be dirtied by WriteBytes but not")
		}

		if !utils.EqualByteSliceOnlyInMinLen(data[i].res, pg.Src()) {
			t.Errorf("page should be %v but it got %v", data[i].res, pg.Src()[:len(data[i].res)])
		}
	}

	pg, err := page.ReadBucketPage(0)
	if err != nil {
		t.Errorf(err.Error())
	}

	pg.FlushPage()
	if pg.Dirty() != false {
		t.Errorf("page should be cleaned by FlushPage but not")
	}
}
