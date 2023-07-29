package command

import (
	"GtBase/pkg/constants"
	"GtBase/src/client"
	"GtBase/src/object"
	"GtBase/src/replic"
	"errors"
)

func Slave(logIdx, logOff, seq int32, client *client.GtBaseClient, rs *replic.ReplicState) {
	s := replic.CreateSlave(logIdx, logOff, seq, client)
	rs.AppendSlaveLock(s)
	s.SendRedoLogToSlave()
}

// if slave is satisfy the FullSyncState then Send Next RedoLog
// otherwise change slave's state to SyncState and don't send next redolog
func GetRedo(logIdx, logOff, seq int32, client *client.GtBaseClient, rs *replic.ReplicState) (object.Object, error) {
	key := client.GenerateKey()
	s, ok := rs.GetSlave(key)
	if !ok {
		return object.CreateGtString(constants.ServerSlaveNotExist), nil
	}

	s.GetSendRedoLogResponseFromSlave(logIdx, logOff, seq)
	state, err := s.CheckFullSyncFinish()
	if err != nil {
		return nil, err
	}

	if state == constants.SlaveFullSync {
		s.SendRedoLogToSlave()
	}

	return object.CreateGtString(constants.ServerOkReturn), nil
}

func Redo(seq int32, redoLog []byte, rs *replic.ReplicState) error {
	return rs.GetMaster().RedoFromMaster(seq, redoLog)
}

func GetHeart(logIdx, logOff, seq, heartSeq int32, client *client.GtBaseClient, rs *replic.ReplicState) error {
	key := client.GenerateKey()
	s, ok := rs.GetSlave(key)
	if !ok {
		return errors.New(constants.ServerSlaveNotExist)
	}

	err := s.GetHeartRespFromSlave(logIdx, logOff, seq, heartSeq)
	if err != nil {
		return err
	}

	return nil
}

func Heart(heartSeq int32, rs *replic.ReplicState) error {
	return rs.GetMaster().HeartFromMaster(heartSeq)
}
