package nextwrite

import (
	"GtBase/pkg/constants"
	"GtBase/pkg/glog"
	"GtBase/utils"
	"encoding/binary"
	"log"
	"os"
	"sync"
	"sync/atomic"
)

// NextWrite tell the Next Write Command where to write
type NextWrite struct {
	pageIndex  int32
	pageOffset int32
}

func (nw *NextWrite) NextWriteInfo() (int32, int32) {
	return nw.pageIndex, nw.pageOffset
}

func CreateNextWrite(pageIndex, pageOffset int32) *NextWrite {
	return &NextWrite{pageIndex: pageIndex, pageOffset: pageOffset}
}

const (
	CMNPathToDo string = "E:/Code/GTCDN/GTbase/temp/gt.cmn"
)

// NextWriteFactory assign CMN to all write command
// and assign NextWrite to all Set command
type NextWriteFactory struct {
	commandNumber int32
	nextWrite     NextWrite
	nwLock        sync.Mutex
	cmnLock       sync.Mutex
}

// getCMN will get the current commandNumber and atomically increase it
func (nwf *NextWriteFactory) getCMN() (int32, error) {
	nwf.cmnLock.Lock()
	defer nwf.cmnLock.Unlock()

	err := nwf.checkCMNandInit()
	if err != nil {
		return -1, err
	}

	result := nwf.commandNumber
	// TODO
	// when writeCMN error file and commandNumber may be inconsistency ?
	atomic.AddInt32(&nwf.commandNumber, 1)
	errw := nwf.writeCMN()
	if errw != nil {
		return -1, err
	}

	return result, nil
}

func (nwf *NextWriteFactory) checkCMNandInit() error {
	if nwf.commandNumber == -1 {
		return nwf.initCMN()
	}

	return nil
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
	file, err := os.OpenFile(CMNPathToDo, os.O_RDWR, 0777)
	if err != nil {
		return -1, glog.Error("ReadCMNFile can't open file %s because %s", CMNPathToDo, err.Error())
	}
	defer file.Close()

	var result int32
	errr := binary.Read(file, binary.LittleEndian, &result)
	if errr != nil {
		return -1, glog.Error("ReadCMNFile can't read file %s because %s", CMNPathToDo, errr.Error())
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

var instance *NextWriteFactory
var once sync.Once

func GetNextWriteFactory() *NextWriteFactory {
	once.Do(func() {
		// TODO init it's nextWrite
		instance = &NextWriteFactory{commandNumber: -1}
	})

	return instance
}

func GetCMN() (int32, error) {
	return GetNextWriteFactory().getCMN()
}

func InitCMNFile() {
	if _, err := os.Stat(CMNPathToDo); os.IsNotExist(err) {
		file, errc := os.Create(CMNPathToDo)
		if errc != nil {
			log.Fatalf("InitCMNFile can't create the CMNFILE because %s\n", err)
		}

		errm := os.Chmod(CMNPathToDo, 0777)
		if errm != nil {
			log.Fatalf("InitCMNFile can't chmod because of %s\n", errm)
		}

		errw := binary.Write(file, binary.LittleEndian, utils.Encodeint32ToBytesSmallEnd(0))
		if errw != nil {
			log.Fatalf("writeCMNFile can't write file %v because %v", CMNPathToDo, errw)
		}
	}
}

func DeleteCMNFile() {
	if _, err := os.Stat(CMNPathToDo); os.IsNotExist(err) {
		return
	}

	errr := os.Remove(CMNPathToDo)
	if errr != nil {
		log.Fatalf("DeletePageFile can't remove the CMNFile because %s\n", errr)
	}
}

func InitNextWrite() error {
	GetNextWriteFactory().nwLock.Lock()
	defer GetNextWriteFactory().nwLock.Unlock()
	fileInfo, err := os.Stat(constants.PageFilePathToDo)
	if err != nil {
		return glog.Error("InitNextWrite can't Stat file %v becasuse %v", constants.PageFilePathToDo, err)
	}

	fileSize := fileInfo.Size()
	GetNextWriteFactory().initNextWriteIndexAndOffset(fileSize)
	return nil
}

// InitNextWrite must initialize a new page no matter what the last page is
func (nwf *NextWriteFactory) initNextWriteIndexAndOffset(fileSize int64) {
	pageIndex := int32(fileSize / constants.PageSize)
	pageOffset := 0
	nwf.nextWrite = *CreateNextWrite(pageIndex, int32(pageOffset))
}

func (nwf *NextWriteFactory) getNextWrite() *NextWrite {
	return CreateNextWrite(nwf.nextWrite.pageIndex, nwf.nextWrite.pageOffset)
}

// if last page's size don't satisfy the size will write in
// return a new page index to it
// if size to write is bigger than PageSize
// refuse it and return an error
func (nwf *NextWriteFactory) checkRestSizeAndChange(off int32) error {
	if off > int32(constants.PageSize) {
		return glog.Error("The Size to Write %v bytes is bigger than %v", off, constants.PageSize)
	}
	idx, offInPage := nwf.getNextWrite().NextWriteInfo()
	if offInPage+off > int32(constants.PageSize) {
		nwf.nextWrite = *CreateNextWrite(idx+1, 0)
	}

	return nil
}

func getNextWrite(off int32) (*NextWrite, error) {
	err := GetNextWriteFactory().checkRestSizeAndChange(off)
	if err != nil {
		return nil, err
	}
	return GetNextWriteFactory().getNextWrite(), nil
}

// Get the NextWrite and increase is locked
// so page don't need to lock except the delete command because every set command has different offset and index
func GetNextWriteAndIncreaseIt(off int32) (*NextWrite, error) {
	GetNextWriteFactory().nwLock.Lock()
	defer GetNextWriteFactory().nwLock.Unlock()

	result, err := getNextWrite(off)
	if err != nil {
		return nil, err
	}

	erri := IncreaseNextWrite(off)
	if erri != nil {
		return nil, err
	}

	return result, nil
}

// GetNextWrite ensure that the size of the write does not exceed the size of the page
func (nwf *NextWriteFactory) increaseNextWrite(off int32) error {
	idx, oldOff := nwf.nextWrite.NextWriteInfo()
	if oldOff+off > int32(constants.PageSize) {
		return glog.Error("increaseNextWrite write %v bytes but page rest %v bytes", off, constants.PageSize-int64(oldOff))
	}

	nwf.nextWrite = *CreateNextWrite(idx, oldOff+off)
	return nil
}

func IncreaseNextWrite(off int32) error {
	return GetNextWriteFactory().increaseNextWrite(off)
}
