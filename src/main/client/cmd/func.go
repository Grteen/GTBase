package cmd

import (
	"GtBase/pkg/constants"
	"GtBase/src/object"
	"GtBase/utils"
	"fmt"
	"os"
	"strconv"
)

func Get(parts []string, c *GtBaseClient) object.Object {
	if len(parts) != 2 {
		return object.CreateGtString(constants.ServerErrorArg)
	}

	req := utils.EncodeFieldsToGtBasePacket(utils.ChangeStringSliceToByteSlic(parts))
	return WriteAndRead(req, c)
}

func Set(parts []string, c *GtBaseClient) object.Object {
	if len(parts) != 3 {
		return object.CreateGtString(constants.ServerErrorArg)
	}

	req := utils.EncodeFieldsToGtBasePacket(utils.ChangeStringSliceToByteSlic(parts))

	return WriteAndRead(req, c)
}

func Del(parts []string, c *GtBaseClient) object.Object {
	if len(parts) != 2 {
		return object.CreateGtString(constants.ServerErrorArg)
	}

	req := utils.EncodeFieldsToGtBasePacket(utils.ChangeStringSliceToByteSlic(parts))

	return WriteAndRead(req, c)
}

func QuitClient(parts []string, c *GtBaseClient) object.Object {
	fmt.Println("bye")
	os.Exit(0)

	return nil
}

func BecomeSlave(parts []string, c *GtBaseClient) object.Object {
	if len(parts) != 3 {
		return object.CreateGtString(constants.ServerErrorArg)
	}

	fields := make([][]byte, 0)
	fields = append(fields, []byte(parts[0]))
	fields = append(fields, []byte(parts[1]))

	p, erra := strconv.Atoi(parts[2])
	if erra != nil {
		return object.CreateGtString(erra.Error())
	}

	fields = append(fields, utils.Encodeint32ToBytesSmallEnd(int32(p)))
	req := utils.EncodeFieldsToGtBasePacket(fields)

	return WriteAndRead(req, c)
}
