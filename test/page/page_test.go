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
	pg := page.ReadPage(ph.PageIndex())

	data := readWritePageCreateData()

	for _, d := range data {
		pg.SetSrc(d)
		page.WritePage(pg)

		spg := page.ReadPage(ph.PageIndex())

		if !utils.EqualByteSlice(spg.Src(), d) {
			t.Errorf("WritePage should write %s but ReadPage reads %s", d, spg.Src())
		}
	}
}
