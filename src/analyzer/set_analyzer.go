package analyzer

import (
	"GtBase/pkg/constants"
	"GtBase/src/command"
	"GtBase/src/object"
)

// SET [KEY] [VAL]
type SetAnalyzer struct {
	parts [][]byte
}

func (a *SetAnalyzer) Analyze() Command {
	result := CreateSetCommand()
	a.getKey(0, result)
	return result
}

func (a *SetAnalyzer) getKey(nowIdx int32, c *SetCommand) {
	c.key = object.ParseObjectType(a.parts[nowIdx])
	a.getVal(nowIdx+1, c)
}

func (a *SetAnalyzer) getVal(nowIdx int32, c *SetCommand) {
	c.val = object.ParseObjectType(a.parts[nowIdx])
}

func CreateSetAnalyzer(parts [][]byte) Analyzer {
	return &SetAnalyzer{parts: parts}
}

type SetCommand struct {
	key object.Object
	val object.Object
}

func (c *SetCommand) Exec() object.Object {
	err := command.Set(c.key, c.val)
	if err != nil {
		return object.CreateGtString(err.Error())
	}

	return object.CreateGtString(constants.ServerOkReturn)
}

func CreateSetCommand() *SetCommand {
	return &SetCommand{}
}
