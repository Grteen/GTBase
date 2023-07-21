package bucket

import (
	"GtBase/src/page"
	"GtBase/utils"
)

const (
	BucketByteLength     int32 = 8
	HashBucketHasBuckets int32 = 256
	HashBucketSize       int32 = HashBucketHasBuckets * BucketByteLength
	PageHasHashBuckets   int32 = int32(page.PageSize) / HashBucketSize
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
	result := make([]byte, 0, BucketByteLength)
	result = append(result, utils.Encodeint32ToBytesSmallEnd(b.firstIndex)...)
	result = append(result, utils.Encodeint32ToBytesSmallEnd(b.firstOffset)...)

	return result
}

// Bucket.WriteInPage should call index in negative
func (b *Bucket) writeInPage(idx, off int32) {
	page.WriteBytesToPageMemory(idx, off, b.ToByte())
}

func (b *Bucket) WriteInPage() {
	b.writeInPage(-b.BucketHeader().CalIndexOfBucketPage(), b.BucketHeader().CalOffsetOfBucketPage())
}

func CreateBucket(bh *BucketHeader, idx, off int32) *Bucket {
	return &Bucket{bh, idx, off}
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
	return int32(bh.firstHashValue / PageHasHashBuckets)
}

func (bh *BucketHeader) CalOffsetOfBucketPage() int32 {
	return (bh.firstHashValue%PageHasHashBuckets)*HashBucketSize + bh.secondHashValue*BucketByteLength
}

func CreateBucketHeader(first, second int32) *BucketHeader {
	return &BucketHeader{first, second}
}
