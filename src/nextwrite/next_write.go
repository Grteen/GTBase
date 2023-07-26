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

// NextWriteFactory assign CMN to all write command
// and assign NextWrite to all Set command
type NextWriteFactory struct {
	commandNumber int32
	cmnLock       sync.Mutex

	nextWrite        NextWrite
	nwLock           sync.Mutex
	redoLogNextWrite NextWrite
	redoLock         sync.Mutex
}

func (nwf *NextWriteFactory) SetNextWriteLock(nw NextWrite) {
	nwf.nwLock.Lock()
	defer nwf.nwLock.Unlock()
	nwf.nextWrite = nw
}

func (nwf *NextWriteFactory) SetRedoLogNextWriteLock(nw NextWrite) {
	nwf.redoLock.Lock()
	defer nwf.redoLock.Unlock()
	nwf.redoLogNextWrite = nw
}

func (nwf *NextWriteFactory) SetNextWrite(nw NextWrite) {
	nwf.nextWrite = nw
}

func (nwf *NextWriteFactory) SetRedoLogNextWrite(nw NextWrite) {
	nwf.redoLogNextWrite = nw
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
	file, err := os.OpenFile(constants.CMNPathToDo, os.O_RDWR, 0777)
	if err != nil {
		return -1, glog.Error("ReadCMNFile can't open file %s because %s", constants.CMNPathToDo, err.Error())
	}
	defer file.Close()

	var result int32
	errr := binary.Read(file, binary.LittleEndian, &result)
	if errr != nil {
		return -1, glog.Error("ReadCMNFile can't read file %s because %s", constants.CMNPathToDo, errr.Error())
	}

	return result, nil
}

func (nwf *NextWriteFactory) writeCMN() error {
	file, err := os.OpenFile(constants.CMNPathToDo, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return glog.Error("writeCMNFile can't open file %v because %v", constants.CMNPathToDo, err)
	}
	defer file.Close()

	errw := binary.Write(file, binary.LittleEndian, nwf.commandNumber)
	if errw != nil {
		return glog.Error("writeCMNFile can't write file %v because %v", constants.CMNPathToDo, errw)
	}

	return nil
}

func GetCMN() (int32, error) {
	return GetNextWriteFactory().getCMN()
}

func InitCMNFile() {
	if _, err := os.Stat(constants.CMNPathToDo); os.IsNotExist(err) {
		file, errc := os.Create(constants.CMNPathToDo)
		if errc != nil {
			log.Fatalf("InitCMNFile can't create the CMNFILE because %s\n", err)
		}

		errm := os.Chmod(constants.CMNPathToDo, 0777)
		if errm != nil {
			log.Fatalf("InitCMNFile can't chmod because of %s\n", errm)
		}

		errw := binary.Write(file, binary.LittleEndian, utils.Encodeint32ToBytesSmallEnd(0))
		if errw != nil {
			log.Fatalf("writeCMNFile can't write file %v because %v", constants.CMNPathToDo, errw)
		}
	}
}

