package server

import (
	"GtBase/pkg/constants"
	"GtBase/src/command"
	"GtBase/src/object"
)

// GET [KEY]
type GetAnalyzer struct {
	parts [][]byte
}

func (a *GetAnalyzer) Analyze() Command {
	result := CreateGetCommand()
	a.getKey(0, result)
	return result
}

func (a *GetAnalyzer) getKey(nowIdx int32, c *GetCommand) {
	c.key = object.ParseObjectType(a.parts[nowIdx])
}

func CreateGetAnalyzer(parts [][]byte) *GetAnalyzer {
	return &GetAnalyzer{parts: parts}
}

type GetCommand struct {
	key object.Object
}

func (c *GetCommand) Exec() object.Object {
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
