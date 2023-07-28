package command

import (
	"GtBase/pkg/constants"
	"GtBase/src/client"
	"GtBase/src/object"
	"GtBase/src/replic"
)

func Slave(logIdx, logOff, seq int32, client *client.GtBaseClient, rs *replic.ReplicState) {
	s := replic.CreateSlave(logIdx, logOff, seq, client)
	rs.AppendSlaveLock(s)
	s.SendRedoLog()
}

func GetRedo(logIdx, logOff, seq int32, client *client.GtBaseClient, rs *replic.ReplicState) object.Object {
	key := client.GenerateKey()
	s, ok := rs.GetSlave(key)
	if !ok {
		return object.CreateGtString(constants.ServerSlaveNotExist)
	}

	s.GetSendRedoLogResponse(logIdx, logOff, seq)
	return object.CreateGtString(constants.ServerOkReturn)
}
