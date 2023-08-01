package page

import (
	"GtBase/pkg/constants"
	"GtBase/pkg/glog"
	"GtBase/utils"
	"os"
)

type RedoPage struct {
	Page
}

func (p *RedoPage) ReadCMN(off int32) int32 {
	return utils.EncodeBytesSmallEndToint32(p.SrcSliceLength(off, constants.RedoLogCMNSize))
}

func (p *RedoPage) ReadCmdLen(off int32) int32 {
	return utils.EncodeBytesSmallEndToint32(p.SrcSliceLength(off, constants.RedoLogCmdLenSize))
}

func (p *RedoPage) ReadCmd(off, cmdLen int32) []byte {
	return p.SrcSliceLength(off, cmdLen)
}

func (p *RedoPage) DirtyPageLock() {
	p.pageHeader.mu.Lock()
	defer p.pageHeader.mu.Unlock()
	p.pageHeader.dirty = true
	GetPagePool().RedoDirtyListPush(p, p.GetCMN())
}

func (p *RedoPage) WriteBytes(off int32, bts []byte) {
	// ToDo should ensure the consistency
	for i := 0; i < len(bts); i++ {
		p.src[i+int(off)] = bts[i]
	}
	p.DirtyPageLock()
}

func CreateRedoPage(idx int32, src []byte, flushPath string) *RedoPage {
	ph := CreatePageHeader(idx)
	result := &RedoPage{Page: Page{pageHeader: &ph, src: src, flushPath: flushPath}}
	return result
}

func ReadRedoPage(idx int32) (*RedoPage, error) {
	return readRedoPage(idx, constants.PageFilePathToDo)
}

func readRedoPage(idx int32, filePath string) (*RedoPage, error) {
	p := readRedoPageFromCache(idx)
	if p != nil {
		return p, nil
	}

	pd, err := readRedoPageFromDisk(idx)
	if err != nil {
		return nil, err
	}

	GetPagePool().CacheRedoPage(pd)

	return pd, nil
}

func readRedoPageFromCache(idx int32) *RedoPage {
	p, ok := GetPagePool().GetRedoPage(idx)
	if !ok {
		return nil
	}

	return p
}

func readRedoPageFromDisk(idx int32) (*RedoPage, error) {
	var pageOffset int64 = CalOffsetOfIndex(idx)
	file, err := os.OpenFile(constants.RedoLogToDo, os.O_RDWR, 0777)
	if err != nil {
		return nil, glog.Error("ReadPage can't open PageFile because %s\n", err)
	}
	defer file.Close()

	src, err := readOnePageOfBytes(file, pageOffset)
	if err != nil {
		return nil, glog.Error("readOnePageOfBytes can't read because %s\n", err)
	}

	return CreateRedoPage(idx, src, constants.RedoLogToDo), nil
}

func WriteBytesToRedoPageMemory(idx, off int32, bts []byte, cmn int32) error {
	pg, err := ReadRedoPage(idx)
	if err != nil {
		return err
	}

	pg.SetCMN(cmn)
	pg.WriteBytes(off, bts)

	return nil
}

func WriteBytesToRedoPageMemoryLock(idx, off int32, bts []byte, cmn int32) error {
	pg, err := ReadRedoPage(idx)
	if err != nil {
		return err
	}

	pg.lock.Lock()
	defer pg.lock.Unlock()
	pg.SetCMN(cmn)
	pg.WriteBytes(off, bts)

	return nil
}

func WriteRedoLogFromReplic(idx, off int32, bts []byte) error {
	nowIdx := idx
	nowOff := off
	for len(bts) != 0 {
		pg, err := ReadRedoPage(nowIdx)
		if err != nil {
			return err
		}

		pg.WriteBytes(nowOff, bts[:utils.MinInt(int(constants.PageSize-int64(nowOff)), len(bts))])
		nowOff = 0
		nowIdx++
		if len(bts) < int(constants.PageSize-int64(nowOff)) {
			bts = bts[:0]
		} else {
			bts = bts[constants.PageSize-int64(nowOff):]
		}
	}

	return nil
}

func GetEndRedoLogIdxAndOff() (int32, int32, error) {
	fileInfor, errs := os.Stat(constants.RedoLogToDo)
	if errs != nil {
		return -1, -1, errs
	}

	fileSizer := fileInfor.Size()
	return int32(fileSizer / constants.PageSize), int32(fileSizer % constants.PageSize), nil
}

func InitRedoLog() {
	initPageFile(constants.RedoLogToDo)
}

func DeleteRedoLog() {
	deletePageFile(constants.RedoLogToDo)
}
