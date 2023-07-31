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

	req := make([]byte, 0)
	req = append(req, []byte(parts[0])...)
	req = append(req, []byte(" ")...)
	req = append(req, []byte(parts[1])...)
	req = append(req, []byte(constants.CommandSep)...)

	return WriteAndRead(req, c)
}

func Set(parts []string, c *GtBaseClient) object.Object {
	if len(parts) != 3 {
		return object.CreateGtString(constants.ServerErrorArg)
	}

	req := make([]byte, 0)
	req = append(req, []byte(parts[0])...)
	req = append(req, []byte(" ")...)
	req = append(req, []byte(parts[1])...)
	req = append(req, []byte(" ")...)
	req = append(req, []byte(parts[2])...)
	req = append(req, []byte(constants.CommandSep)...)

	return WriteAndRead(req, c)
}

func Del(parts []string, c *GtBaseClient) object.Object {
	if len(parts) != 2 {
		return object.CreateGtString(constants.ServerErrorArg)
	}

	req := make([]byte, 0)
	req = append(req, []byte(parts[0])...)
	req = append(req, []byte(" ")...)
	req = append(req, []byte(parts[1])...)
	req = append(req, []byte(constants.CommandSep)...)

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

	req := make([]byte, 0)
	req = append(req, []byte(parts[0])...)
	req = append(req, []byte(" ")...)
	req = append(req, []byte(parts[1])...)
	req = append(req, []byte(" ")...)
	p, erra := strconv.Atoi(parts[2])
	if erra != nil {
		return object.CreateGtString(erra.Error())
	}
	req = append(req, utils.Encodeint32ToBytesSmallEnd(int32(p))...)
	req = append(req, []byte(constants.CommandSep)...)

	return WriteAndRead(req, c)
}
