package page

import (
	"log"
	"os"
)

const (
	PageFilePathToDo string = "./temp/gt.pf"
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

func (ph *PageHeader) CalOffsetOfIndex() int64 {
	return int64(ph.pageIndex) * PageSize
}

func InitPageFile() {
	if _, err := os.Stat(PageFilePathToDo); os.IsNotExist(err) {
		_, errc := os.Create(PageFilePathToDo)
		if errc != nil {
			log.Fatalf("InitPageFile can't create the PageFile because %s\n", err)
		}
	}
}

// // read the page from disk according to the pageIndex
// func ReadPage(ph PageHeader) *Page {

// 	// var pageOffset int64 = int64(ph.pageIndex) * PageSize

// }

// write the page back to the disk
func WritePage(page *Page) {
	file, err := os.Open(PageFilePathToDo)
	if err != nil {
		log.Fatalf("WritePage can't open PageFile because %s\n", err)
	}
	defer file.Close()

	_, err = file.WriteAt(page.src, 1)
	if err != nil {
		log.Fatalf("WritePage can't write because %s\n", err)
	}
}
