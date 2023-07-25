package analyzer

import (
	"GtBase/pkg/constants"
	"bytes"
)

type CommandAssign struct {
	bts  []byte
	dict map[string]func([][]byte) Analyzer
}

func (c *CommandAssign) InitDict() {
	c.dict = map[string]func([][]byte) Analyzer{
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

	return f(split[1:])
}

func CreateCommandAssign(bts []byte) *CommandAssign {
	result := &CommandAssign{bts: bts}
	result.InitDict()
	return result
}
