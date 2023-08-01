package analyzer

import (
	"GtBase/pkg/constants"
	"GtBase/src/command"
	"GtBase/src/object"
	"GtBase/utils"
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

func CreateGetAnalyzer(parts [][]byte, cmd []byte, cmn int32, args map[string]interface{}) Analyzer {
	return &GetAnalyzer{parts: parts, cmd: cmd, cmn: cmn}
}

type GetCommand struct {
	key object.Object
}

func (c *GetCommand) Exec() (object.Object, *utils.Message) {
	return c.ExecWithOutRedoLog()
}

func (c *GetCommand) ExecWithOutRedoLog() (object.Object, *utils.Message) {
	result, err := command.Get(c.key)
	if err != nil {
		return object.CreateGtString(err.Error()), nil
	}
	if result == nil {
		return object.CreateGtString(constants.ServerGetNilReturn), nil
	}

	return result, nil
}

func CreateGetCommand() *GetCommand {
	return &GetCommand{}
}
