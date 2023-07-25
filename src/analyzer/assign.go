package analyzer

import (
	"GtBase/pkg/constants"
	"bytes"
)

type CommandAssign struct {
	bts  []byte
	dict map[string]func([][]byte, []byte, int32) Analyzer
	cmn  int32
}

func (c *CommandAssign) InitDict() {
	c.dict = map[string]func([][]byte, []byte, int32) Analyzer{
		constants.SetCommand: CreateSetAnalyzer,
		constants.GetCommand: CreateGetAnalyzer,
		constants.DelCommand: CreateDelAnalyzer,
	}
}

func (c *CommandAssign) splitCommand() [][]byte {
	return bytes.Fields(c.bts)
}

func (c *CommandAssign) Assign() Analyzer {
	split := c.splitCommand()
	cmd := split[0]

	f, ok := c.dict[string(cmd)]
	if !ok {
		return CreateUnknownCommandAnalyzer(cmd)
	}

	return f(split[1:], c.bts, c.cmn)
}

func CreateCommandAssign(bts []byte, cmn int32) *CommandAssign {
	result := &CommandAssign{bts: bts, cmn: cmn}
	result.InitDict()
	return result
}
