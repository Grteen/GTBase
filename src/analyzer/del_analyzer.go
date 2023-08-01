package analyzer

import (
	"GtBase/pkg/constants"
	"GtBase/src/command"
	"GtBase/src/object"
	"GtBase/src/redo"
	"GtBase/utils"
)

// GET [KEY]
type DelAnalyzer struct {
	parts [][]byte
	cmd   []byte
	cmn   int32
}

func (a *DelAnalyzer) Analyze() Command {
	cmd := CreateDelCommand(a.cmd, a.cmn)
	return a.getKey(0, cmd)
}

func (a *DelAnalyzer) getKey(nowIdx int32, c *DelCommand) Command {
	if len(a.parts) <= int(nowIdx) {
		return CreateErrorArgCommand()
	}
	c.key = object.ParseObjectType(a.parts[nowIdx])
	return c
}

func CreateDelAnalyzer(parts [][]byte, cmd []byte, cmn int32, args map[string]interface{}) Analyzer {
	return &DelAnalyzer{parts: parts, cmd: cmd, cmn: cmn}
}

type DelCommand struct {
	key object.Object
	cmd []byte
	cmn int32
}

func (c *DelCommand) Exec() (object.Object, *utils.Message) {
	redo.WriteRedoLog(c.cmn, c.cmd)

	return c.ExecWithOutRedoLog()
}

func (c *DelCommand) ExecWithOutRedoLog() (object.Object, *utils.Message) {
	err := command.Del(c.key, c.cmn)
	if err != nil {
		return object.CreateGtString(err.Error()), nil
	}

	return object.CreateGtString(constants.ServerOkReturn), nil
}

func CreateDelCommand(cmd []byte, cmn int32) *DelCommand {
	return &DelCommand{cmd: cmd, cmn: cmn}
}
