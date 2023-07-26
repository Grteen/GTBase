package page

import (
	"GtBase/pkg/constants"
	"GtBase/pkg/glog"
	"log"
	"math"
	"os"
	"sync"
)

type PageItf interface {
	FlushPageLock() error
}

// Page is the basic unit store in disk and in xxx.pf file
// It is always 16KB
type Page struct {
	pageHeader *PageHeader
	src        []byte
	lock       sync.RWMutex
	flushPath  string
}

func (p *Page) Src() []byte {
	return p.src
}

func (p *Page) SetSrc(bts []byte) {
	p.src = bts
}

func (p *Page) SetCMN(cmn int32) {
	p.pageHeader.mu.Lock()
	defer p.pageHeader.mu.Unlock()
	p.pageHeader.cmn = cmn
}

func (p *Page) GetCMN() int32 {
	return p.pageHeader.cmn
}

func (p *Page) Lock() {
	p.lock.Lock()
}

func (p *Page) RLock() {
	p.lock.RLock()
}

func (p *Page) UnLock() {
	p.lock.Unlock()
}

func (p *Page) RUnLock() {
	p.lock.RUnlock()
}

// Dirty the page and push it to the dirtyList
func (p *Page) DirtyPageLock() {
	p.pageHeader.mu.Lock()
	defer p.pageHeader.mu.Unlock()
	p.pageHeader.dirty = true
	// TODO
	// should be push in cmn
	GetPagePool().DirtyListPush(p, p.GetCMN())
}

func (p *Page) CleanPageLock() {
	p.pageHeader.mu.Lock()
	defer p.pageHeader.mu.Unlock()
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

func (p *Page) IsBucket() bool {
	return p.flushPath == constants.BucketPageFilePathToDo
}

func (p *Page) SetFlushPath(flushPath string) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.flushPath = flushPath
}

func (p *Page) SrcSlice(start, end int32) []byte {
	return p.src[start:end]
}

func (p *Page) SrcSliceLength(start, length int32) []byte {
	return p.src[start : start+length]
}

// also Dirty the Page
func (p *Page) WriteBytes(off int32, bts []byte) {
	// ToDo should ensure the consistency
	for i := 0; i < len(bts); i++ {
		p.src[i+int(off)] = bts[i]
	}
	p.DirtyPageLock()
}

// write the page back to the disk
// also clean the page
func (p *Page) writePageLock() error {
	p.lock.Lock()
	defer p.lock.Unlock()
	file, err := os.OpenFile(p.GetFlushPath(), os.O_RDWR, 0777)
	if err != nil {
		return glog.Error("WritePage can't open PageFile because %s\n", err)
	}
	defer file.Close()

	_, err = file.WriteAt(p.Src(), CalOffsetOfIndex(p.GetIndex()))
	if err != nil {
		return glog.Error("WritePage can't write because %s\n", err)
	}

	p.CleanPageLock()

	return nil
}

// also clean the page
func (p *Page) FlushPageLock() error {
	// if !p.Dirty() {
	// 	return glog.Error("FlushPage don't need to flush because page%v not dirty", p.GetIndex())
	// }
	return p.writePageLock()
}

func IsBucketFilePath(filePath string) bool {
	return filePath == constants.BucketPageFilePathToDo
}

func CreatePage(idx int32, src []byte, flushPath string) *Page {
	ph := CreatePageHeader(idx)
	result := &Page{pageHeader: &ph, src: src, flushPath: flushPath}
	return result
}

// PageHeader is the header info of a Page
type PageHeader struct {
	pageIndex int32
	dirty     bool
	mu        sync.Mutex
	cmn       int32
}

func (ph *PageHeader) PageIndex() int32 {
	return ph.pageIndex
}

func (ph *PageHeader) GetCMN() int32 {
	return ph.cmn
}

func CalOffsetOfIndex(idx int32) int64 {
	idxc := idx
	if idx < 0 {
		idxc += 1
	}
	return int64(math.Abs(float64(idxc))) * constants.PageSize
}

func CreatePageHeader(idx int32) PageHeader {
	return PageHeader{pageIndex: idx, dirty: false}
}

func InitPageFile() {
	initPageFile(constants.PageFilePathToDo)
}

func InitBucketPageFile() {
	initPageFile(constants.BucketPageFilePathToDo)
}

func InitRedoLog() {
	initPageFile(constants.RedoLogToDo)
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
	deletePageFile(constants.PageFilePathToDo)
}

func DeleteBucketPageFile() {
	deletePageFile(constants.BucketPageFilePathToDo)
}

func DeleteRedoLog() {
	deletePageFile(constants.RedoLogToDo)
}

func deletePageFile(filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return
	}

	errr := os.Remove(filePath)
	if errr != nil {
		log.Fatalf("DeletePageFile can't remove the PageFile because %s\n", errr)
	}
}
