package client

import (
	"GtBase/pkg/constants"
	"GtBase/utils"
)

func Heart(client *GtBaseClient) error {
	result := make([]byte, 0)
	result = append(result, []byte(constants.HeartCommand)...)
	result = append(result, []byte(constants.CommandSep)...)

	err := client.Write(result)
	if err != nil {
		return err
	}

	return nil
}

func GetHeart(client *GtBaseClient, logIdx, logOff int32) error {
	result := make([]byte, 0)
	result = append(result, []byte(constants.GetHeartCommand)...)
	result = append(result, []byte(" ")...)
	result = append(result, utils.Encodeint32ToBytesSmallEnd(logIdx)...)
	result = append(result, []byte(" ")...)
	result = append(result, utils.Encodeint32ToBytesSmallEnd(logOff)...)
	result = append(result, []byte(constants.CommandSep)...)
	err := client.Write(result)
	if err != nil {
		return err
	}

	return nil
}

func Redo(client *GtBaseClient, redoLog []byte, seq int32) error {
	result := make([]byte, 0)
	result = append(result, []byte(constants.RedoCommand)...)
	result = append(result, []byte(" ")...)
	result = append(result, utils.Encodeint32ToBytesSmallEnd(seq)...)
	result = append(result, []byte(" ")...)
	result = append(result, redoLog...)
	result = append(result, []byte(constants.CommandSep)...)

	err := client.Write(result)
	if err != nil {
		return err
	}

	return nil
}
