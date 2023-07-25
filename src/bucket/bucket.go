package bucket

import (
	"GtBase/pkg/constants"
	"GtBase/src/object"
	"GtBase/src/page"
	"GtBase/utils"
)

// Bucket is used to store the first record's index and offset in page
type Bucket struct {
	bh          *BucketHeader
	firstIndex  int32
	firstOffset int32
}

func (b *Bucket) FirstIndex() int32 {
	return b.firstIndex
}

func (b *Bucket) FirstOffset() int32 {
	return b.firstOffset
}

func (b *Bucket) BucketHeader() *BucketHeader {
	return b.bh
}

func (b *Bucket) ToByte() []byte {
	result := make([]byte, 0, constants.BucketByteLength)
	result = append(result, utils.Encodeint32ToBytesSmallEnd(b.firstIndex)...)
	result = append(result, utils.Encodeint32ToBytesSmallEnd(b.firstOffset)...)

	return result
}

func (b *Bucket) WriteInPage() {
	page.WriteBytesToBucketPageMemoryLock(b.BucketHeader().CalIndexOfBucketPage(), b.BucketHeader().CalOffsetOfBucketPage(), b.ToByte())
}

func CreateBucket(bh *BucketHeader, firstRecordIdx, firstRecordOff int32) *Bucket {
	return &Bucket{bh, firstRecordIdx, firstRecordOff}
}

func CreateBucketByKey(key object.Object, firstRecordIdx, firstRecordOff int32) *Bucket {
	hashBucketIndex := utils.FirstHash(key.ToByte())
	bucketIndex := utils.SecondHash(hashBucketIndex)

	return CreateBucket(CreateBucketHeader(hashBucketIndex, bucketIndex), firstRecordIdx, firstRecordOff)
}

// find this key's hash's first record
func FindFirstRecordRLock(key object.Object) (int32, int32, error) {
	hashBucketIndex := utils.FirstHash(key.ToByte())
	bucketIndex := utils.SecondHash(hashBucketIndex)

	idx, off, err := findFirstRecordRLock(hashBucketIndex, bucketIndex)
	if err != nil {
		return -1, -1, err
	}

	return idx, off, nil
}

func findFirstRecordRLock(hashBucketIndex, bucketIndex int32) (int32, int32, error) {
	bh := CreateBucketHeader(hashBucketIndex, bucketIndex)

	pg, err := page.ReadBucketPage(bh.CalIndexOfBucketPage())
	if err != nil {
		return -1, -1, err
	}

	pg.RLock()
	defer pg.RUnLock()

	bts := pg.SrcSlice(bh.CalOffsetOfBucketPage(), bh.CalOffsetOfBucketPage()+constants.BucketByteLength)

	idxbts := bts[:constants.BucketByteLength/2]
	offbts := bts[constants.BucketByteLength/2 : constants.BucketByteLength]

	idx := utils.EncodeBytesSmallEndToint32(idxbts)
	off := utils.EncodeBytesSmallEndToint32(offbts)

	return idx, off, nil
}

// FirstRecord index and offset is zero so this bucket doesn't has first record
func IsNilFirstRecord(idx, off int32) bool {
	return idx == 0 && off == 0
}

type BucketHeader struct {
	firstHashValue  int32
	secondHashValue int32
}

func (bh *BucketHeader) FirstHashValue() int32 {
	return bh.firstHashValue
}

func (bh *BucketHeader) SecondHashValue() int32 {
	return bh.secondHashValue
}

func (bh *BucketHeader) CalIndexOfBucketPage() int32 {
	return int32(bh.firstHashValue/constants.PageHasHashBuckets) + 1
}

func (bh *BucketHeader) CalOffsetOfBucketPage() int32 {
	return (bh.firstHashValue%constants.PageHasHashBuckets)*constants.HashBucketSize + bh.secondHashValue*constants.BucketByteLength
}

func CreateBucketHeader(first, second int32) *BucketHeader {
	return &BucketHeader{first, second}
}
