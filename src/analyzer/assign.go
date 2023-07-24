package analyzer

import (
	"GtBase/pkg/constants"
	"GtBase/utils"
	"bytes"
)

type CommandAssign struct {
	bts []byte
}

func (c *CommandAssign) splitCommand() [][]byte {
	return bytes.Fields(c.bts)
}

func (c *CommandAssign) Assign() Analyzer {
	split := c.splitCommand()
	cmd := split[0]

	if utils.EqualByteSlice(cmd, []byte(constants.GetCommand)) {
		return CreateGetAnalyzer(split[1:])
	} else if utils.EqualByteSlice(cmd, []byte(constants.SetCommand)) {
		return CreateSetAnalyzer(split[1:])
	} else if utils.EqualByteSlice(cmd, []byte(constants.DelCommand)) {
		return CreateDelAnalyzer(split[1:])
	} else {
		return nil
	}
}

func CreateCommandAssign(bts []byte) *CommandAssign {
	return &CommandAssign{bts: bts}
}
