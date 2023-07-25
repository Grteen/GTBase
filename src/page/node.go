package page

import (
	"container/list"
	"sync"
)

type DirtyListNode struct {
	pg        *Page
	oldestCMN int32
}

func (n *DirtyListNode) GetOldestCMN() int32 {
	return n.oldestCMN
}

func (n *DirtyListNode) GetPage() *Page {
	return n.pg
}

func CreateDirtyListNode(pg *Page, oldestCMN int32) *DirtyListNode {
	return &DirtyListNode{pg, oldestCMN}
}

type LRUListNode struct {
	pgIdx int32
}

func (n *LRUListNode) GetPageIndex() int32 {
	return n.pgIdx
}

func CreateLRUListNode(pgIdx int32) *LRUListNode {
	return &LRUListNode{pgIdx: pgIdx}
}

type LRUList struct {
	cap  int32
	list *list.List
	mp   map[int32]*list.Element
	mu   sync.Mutex
}

// Put page's index into LRUList and return the page's index that should be deleted in cache
func (l *LRUList) Put(pgIdx int32) *int32 {
	l.mu.Lock()
	defer l.mu.Unlock()

	node, ok := l.mp[pgIdx]
	if ok {
		l.list.MoveToBack(node)
		return nil
	}

	ele := l.list.PushBack(CreateLRUListNode(pgIdx))
	l.mp[pgIdx] = ele

	if l.list.Len() > int(l.cap) {
		if e := l.list.Front(); e != nil {
			l.list.Remove(e)
			node := e.Value.(*LRUListNode)
			delete(l.mp, node.GetPageIndex())

			result := node.GetPageIndex()
			return &result
		}
	}

	return nil
}

func CreateLRUList(cap int32) *LRUList {
	return &LRUList{
		cap:  cap,
		list: list.New(),
		mp:   make(map[int32]*list.Element),
	}
}
