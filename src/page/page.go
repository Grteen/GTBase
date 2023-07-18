package page

import (
	"io"
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

func CreatePage(ph *PageHeader, src []byte) *Page {
	result := &Page{}
	result.SetPageHeader(ph)
	result.SetSrc(src)
	return result
}

// PageHeader is the header info of a Page
type PageHeader struct {
	pageIndex int32
}

func (ph *PageHeader) PageIndex() int32 {
	return ph.pageIndex
}

func (ph *PageHeader) SetPageIndex(idx int32) {
	ph.pageIndex = idx
}

func (ph *PageHeader) CalOffsetOfIndex() int64 {
	return int64(ph.PageIndex()) * PageSize
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

// read the page from disk according to the pageIndex
func ReadPage(ph *PageHeader) *Page {
	// TODO: should read page from pageCache First
	var pageOffset int64 = ph.CalOffsetOfIndex()
	file, err := os.OpenFile(PageFilePathToDo, os.O_RDWR, 0777)
	if err != nil {
		log.Fatalf("ReadPage can't open PageFile because %s\n", err)
	}
	defer file.Close()

	return CreatePage(ph, readOnePageOfBytes(file, pageOffset))
}

func readOnePageOfBytes(f *os.File, offset int64) []byte {
	result := make([]byte, PageSize)
	_, err := f.ReadAt(result, offset)
	if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
		log.Fatalf("readOnePageOfBytes can't read because %s\n", err)
	}

	return result
}

// write the page back to the disk
func WritePage(page *Page) {
	file, err := os.OpenFile(PageFilePathToDo, os.O_RDWR, 0777)
	if err != nil {
		log.Fatalf("WritePage can't open PageFile because %s\n", err)
	}
	defer file.Close()

	_, err = file.WriteAt(page.Src(), page.pageHeader.CalOffsetOfIndex())
	if err != nil {
		log.Fatalf("WritePage can't write because %s\n", err)
	}
}
