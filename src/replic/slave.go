package replic

import (
	"GtBase/src/client"
	"sync"
)

type NextSeq struct {
	seq   int32
	count int32
	cLock sync.Mutex
}

type Slave struct {
	client.GtBaseClient
	logIdx  int32
	logOff  int32
	nextSeq *NextSeq
}

// func CreateSlave(logIdx, logOff int32) *Slave {

// }
