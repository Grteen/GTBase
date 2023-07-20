package page

import (
	"log"
	"os"
)

const (
	PageFilePathToDo string = "E:/Code/GTCDN/GTbase/temp/gt.pf"
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

func (p *Page) Dirty() {
	p.pageHeader.dirty = true
}

func (p *Page) WriteBytes(off int32, bts []byte) {
	// ToDo should ensure the consistency
	for i := 0; i < len(bts); i++ {
		p.src[i] = bts[i]
	}
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
	return PageHeader{pageIndex: idx, dirty: false}
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

func DeletePageFile() {
	if _, err := os.Stat(PageFilePathToDo); os.IsNotExist(err) {
		log.Fatalf("PageFile not exist")
	}

	errr := os.Remove(PageFilePathToDo)
	if errr != nil {
		log.Fatalf("DeletePageFile can't remove the PageFile because %s\n", errr)
	}
}
