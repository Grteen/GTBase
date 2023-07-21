package utils

import (
	"GtBase/pkg/constants"

	"github.com/spaolacci/murmur3"
)

func FirstHash(target []byte) int32 {
	return int32(murmur3.Sum64(target) % uint64(constants.HashBucketNumber))
}

func SecondHash(target int32) int32 {
	return int32(murmur3.Sum64(Encodeint32ToBytesSmallEnd(target)) % uint64(constants.HashBucketHasBuckets))
}
