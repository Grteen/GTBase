package constants

const (
	PageFilePathToDo       string = "E:/Code/GTCDN/GTbase/temp/gt.pf"
	BucketPageFilePathToDo string = "E:/Code/GTCDN/GTbase/temp/gt.bf"
	PageSize               int64  = 16384

	BucketByteLength     int32 = 8
	HashBucketHasBuckets int32 = 256
	HashBucketSize       int32 = HashBucketHasBuckets * BucketByteLength
	PageHasHashBuckets   int32 = int32(PageSize) / HashBucketSize
)
