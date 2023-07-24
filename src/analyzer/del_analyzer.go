package analyzer

import (
	"GtBase/pkg/constants"
	"GtBase/src/command"
	"GtBase/src/object"
)

// GET [KEY]
type DelAnalyzer struct {
	parts [][]byte
}

func (a *DelAnalyzer) Analyze() Command {
	result := CreateDelCommand()
	a.getKey(0, result)
	return result
}

func (a *DelAnalyzer) getKey(nowIdx int32, c *DelCommand) {
	c.key = object.ParseObjectType(a.parts[nowIdx])
}

func CreateDelAnalyzer(parts [][]byte) *DelAnalyzer {
	return &DelAnalyzer{parts: parts}
}

type DelCommand struct {
	key object.Object
}

func (c *DelCommand) Exec() object.Object {
	err := command.Del(c.key)
	if err != nil {
		return object.CreateGtString(err.Error())
	}

	return object.CreateGtString(constants.ServerOkReturn)
}

func CreateDelCommand() *DelCommand {
	return &DelCommand{}
}
