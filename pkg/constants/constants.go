package constants

const (
	PageFilePathToDo       string = "/root/GtBase/temp/gt.pf"
	BucketPageFilePathToDo string = "/root/GtBase/temp/gt.bf"
	RedoLogToDo            string = "/root/GtBase/temp/redo.log"
	CMNPathToDo            string = "/root/GtBase/temp/gt.cmn"
	CheckPointPathToDo     string = "/root/GtBase/temp/gt.cp"
	PageSize               int64  = 16384

	BucketByteLength     int32 = 8
	HashBucketHasBuckets int32 = 256
	HashBucketSize       int32 = HashBucketHasBuckets * BucketByteLength
	PageHasHashBuckets   int32 = int32(PageSize) / HashBucketSize

	HashBucketNumber int32 = 256

	PairFlagSize           int32 = 1
	PairKeyLengthSize      int32 = 4
	PairValLengthSize      int32 = 4
	PairOverFlowIndexSize  int32 = 4
	PairOverFlowOffsetSize int32 = 4

	RedoLogCMNSize    int32 = 4
	RedoLogCmdLenSize int32 = 4

	ServerGetNilReturn         string = "Nil"
	ServerOkReturn             string = "Ok"
	ServerUnknownCommandFormat string = "Unknown Command %v"
	ServerErrorArg             string = "Invalid Argument"

	IoerRead   int32 = 1
	IoerAccept int32 = 2

	SetCommand string = "Set"
	GetCommand string = "Get"
	DelCommand string = "Del"

	PagePoolDefaultCapcity = 1024
)
