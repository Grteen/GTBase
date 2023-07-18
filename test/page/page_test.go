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

func TestReadWritePage(t *testing.T) {
	ph := page.CreatePageHeader(0)
	pg := page.ReadPage(&ph)

	data := readWritePageCreateData()

	for _, d := range data {
		pg.SetSrc(d)
		page.WritePage(pg)

		spg := page.ReadPage(&ph)

		if !utils.EqualByteSlice(spg.Src(), d) {
			t.Errorf("WritePage should write %s but ReadPage reads %s", d, spg.Src())
		}
	}

}
