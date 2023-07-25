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
	cmd := CreateDelCommand()
	return a.getKey(0, cmd)
}

func (a *DelAnalyzer) getKey(nowIdx int32, c *DelCommand) Command {
	if len(a.parts) <= int(nowIdx) {
		return CreateErrorArgCommand()
	}
	c.key = object.ParseObjectType(a.parts[nowIdx])
	return c
}

func CreateDelAnalyzer(parts [][]byte) Analyzer {
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
