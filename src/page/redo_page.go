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

func InitRedoLog() {
	initPageFile(constants.RedoLogToDo)
}

func DeleteRedoLog() {
	deletePageFile(constants.RedoLogToDo)
}
