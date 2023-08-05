package page

import (
	"GtBase/pkg/constants"
	"GtBase/pkg/glog"
	"GtBase/src/object"
	"GtBase/src/option"
	"GtBase/utils"
	"os"
)

type PairPage struct {
	Page
}

func (p *PairPage) ReadFlag(off int32) int8 {
	return utils.EncodeBytesSmallEndToInt8(p.SrcSliceLength(off-constants.PairFlagSize, constants.PairFlagSize))
}

func (p *PairPage) ReadKeyLength(off int32) int32 {
	return utils.EncodeBytesSmallEndToint32(p.SrcSliceLength(off-constants.PairKeyLengthSize, constants.PairKeyLengthSize))
}

func (p *PairPage) ReadKey(off int32, keyLength int32) object.Object {
	return object.CreateGtStringByBytes(p.SrcSliceLength(off-keyLength, keyLength))
}

func (p *PairPage) ReadValLength(off int32) int32 {
	return utils.EncodeBytesSmallEndToint32(p.SrcSliceLength(off-constants.PairKeyLengthSize, constants.PairKeyLengthSize))
}

func (p *PairPage) ReadVal(off int32, valLength int32) object.Object {
	return object.CreateGtStringByBytes(p.SrcSliceLength(off-valLength, valLength))
}

func (p *PairPage) ReadOverFlow(off int32) (int32, int32) {
	overFlowIdx := utils.EncodeBytesSmallEndToint32(p.SrcSliceLength(off, constants.PairOverFlowIndexSize))
	overFlowOffset := utils.EncodeBytesSmallEndToint32(p.SrcSliceLength(off+constants.PairOverFlowIndexSize, constants.PairOverFlowOffsetSize))
	return overFlowIdx, overFlowOffset
}

func CreatePairPage(idx int32, src []byte, flushPath string) *PairPage {
	ph := CreatePageHeader(idx)
	result := &PairPage{Page: Page{pageHeader: &ph, src: src, flushPath: flushPath}}
	return result
}

func ReadPairPage(idx int32) (*PairPage, error) {
	return readPairPage(idx, constants.PageFilePathToDo)
}

func readPairPage(idx int32, filePath string) (*PairPage, error) {
	p := readPairPageFromCache(idx)
	if p != nil {
		return p, nil
	}

	pd, err := readPairPageFromDisk(idx)
	if err != nil {
		return nil, err
	}

	GetPagePool().CachePairPage(pd)

	return pd, nil
}

func readPairPageFromCache(idx int32) *PairPage {
	p, ok := GetPagePool().GetPairPage(idx)
	if !ok {
		if option.IsCache() {
			GetPagePool().CachePairPage(CreatePairPage(idx, make([]byte, constants.PageSize), ""))
		}
		return nil
	}

	return p
}

func readPairPageFromDisk(idx int32) (*PairPage, error) {
	var pageOffset int64 = CalOffsetOfIndex(idx)
	file, err := os.OpenFile(constants.PageFilePathToDo, os.O_RDWR, 0777)
	if err != nil {
		return nil, glog.Error("ReadPage can't open PageFile because %s\n", err)
	}
	defer file.Close()

	src, err := readOnePageOfBytes(file, pageOffset)
	if err != nil {
		return nil, glog.Error("readOnePageOfBytes can't read because %s\n", err)
	}

	return CreatePairPage(idx, src, constants.PageFilePathToDo), nil
}

func WriteBytesToPairPageMemory(idx, off int32, bts []byte, cmn int32) error {
	pg, err := ReadPairPage(idx)
	if err != nil {
		return err
	}

	pg.SetCMN(cmn)
	pg.WriteBytes(off, bts)

	return nil
}

func WriteBytesToPairPageMemoryLock(idx, off int32, bts []byte, cmn int32) error {
	pg, err := ReadPairPage(idx)
	if err != nil {
		return err
	}

	pg.lock.Lock()
	defer pg.lock.Unlock()
	pg.SetCMN(cmn)
	pg.WriteBytes(off, bts)

	return nil
}

func InitPageFile() {
	initPageFile(constants.PageFilePathToDo)
}

func DeletePageFile() {
	deletePageFile(constants.PageFilePathToDo)
}
