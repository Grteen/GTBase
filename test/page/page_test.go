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
		page.FlushPage(int32(idx))

		spg, err := page.ReadPage(ph.PageIndex())
		if err != nil {
			t.Errorf(err.Error())
		}

		if !utils.EqualByteSlice(spg.Src(), d) {
			t.Errorf("WritePage should write %s but ReadPage reads %s", d, spg.Src())
		}
	}
}

func TestWritePage(t *testing.T) {
	pg, err := page.ReadPage(0)
	if err != nil {
		t.Errorf(err.Error())
	}

}
