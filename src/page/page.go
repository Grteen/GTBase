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
	pageHeader *PageHeader
	src        []byte
}

func (p *Page) Src() []byte {
	return p.src
}

func (p *Page) SetSrc(bts []byte) {
	p.src = bts
}

func (p *Page) SetPageHeader(ph *PageHeader) {
	p.pageHeader = ph
}

func CreatePage(idx int32, src []byte) *Page {
	result := &Page{}
	ph := CreatePageHeader(idx)
	result.SetPageHeader(&ph)
	result.SetSrc(src)
	return result
}

// PageHeader is the header info of a Page
type PageHeader struct {
	pageIndex int32
	dirty     bool
}

func (ph *PageHeader) PageIndex() int32 {
	return ph.pageIndex
}

func (ph *PageHeader) SetPageIndex(idx int32) {
	ph.pageIndex = idx
}

func CalOffsetOfIndex(idx int32) int64 {
	return int64(idx) * PageSize
}

func CreatePageHeader(idx int32) PageHeader {
	var result PageHeader
	result.SetPageIndex(idx)
	return result
}

func InitPageFile() {
	if _, err := os.Stat(PageFilePathToDo); os.IsNotExist(err) {
		_, errc := os.Create(PageFilePathToDo)
		if errc != nil {
			log.Fatalf("InitPageFile can't create the PageFile because %s\n", err)
		}

		errm := os.Chmod(PageFilePathToDo, 0777)
		if errm != nil {
			log.Fatalf("InitPageFile can't chmod because of %s\n", errm)
		}
	}
}

// write the page back to the disk
func WritePage(page *Page) error {
	file, err := os.OpenFile(PageFilePathToDo, os.O_RDWR, 0777)
	if err != nil {
		return err
		log.Fatalf("WritePage can't open PageFile because %s\n", err)
	}
	defer file.Close()

	_, err = file.WriteAt(page.Src(), CalOffsetOfIndex(page.pageHeader.PageIndex()))
	if err != nil {
		return err
		log.Fatalf("WritePage can't write because %s\n", err)
	}

	return nil
}
