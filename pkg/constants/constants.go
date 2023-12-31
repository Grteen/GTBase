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
	ServerUnknownCommand       string = "Unknown Command"
	ServerErrorArg             string = "Invalid Argument"
	ServerSlaveNotExist        string = "Slave Not Exist"

	IoerRead   int32 = 1
	IoerAccept int32 = 2

	SetCommand      string = "Set"
	GetCommand      string = "Get"
	DelCommand      string = "Del"
	SlaveCommand    string = "Slave"
	RedoCommand     string = "Redo"
	GetRedoCommand  string = "GetRedo"
	HeartCommand    string = "Heart"
	GetHeartCommand string = "GetHeart"
	BecomeSlave     string = "BecomeSlave"

	PagePoolDefaultCapcity int32 = 1024

	ReadNextRedoPageError string = "should read next redo page"
	ClientExitError       string = "client exits"

	CommandSep       string = "\r\n"
	ReplicRedoLogEnd string = "\r\n"

	MaxRedoLogToSendOnceint32       = 100 * PageSize
	SendRedoLogSeqSize        int32 = 4
	SlaveFullSyncThreshold    int64 = PageSize

	SlaveFullSync   int32 = 1
	SlaveSync       int32 = 2
	SlaveDisConnect int32 = 3

	AssignArgClient      string = "AssignClient"
	AssignArgReplicState string = "AssignReplicState"
	AssignArgHostSelf    string = "AssignHostSelf"
	AssignArgPortSelf    string = "AssignPortSelf"
	AssignArgUUIDSelf    string = "AssignUUIDSelf"

	HeartCountLimit int32 = 10
	HeartSeqSize    int32 = 4

	GtBasePacketLengthSize int32 = 4

	MessageNeedRedo int = 1
)
