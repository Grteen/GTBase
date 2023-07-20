package page

import (
	"log"
	"os"
)

const (
	PageFilePathToDo       string = "E:/Code/GTCDN/GTbase/temp/gt.pf"
	BucketPageFilePathToDo string = "E:/Code/GTCDN/GTbase/temp/gt.bf"
	PageSize               int64  = 16384
)

// Page is the basic unit store in disk and in xxx.pf file
// It is always 16KB
type Page struct {
	pageHeader *PageHeader
	src        []byte
	flushPath  string
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

func (p *Page) DirtyPage() {
	p.pageHeader.dirty = true
}

func (p *Page) CleanPage() {
	p.pageHeader.dirty = false
}

func (p *Page) Dirty() bool {
	return p.pageHeader.dirty
}

func (p *Page) GetIndex() int32 {
	return p.pageHeader.PageIndex()
}

func (p *Page) GetFlushPath() string {
	return p.flushPath
}

// also Dirty the Page
func (p *Page) WriteBytes(off int32, bts []byte) {
	// ToDo should ensure the consistency
	for i := 0; i < len(bts); i++ {
		p.src[i+int(off)] = bts[i]
	}
	p.DirtyPage()
}

// write the page back to the disk
// also clean the page
func (p *Page) writePage() error {
	return writePage(p, p.flushPath)
}

func CreatePage(idx int32, src []byte) *Page {
	return createPage(idx, src, PageFilePathToDo)
}

func CreateBucketPage(idx int32, src []byte) *Page {
	return createPage(idx, src, BucketPageFilePathToDo)
}

func createPage(idx int32, src []byte, flushPath string) *Page {
	ph := CreatePageHeader(idx)
	result := &Page{pageHeader: &ph, src: src, flushPath: flushPath}
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
	initPageFile(PageFilePathToDo)
}

func InitBucketPageFile() {
	initPageFile(BucketPageFilePathToDo)
}

func initPageFile(filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		_, errc := os.Create(filePath)
		if errc != nil {
			log.Fatalf("InitPageFile can't create the PageFile because %s\n", err)
		}

		errm := os.Chmod(filePath, 0777)
		if errm != nil {
			log.Fatalf("InitPageFile can't chmod because of %s\n", errm)
		}
	}
}

func DeletePageFile() {
	deletePageFile(PageFilePathToDo)
}

func DeleteBucketPageFile() {
	deletePageFile(BucketPageFilePathToDo)
}

func deletePageFile(filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Fatalf("PageFile not exist")
	}

	errr := os.Remove(filePath)
	if errr != nil {
		log.Fatalf("DeletePageFile can't remove the PageFile because %s\n", errr)
	}
}
