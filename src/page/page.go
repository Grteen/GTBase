package page

import (
	"log"
	"os"
	"strings"
)

const (
	TempFileNameToDo string = "ttt.pf"
	FilePath         string = "./temp/"
	PageSize         int64  = 16384
)

// Page is the basic unit store in disk and in xxx.pf file
// It is always 16KB
type Page struct {
	pageHeader PageHeader
	src        []byte
}

// PageHeader is the header info of a Page
type PageHeader struct {
	pageIndex int32
}

func InitPageFile() {
	var filePath strings.Builder
	filePath.WriteString(FilePath)
	filePath.WriteString(TempFileNameToDo)

	if _, err := os.Stat(filePath.String()); os.IsNotExist(err) {
		_, errc := os.Create(filePath.String())
		if errc != nil {
			log.Fatalf("InitPageFile can't create the PageFile because %s\n", err)
		}
	}
}

// read the page from disk according to the pageIndex
// func ReadPage(ph PageHeader) *Page {

// 	// var pageOffset int64 = int64(ph.pageIndex) * PageSize

// }
