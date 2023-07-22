package page

import (
	"GtBase/pkg/constants"
	"GtBase/pkg/glog"
	"GtBase/src/object"
	"GtBase/utils"
	"log"
	"math"
	"os"
)

const ()

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

func (p *Page) IsBucket() bool {
	return p.flushPath == constants.BucketPageFilePathToDo
}

func (p *Page) SetFlushPath(flushPath string) {
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
	p.DirtyPage()
}

// write the page back to the disk
// also clean the page
func (p *Page) writePage() error {
	file, err := os.OpenFile(p.GetFlushPath(), os.O_RDWR, 0777)
	if err != nil {
		return glog.Error("WritePage can't open PageFile because %s\n", err)
	}
	defer file.Close()

	_, err = file.WriteAt(p.Src(), CalOffsetOfIndex(p.GetIndex()))
	if err != nil {
		return glog.Error("WritePage can't write because %s\n", err)
	}

	p.CleanPage()

	return nil
}

// also clean the page
func (p *Page) FlushPage() error {
	if !p.Dirty() {
		return glog.Error("FlushPage don't need to flush because page%v not dirty", p.GetIndex())
	}

	return p.writePage()
}

func (p *Page) ReadFlag(off int32) int8 {
	return utils.EncodeBytesSmallEndToInt8(p.SrcSliceLength(off-constants.PairFlagSize, constants.PairFlagSize))
}

func (p *Page) ReadKeyLength(off int32) int32 {
	return utils.EncodeBytesSmallEndToint32(p.SrcSliceLength(off-constants.PairKeyLengthSize, constants.PairKeyLengthSize))
}

func (p *Page) ReadKey(off int32, keyLength int32) object.Object {
	return object.CreateGtStringByBytes(p.SrcSliceLength(off-keyLength, keyLength))
}

func (p *Page) ReadValLength(off int32) int32 {
	return utils.EncodeBytesSmallEndToint32(p.SrcSliceLength(off-constants.PairKeyLengthSize, constants.PairKeyLengthSize))
}

func (p *Page) ReadVal(off int32, valLength int32) object.Object {
	return object.CreateGtStringByBytes(p.SrcSliceLength(off-valLength, valLength))
}

func (p *Page) ReadOverFlow(off int32) (int32, int32) {
	overFlowIdx := utils.EncodeBytesSmallEndToint32(p.SrcSliceLength(off, constants.PairOverFlowIndexSize))
	overFlowOffset := utils.EncodeBytesSmallEndToint32(p.SrcSliceLength(off+constants.PairOverFlowIndexSize, constants.PairOverFlowOffsetSize))
	return overFlowIdx, overFlowOffset
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
}

func (ph *PageHeader) PageIndex() int32 {
	return ph.pageIndex
}

func (ph *PageHeader) SetPageIndex(idx int32) {
	ph.pageIndex = idx
}

func CalOffsetOfIndex(idx int32) int64 {
	idxc := idx
	if idx < 0 {
		idxc += 1
	}
	return int64(math.Abs(float64(idx))) * constants.PageSize
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

func deletePageFile(filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return
	}

	errr := os.Remove(filePath)
	if errr != nil {
		log.Fatalf("DeletePageFile can't remove the PageFile because %s\n", errr)
	}
}
