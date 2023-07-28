package analyzer

import (
	"GtBase/src/client"
	"GtBase/src/replic"
)

// Slave [LogIdx] [LogOff] [Seq]
type SlaveAnalyzer struct {
	parts [][]byte
	cmd   []byte
	cmn   int32
	c     *client.GtBaseClient
	rs    *replic.ReplicState
}
