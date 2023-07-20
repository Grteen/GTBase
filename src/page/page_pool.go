package page

import (
	"GtBase/pkg/glog"
	"io"
	"os"
	"sync"
)

// PagePool caches all pages
// every read should read PagePool first
// if no cache in PagePool, it will read from disk and cache it
type PagePool struct {
	caches map[int32]*Page
}

func (pool *PagePool) GetPage(idx int32) (*Page, bool) {
	p, ok := pool.caches[idx]
	return p, ok
}

func (pool *PagePool) CachePage(p *Page) {
	pool.caches[p.GetIndex()] = p
}

func CreatePagePool() *PagePool {
	return &PagePool{caches: map[int32]*Page{}}
}

var instance *PagePool
var once sync.Once

func GetPagePool() *PagePool {
	once.Do(func() {
		instance = CreatePagePool()
	})
	return instance
}

// read the page from cache first
// if it not exist, read page from disk and cache it
func ReadPage(idx int32) (*Page, error) {
	return readPage(idx, PageFilePathToDo)
}

func ReadBucketPage(idx int32) (*Page, error) {
	return readPage(idx, BucketPageFilePathToDo)
}

func readPage(idx int32, filePath string) (*Page, error) {
	p := readPageFromCache(idx)
	if p != nil {
		return p, nil
	}

	pd, err := readPageFromDisk(idx, filePath)

	if err != nil {
		return nil, err
	}

	GetPagePool().CachePage(pd)

	return pd, nil
}

func readPageFromCache(idx int32) *Page {
	p, ok := GetPagePool().GetPage(idx)
	if !ok {
		return nil
	}

	return p
}

func readPageFromDisk(idx int32, filePath string) (*Page, error) {
	var pageOffset int64 = CalOffsetOfIndex(idx)
	file, err := os.OpenFile(filePath, os.O_RDWR, 0777)
	if err != nil {
		return nil, glog.Error("ReadPage can't open PageFile because %s\n", err)
	}
	defer file.Close()

	src, err := readOnePageOfBytes(file, pageOffset)
	if err != nil {
		return nil, glog.Error("readOnePageOfBytes can't read because %s\n", err)
	}

	return CreatePage(idx, src, filePath), nil
}

func readOnePageOfBytes(f *os.File, offset int64) ([]byte, error) {
	result := make([]byte, PageSize)
	_, err := f.ReadAt(result, offset)
	if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
		return nil, err
	}

	return result, nil
}

// also clean the page
// func FlushPage(idx int32) error {
// 	pg := readPageFromCache(idx)
// 	if pg == nil {
// 		return glog.Error("FlushPage can't flush because page%v not in PagePool", idx)
// 	}

// 	return pg.flushPage()
// }

func (p *Page) FlushPage() error {
	if !p.Dirty() {
		return glog.Error("FlushPage don't need to flush because page%v not dirty", p.GetIndex())
	}

	return p.writePage()
}

// if page not in cache, it will read page from disk
// Dirty the target page too
func WriteBytesToPageMemory(idx, off int32, bts []byte) error {
	pg, err := ReadPage(idx)
	if err != nil {
		return err
	}

	pg.WriteBytes(off, bts)

	return nil
}
