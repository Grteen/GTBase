package bucket

import "GtBase/utils"

const (
	BucketByteLength int32 = 8
)

// Bucket is used to store the first record's index and offset in page
type Bucket struct {
	firstIndex  int32
	firstOffset int32
}

func (b *Bucket) FirstIndex() int32 {
	return b.firstIndex
}

func (b *Bucket) FirstOffset() int32 {
	return b.firstOffset
}

func (b *Bucket) ToByte() []byte {
	result := make([]byte, 0, BucketByteLength)
	result = append(result, utils.Encodeint32ToBytesSmallEnd(b.firstIndex)...)
	result = append(result, utils.Encodeint32ToBytesSmallEnd(b.firstOffset)...)

	return result
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
