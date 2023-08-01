package command

import (
	"GtBase/pkg/constants"
	"GtBase/src/client"
	"GtBase/src/object"
	"GtBase/src/page"
	"GtBase/src/replic"
	"GtBase/utils"
	"errors"
)

func Slave(logIdx, logOff, seq int32, host string, port int, client *client.GtBaseClient, rs *replic.ReplicState) error {
	s := replic.CreateSlave(logIdx, logOff, seq, nil)
	err := s.InitClient(host, port)
	if err != nil {
		return err
	}
	exist := rs.AppendSlaveLock(s)
	if !exist {
		err := s.SendRedoLogToSlave()
		if err != nil {
			return err
		}
		go s.HeartBeat(rs)
	}

	return nil
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

func BecomeSlave(host, hostSelf string, port, portSelf int, rs *replic.ReplicState) error {
	fd, err := utils.Dial(host, port)
	if err != nil {
		return err
	}

	c := client.CreateGtBaseClient(fd, client.CreateAddress(host, port))
	idx, off, errg := page.GetEndRedoLogIdxAndOff()
	if errg != nil {
		return errg
	}

	rs.SetMasterLock(replic.CreateMaster(idx, off, 0, c))

	errs := client.Slave(c, idx, off, 0, hostSelf, portSelf)
	if errs != nil {
		return errs
	}

	return nil
}
