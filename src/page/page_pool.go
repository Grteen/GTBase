package page

import (
	"GtBase/pkg/constants"
	"GtBase/pkg/glog"
	"container/list"
	"context"
	"io"
	"os"
	"sync"
	"time"
)

// PagePool caches all pages
// every read should read PagePool first
// if no cache in PagePool, it will read from disk and cache it
type PagePool struct {
	caches          map[int32]*PairPage
	cacheLock       sync.Mutex
	bucketCaches    map[int32]*BucketPage
	bucketCacheLock sync.Mutex
	redoCaches      map[int32]*RedoPage
	redoCacheLock   sync.Mutex

	dirtyList     *list.List
	dirtyListLock sync.Mutex
	lruList       *LRUList
}

func (pool *PagePool) GetPairPage(idx int32) (*PairPage, bool) {
	pool.cacheLock.Lock()
	defer pool.cacheLock.Unlock()
	p, ok := pool.caches[idx]
	return p, ok
}

func (pool *PagePool) GetBucketPage(idx int32) (*BucketPage, bool) {
	pool.bucketCacheLock.Lock()
	defer pool.bucketCacheLock.Unlock()
	p, ok := pool.bucketCaches[idx]
	return p, ok
}

func (pool *PagePool) GetRedoPage(idx int32) (*RedoPage, bool) {
	pool.redoCacheLock.Lock()
	defer pool.redoCacheLock.Unlock()
	p, ok := pool.redoCaches[idx]
	return p, ok
}

func (pool *PagePool) CachePairPage(p *PairPage) {
	pool.cacheLock.Lock()
	defer pool.cacheLock.Unlock()
	pool.caches[p.GetIndex()] = p

	delpageIdx := pool.lruList.Put(p.GetIndex())
	if delpageIdx != nil {
		delpg, ok := pool.GetPairPage(*delpageIdx)
		if !ok {
			return
		}
		go delpg.FlushPageLock()

		delete(pool.caches, *delpageIdx)
	}
}

func (pool *PagePool) CacheBucketPage(p *BucketPage) {
	pool.bucketCacheLock.Lock()
	defer pool.bucketCacheLock.Unlock()
	pool.bucketCaches[p.GetIndex()] = p

	delpageIdx := pool.lruList.Put(p.GetIndex())
	if delpageIdx != nil {
		delpg, ok := pool.GetPairPage(*delpageIdx)
		if !ok {
			return
		}
		go delpg.FlushPageLock()

		delete(pool.caches, *delpageIdx)
	}
}

func (pool *PagePool) CacheRedoPage(p *RedoPage) {
	pool.redoCacheLock.Lock()
	defer pool.redoCacheLock.Unlock()
	pool.redoCaches[p.GetIndex()] = p

	delpageIdx := pool.lruList.Put(p.GetIndex())
	if delpageIdx != nil {
		delpg, ok := pool.GetPairPage(*delpageIdx)
		if !ok {
			return
		}
		go delpg.FlushPageLock()

		delete(pool.caches, *delpageIdx)
	}
}

func (pool *PagePool) DirtyListPush(pg PageItf, cmn int32) {
	pool.dirtyListLock.Lock()
	defer pool.dirtyListLock.Unlock()

	pool.dirtyList.PushBack(CreateDirtyListNode(pg, cmn))
}

func (pool *PagePool) DirtyListGet() (*DirtyListNode, error) {
	pool.dirtyListLock.Lock()
	defer pool.dirtyListLock.Unlock()

	front := pool.dirtyList.Front()
	if front != nil {
		result, ok := pool.dirtyList.Remove(front).(*DirtyListNode)
		if !ok {
			return nil, glog.Error("can't transform to *Page")
		}
		return result, nil
	}

	return nil, nil
}

func CreatePagePool() *PagePool {
	return &PagePool{caches: map[int32]*PairPage{}, bucketCaches: map[int32]*BucketPage{}, dirtyList: list.New(), lruList: CreateLRUList(constants.PagePoolDefaultCapcity)}
}

var instance *PagePool
var once sync.Once

func GetPagePool() *PagePool {
	once.Do(func() {
		instance = CreatePagePool()
	})
	return instance
}

func readOnePageOfBytes(f *os.File, offset int64) ([]byte, error) {
	result := make([]byte, constants.PageSize)
	_, err := f.ReadAt(result, offset)
	if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
		return nil, err
	}

	return result, nil
}

func FlushDirtyList(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(100 * time.Millisecond)
			node, err := GetPagePool().DirtyListGet()
			if err != nil {
				return
			}
			if node == nil {
				continue
			}

			node.GetPage().FlushPageLock()
		}
	}
}
