package nextwrite

import (
	"GtBase/pkg/glog"
	"encoding/binary"
	"os"
	"sync"
)

// NextWrite tell the Next Write Command where to write
type NextWrite struct {
	pageIndex  int32
	pageOffset int32
}

func (nw *NextWrite) NextWriteInfo() (int32, int32) {
	return nw.pageIndex, nw.pageOffset
}

const (
	CMNPathToDo string = "./temp/gt.cmn"
)

// NextWriteFactory assign CMN to all write command
// and assign NextWrite to all Set command
type NextWriteFactory struct {
	commandNumber int32
	nextWrite     NextWrite
	nwLock        sync.Mutex
	cmnLock       sync.Mutex
}

var instance *NextWriteFactory
var once sync.Once

func GetNextWriteFactory() *NextWriteFactory {
	once.Do(func() {
		// TODO init it's nextWrite
		instance = &NextWriteFactory{commandNumber: -1}
	})

	return instance
}

// getCMN will get the current commandNumber and atomically increase it
// func (nwf *NextWriteFactory) getCMN() int32 {
// 	nwf.cmnLock.Lock()
// 	defer nwf.cmnLock.Unlock()

// }

func (nwf *NextWriteFactory) checkCMNandInit() error {
	if !nwf.checkCMNInit() {
		return nwf.initCMN()
	}

	return nil
}

func (nwf *NextWriteFactory) checkCMNInit() bool {
	if nwf.commandNumber == -1 {
		return false
	}

	return true
}

func (nwf *NextWriteFactory) initCMN() error {
	cmn, err := nwf.readCMN()
	if err != nil {
		return err
	}

	nwf.commandNumber = cmn
	return nil
}

func (nwf *NextWriteFactory) readCMN() (int32, error) {
	file, err := os.Open(CMNPathToDo)
	if err != nil {
		return -1, glog.Error("ReadCMNFile can't open file %v because %v", CMNPathToDo, err)
	}
	defer file.Close()

	var result int32
	errr := binary.Read(file, binary.LittleEndian, &result)
	if errr != nil {
		return -1, glog.Error("ReadCMNFile can't read file %v because %v", CMNPathToDo, errr)
	}

	return result, nil
}

func (nwf *NextWriteFactory) writeCMN() error {
	file, err := os.OpenFile(CMNPathToDo, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return glog.Error("writeCMNFile can't open file %v because %v", CMNPathToDo, err)
	}
	defer file.Close()

	errw := binary.Write(file, binary.LittleEndian, nwf.commandNumber)
	if errw != nil {
		return glog.Error("writeCMNFile can't write file %v because %v", CMNPathToDo, errw)
	}

	return nil
}
