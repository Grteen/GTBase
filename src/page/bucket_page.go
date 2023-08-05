package page

import (
	"GtBase/pkg/constants"
	"GtBase/pkg/glog"
	"GtBase/src/option"
	"os"
)

type BucketPage struct {
	Page
}

func CreateBucketPage(idx int32, src []byte, flushPath string) *BucketPage {
	ph := CreatePageHeader(idx)
	result := &BucketPage{Page: Page{pageHeader: &ph, src: src, flushPath: flushPath}}
	return result
}

func ReadBucketPage(idx int32) (*BucketPage, error) {
	return readBucketPage(idx, constants.PageFilePathToDo)
}

func readBucketPage(idx int32, filePath string) (*BucketPage, error) {
	p := readBucketPageFromCache(idx)
	if p != nil {
		return p, nil
	}

	pd, err := readBucketPageFromDisk(idx)
	if err != nil {
		return nil, err
	}

	GetPagePool().CacheBucketPage(pd)

	return pd, nil
}

func readBucketPageFromCache(idx int32) *BucketPage {
	p, ok := GetPagePool().GetBucketPage(idx)
	if !ok {
		if option.IsCache() {
			GetPagePool().CacheBucketPage(CreateBucketPage(idx, make([]byte, constants.PageSize), ""))
		}
		return nil
	}

	return p
}

func readBucketPageFromDisk(idx int32) (*BucketPage, error) {
	var pageOffset int64 = CalOffsetOfIndex(idx)
	file, err := os.OpenFile(constants.BucketPageFilePathToDo, os.O_RDWR, 0777)
	if err != nil {
		return nil, glog.Error("ReadPage can't open PageFile because %s\n", err)
	}
	defer file.Close()

	src, err := readOnePageOfBytes(file, pageOffset)
	if err != nil {
		return nil, glog.Error("readOnePageOfBytes can't read because %s\n", err)
	}

	return CreateBucketPage(idx, src, constants.BucketPageFilePathToDo), nil
}

func WriteBytesToBucketrPageMemory(idx, off int32, bts []byte, cmn int32) error {
	pg, err := ReadBucketPage(idx)
	if err != nil {
		return err
	}

	pg.SetCMN(cmn)
	pg.WriteBytes(off, bts)

	return nil
}

func WriteBytesToBucketPageMemoryLock(idx, off int32, bts []byte, cmn int32) error {
	pg, err := ReadBucketPage(idx)
	if err != nil {
		return err
	}

	pg.lock.Lock()
	defer pg.lock.Unlock()
	pg.SetCMN(cmn)
	pg.WriteBytes(off, bts)

	return nil
}

func InitBucketPageFile() {
	initPageFile(constants.BucketPageFilePathToDo)
}

func DeleteBucketPageFile() {
	deletePageFile(constants.BucketPageFilePathToDo)
}
