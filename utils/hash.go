package utils

import "github.com/spaolacci/murmur3"

func FirstHash(target []byte) int32 {
	return int32(murmur3.Sum64(target) % 256)
}
