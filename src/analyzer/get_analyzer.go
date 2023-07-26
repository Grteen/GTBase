package analyzer

import (
	"GtBase/pkg/constants"
	"GtBase/src/command"
	"GtBase/src/object"
)

// GET [KEY]
type GetAnalyzer struct {
	parts [][]byte
	cmd   []byte
	cmn   int32
}

func (a *GetAnalyzer) Analyze() Command {
	cmd := CreateGetCommand()
	return a.getKey(0, cmd)
}

func (a *GetAnalyzer) getKey(nowIdx int32, c *GetCommand) Command {
	if len(a.parts) <= int(nowIdx) {
		return CreateErrorArgCommand()
	}
	c.key = object.ParseObjectType(a.parts[nowIdx])
	return c
}

func CreateGetAnalyzer(parts [][]byte, cmd []byte, cmn int32) Analyzer {
	return &GetAnalyzer{parts: parts, cmd: cmd, cmn: cmn}
}

type GetCommand struct {
	key object.Object
}

func (c *GetCommand) Exec() object.Object {
	return c.ExecWithOutRedoLog()
}

func (c *GetCommand) ExecWithOutRedoLog() object.Object {
	result, err := command.Get(c.key)
	if err != nil {
		return object.CreateGtString(err.Error())
	}
	if result == nil {
		return object.CreateGtString(constants.ServerGetNilReturn)
	}

	return result
}

func CreateGetCommand() *GetCommand {
	return &GetCommand{}
}
