package analyzer

import (
	"GtBase/pkg/constants"
	"GtBase/src/command"
	"GtBase/src/object"
	"GtBase/src/redo"
)

// SET [KEY] [VAL]
type SetAnalyzer struct {
	parts [][]byte
	cmd   []byte
	cmn   int32
}

func (a *SetAnalyzer) Analyze() Command {
	cmd := CreateSetCommand(a.cmd, a.cmn)
	return a.getKey(0, cmd)
}

func (a *SetAnalyzer) getKey(nowIdx int32, c *SetCommand) Command {
	if len(a.parts) <= int(nowIdx) {
		return CreateErrorArgCommand()
	}
	c.key = object.ParseObjectType(a.parts[nowIdx])
	return a.getVal(nowIdx+1, c)
}

func (a *SetAnalyzer) getVal(nowIdx int32, c *SetCommand) Command {
	if len(a.parts) <= int(nowIdx) {
		return CreateErrorArgCommand()
	}
	c.val = object.ParseObjectType(a.parts[nowIdx])
	return c
}

func CreateSetAnalyzer(parts [][]byte, cmd []byte, cmn int32, args map[string]interface{}) Analyzer {
	return &SetAnalyzer{parts: parts, cmd: cmd, cmn: cmn}
}

type SetCommand struct {
	key object.Object
	val object.Object
	cmd []byte
	cmn int32
}

func (c *SetCommand) Exec() object.Object {
	redo.WriteRedoLog(c.cmn, c.cmd)

	return c.ExecWithOutRedoLog()
}

func (c *SetCommand) ExecWithOutRedoLog() object.Object {
	err := command.Set(c.key, c.val, c.cmn)
	if err != nil {
		return object.CreateGtString(err.Error())
	}

	return object.CreateGtString(constants.ServerOkReturn)
}

func CreateSetCommand(cmd []byte, cmn int32) *SetCommand {
	return &SetCommand{cmd: cmd, cmn: cmn}
}
