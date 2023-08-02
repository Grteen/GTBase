package client

import (
	"GtBase/pkg/constants"
	"GtBase/utils"
)

func Heart(client *GtBaseClient, heartSeq int32, uuid string) error {
	fileds := make([][]byte, 0)
	fileds = append(fileds, []byte(constants.HeartCommand))
	fileds = append(fileds, utils.Encodeint32ToBytesSmallEnd(heartSeq))
	fileds = append(fileds, []byte(uuid))
	result := utils.EncodeFieldsToGtBasePacket(fileds)

	err := client.Write(result)
	if err != nil {
		return err
	}

	return nil
}

func GetHeart(client *GtBaseClient, logIdx, logOff, heartSeq, seq int32, uuid string) error {
	fileds := make([][]byte, 0)
	fileds = append(fileds, []byte(constants.GetHeartCommand))
	fileds = append(fileds, utils.Encodeint32ToBytesSmallEnd(heartSeq))
	fileds = append(fileds, utils.Encodeint32ToBytesSmallEnd(logIdx))
	fileds = append(fileds, utils.Encodeint32ToBytesSmallEnd(logOff))
	fileds = append(fileds, utils.Encodeint32ToBytesSmallEnd(seq))
	fileds = append(fileds, []byte(uuid))
	result := utils.EncodeFieldsToGtBasePacket(fileds)

	err := client.Write(result)
	if err != nil {
		return err
	}
	return nil
}

func Redo(client *GtBaseClient, redoLog []byte, seq int32, uuid string) error {
	fileds := make([][]byte, 0)
	fileds = append(fileds, []byte(constants.RedoCommand))
	fileds = append(fileds, utils.Encodeint32ToBytesSmallEnd(seq))
	fileds = append(fileds, redoLog)
	fileds = append(fileds, []byte(uuid))
	result := utils.EncodeFieldsToGtBasePacket(fileds)

	err := client.Write(result)
	if err != nil {
		return err
	}

	return nil
}

func GetRedo(client *GtBaseClient, logIdx, logOff, seq int32, uuid string) error {
	fileds := make([][]byte, 0)
	fileds = append(fileds, []byte(constants.GetRedoCommand))
	fileds = append(fileds, utils.Encodeint32ToBytesSmallEnd(logIdx))
	fileds = append(fileds, utils.Encodeint32ToBytesSmallEnd(logOff))
	fileds = append(fileds, utils.Encodeint32ToBytesSmallEnd(seq))
	fileds = append(fileds, []byte(uuid))
	result := utils.EncodeFieldsToGtBasePacket(fileds)

	err := client.Write(result)
	if err != nil {
		return err
	}

	return nil
}

func Slave(client *GtBaseClient, logIdx, logOff, seq int32, host string, port int, uuid string) error {
	fileds := make([][]byte, 0)
	fileds = append(fileds, []byte(constants.SlaveCommand))
	fileds = append(fileds, utils.Encodeint32ToBytesSmallEnd(logIdx))
	fileds = append(fileds, utils.Encodeint32ToBytesSmallEnd(logOff))
	fileds = append(fileds, utils.Encodeint32ToBytesSmallEnd(seq))
	fileds = append(fileds, []byte(host))
	fileds = append(fileds, utils.Encodeint32ToBytesSmallEnd(int32(port)))
	fileds = append(fileds, []byte(uuid))
	result := utils.EncodeFieldsToGtBasePacket(fileds)

	err := client.Write(result)
	if err != nil {
		return err
	}

	return nil
}
