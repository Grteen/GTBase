package page

import (
	"io"
	"log"
	"os"
	"sync"
)

// PagePool caches all pages
// every read should read PagePool first
// if no cache in PagePool, it will read from disk and cache it
type PagePool struct {
	caches map[int32]*Page
}

func CreatePagePool() *PagePool {
	return &PagePool{}
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
func ReadPage(idx int32) *Page {
	// TODO: should read page from PagePool First
	return readPageFromDisk(idx)
}

// func readPageFromCache(idx int32) *Page {

// }

func readPageFromDisk(idx int32) *Page {
	var pageOffset int64 = CalOffsetOfIndex(idx)
	file, err := os.OpenFile(PageFilePathToDo, os.O_RDWR, 0777)
	if err != nil {
		log.Fatalf("ReadPage can't open PageFile because %s\n", err)
	}
	defer file.Close()

	return CreatePage(idx, readOnePageOfBytes(file, pageOffset))
}

func readOnePageOfBytes(f *os.File, offset int64) []byte {
	result := make([]byte, PageSize)
	_, err := f.ReadAt(result, offset)
	if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
		log.Fatalf("readOnePageOfBytes can't read because %s\n", err)
	}

	return result
}