func DeleteCMNFile() {
	if _, err := os.Stat(constants.CMNPathToDo); os.IsNotExist(err) {
		return
	}

	errr := os.Remove(constants.CMNPathToDo)
	if errr != nil {
		log.Fatalf("DeletePageFile can't remove the CMNFile because %s\n", errr)
	}
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

func InitNextWrite() error {
	err := initNextWrite()
	if err != nil {
		return err
	}

	erri := initRedoNextWrite()
	if erri != nil {
		return erri
	}

	return nil
}

func initNextWrite() error {
	GetNextWriteFactory().nwLock.Lock()
	defer GetNextWriteFactory().nwLock.Unlock()
	fileInfo, err := os.Stat(constants.PageFilePathToDo)
	if err != nil {
		return glog.Error("InitNextWrite can't Stat file %v becasuse %v", constants.PageFilePathToDo, err)
	}

	fileSize := fileInfo.Size()
	nw := GetNextWriteFactory().initNextWriteIndexAndOffset(fileSize)
	GetNextWriteFactory().SetNextWrite(*nw)

	return nil
}

func initRedoNextWrite() error {
	GetNextWriteFactory().redoLock.Lock()
	defer GetNextWriteFactory().redoLock.Unlock()
	fileInfor, errs := os.Stat(constants.RedoLogToDo)
	if errs != nil {
		return glog.Error("InitNextWrite can't Stat file %v becasuse %v", constants.RedoLogToDo, errs)
	}

	fileSizer := fileInfor.Size()
	nwr := GetNextWriteFactory().initNextWriteIndexAndOffset(fileSizer)
	GetNextWriteFactory().SetRedoLogNextWrite(*nwr)

	return nil
}

// InitNextWrite must initialize a new page no matter what the last page is
func (nwf *NextWriteFactory) initNextWriteIndexAndOffset(fileSize int64) *NextWrite {
	pageIndex := int32(fileSize / constants.PageSize)
	pageOffset := 0
	return CreateNextWrite(pageIndex, int32(pageOffset))
}

func (nwf *NextWriteFactory) getNextWrite() *NextWrite {
	return CreateNextWrite(nwf.nextWrite.pageIndex, nwf.nextWrite.pageOffset)
}

func (nwf *NextWriteFactory) getRedoNextWrite() *NextWrite {
	return CreateNextWrite(nwf.redoLogNextWrite.pageIndex, nwf.redoLogNextWrite.pageOffset)
}

// if last page's size don't satisfy the size will write in
// return a new page index to it
// if size to write is bigger than PageSize
// refuse it and return an error
func checkRestSizeAndChange(off, idx, offInPage int32) (int32, int32, error) {
	if off > int32(constants.PageSize) {
		return -1, -1, glog.Error("The Size to Write %v bytes is bigger than %v", off, constants.PageSize)
	}
	if offInPage+off > int32(constants.PageSize) {
		return idx + 1, 0, nil
	}

	return -1, -1, nil
}

func getNextWrite(off int32) (*NextWrite, error) {
	idx, offInPage := GetNextWriteFactory().getNextWrite().NextWriteInfo()
	newidx, newoff, err := checkRestSizeAndChange(off, idx, offInPage)
	if err != nil {
		return nil, err
	}

	if newidx != -1 && newoff != -1 {
		GetNextWriteFactory().SetNextWrite(*CreateNextWrite(newidx, newoff))
	}

	return GetNextWriteFactory().getNextWrite(), nil
}

func getRedoNextWrite(off int32) (*NextWrite, error) {
	idx, offInPage := GetNextWriteFactory().getRedoNextWrite().NextWriteInfo()
	newidx, newoff, err := checkRestSizeAndChange(off, idx, offInPage)
	if err != nil {
		return nil, err
	}

	if newidx != -1 && newoff != -1 {
		GetNextWriteFactory().SetRedoLogNextWrite(*CreateNextWrite(newidx, newoff))
	}

	return GetNextWriteFactory().getRedoNextWrite(), nil
}

// Get the NextWrite and increase is locked
// so page don't need to lock except the delete command because every set command has different offset and index
func GetNextWriteAndIncreaseIt(btsLength int32) (*NextWrite, error) {
	GetNextWriteFactory().nwLock.Lock()
	defer GetNextWriteFactory().nwLock.Unlock()

	result, err := getNextWrite(btsLength)
	if err != nil {
		return nil, err
	}

	erri := IncreaseNextWrite(btsLength)
	if erri != nil {
		return nil, err
	}

	return result, nil
}

func GetRedoNextWriteAndIncreaseIt(btsLength int32) (*NextWrite, error) {
	GetNextWriteFactory().redoLock.Lock()
	defer GetNextWriteFactory().redoLock.Unlock()

	result, err := getRedoNextWrite(btsLength)
	if err != nil {
		return nil, err
	}

	erri := IncreaseRedoNextWrite(btsLength)
	if erri != nil {
		return nil, err
	}

	return result, nil
}

// GetNextWrite ensure that the size of the write does not exceed the size of the page
func increaseNextWrite(off, idx, oldOff int32) (*NextWrite, error) {
	if oldOff+off > int32(constants.PageSize) {
		return nil, glog.Error("increaseNextWrite write %v bytes but page rest %v bytes", off, constants.PageSize-int64(oldOff))
	}

	return CreateNextWrite(idx, oldOff+off), nil
}

func IncreaseNextWrite(off int32) error {
	idx, oldOff := GetNextWriteFactory().nextWrite.NextWriteInfo()
	nw, err := increaseNextWrite(off, idx, oldOff)
	if err != nil {
		return err
	}

	GetNextWriteFactory().SetNextWrite(*nw)
	return nil
}

func IncreaseRedoNextWrite(off int32) error {
	idx, oldOff := GetNextWriteFactory().redoLogNextWrite.NextWriteInfo()
	nw, err := increaseNextWrite(off, idx, oldOff)
	if err != nil {
		return err
	}

	GetNextWriteFactory().SetRedoLogNextWrite(*nw)
	return nil
}
