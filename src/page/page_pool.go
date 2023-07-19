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
	pool.caches[p.pageHeader.PageIndex()] = p
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

// read the page from disk according to the pageIndex
func ReadPage(idx int32) (*Page, error) {
	p := readPageFromCache(idx)
	if p != nil {
		return p, nil
	}

	pd, err := readPageFromDisk(idx)

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

func readPageFromDisk(idx int32) (*Page, error) {
	var pageOffset int64 = CalOffsetOfIndex(idx)
	file, err := os.OpenFile(PageFilePathToDo, os.O_RDWR, 0777)
	if err != nil {
		return nil, glog.Error("ReadPage can't open PageFile because %s\n", err)
	}
	defer file.Close()

	src, err := readOnePageOfBytes(file, pageOffset)
	if err != nil {
		return nil, glog.Error("readOnePageOfBytes can't read because %s\n", err)
	}

	return CreatePage(idx, src), nil
}

func readOnePageOfBytes(f *os.File, offset int64) ([]byte, error) {
	result := make([]byte, PageSize)
	_, err := f.ReadAt(result, offset)
	if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
		return nil, err
	}

	return result, nil
}

func FlushPage(idx int32) error {
	pg, ok := GetPagePool().GetPage(idx)
	if !ok {
		return glog.Error("FlushPage can't flush because page%v not in PagePool", idx)
	}

	if !pg.pageHeader.dirty {
		return glog.Error("FlushPage don't need to flush because page%v not dirty", idx)
	}

	return WritePage(pg)
}

// write the page back to the disk
func WritePage(page *Page) error {
	file, err := os.OpenFile(PageFilePathToDo, os.O_RDWR, 0777)
	if err != nil {
		return glog.Error("WritePage can't open PageFile because %s\n", err)
	}
	defer file.Close()

	_, err = file.WriteAt(page.Src(), CalOffsetOfIndex(page.pageHeader.PageIndex()))
	if err != nil {
		return glog.Error("WritePage can't write because %s\n", err)
	}

	return nil
}
