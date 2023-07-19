package nextset

import "sync"

// NextWrite tell the Next Write Command where to write
type NextWrite struct {
	pageIndex  int32
	pageOffset int32
}

func (nw *NextWrite) NextSetInfo() (int32, int32) {
	return nw.pageIndex, nw.pageOffset
}

// NextWriteFactory assign CMN to all write command
// and assign NextWrite to all Set command
type NextWriteFactory struct {
	// commandNumber int32
	nextWrite NextWrite
	nwLock    sync.Mutex
}

var instance *NextWriteFactory
var once sync.Once

func getNextWriteFactory() *NextWriteFactory {
	once.Do(func() {
		// TODO init it's nextWrite
		instance = &NextWriteFactory{}
	})

	return instance
}
