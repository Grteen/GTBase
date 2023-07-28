package command

import (
	"GtBase/pkg/constants"
	"GtBase/src/client"
	"GtBase/src/replic"
	"errors"
)

func GetHeart(logIdx, logOff, seq int32, client *client.GtBaseClient, rs *replic.ReplicState) error {
	key := client.GenerateKey()
	s, ok := rs.GetSlave(key)
	if !ok {
		return errors.New(constants.ServerSlaveNotExist)
	}

	err := s.GetHeartRespFromSlave(logIdx, logOff, seq)
	if err != nil {
		return err
	}

	return nil
}
