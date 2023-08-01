package analyzer

import (
	"GtBase/pkg/constants"
	"GtBase/src/object"
	"GtBase/utils"
	"fmt"
)

type UnKnownCommandAnalyzer struct {
	parts []byte
}

func (a *UnKnownCommandAnalyzer) Analyze() Command {
	return CreateUnknownCommandCommand(a.parts)
}

func CreateUnknownCommandAnalyzer(parts []byte) *UnKnownCommandAnalyzer {
	return &UnKnownCommandAnalyzer{parts: parts}
}

type UnknownCommandCommand struct {
	cmd []byte
}

func (c *UnknownCommandCommand) Exec() (object.Object, *utils.Message) {
	response := fmt.Sprintf(constants.ServerUnknownCommandFormat, string(c.cmd))
	return object.CreateGtString(response), nil
}

func (c *UnknownCommandCommand) ExecWithOutRedoLog() (object.Object, *utils.Message) {
	response := fmt.Sprintf(constants.ServerUnknownCommandFormat, string(c.cmd))
	return object.CreateGtString(response), nil
}

func CreateUnknownCommandCommand(cmd []byte) *UnknownCommandCommand {
	return &UnknownCommandCommand{cmd: cmd}
}

type ErrorArgCommand struct {
	// cmd []byte
}

func (c *ErrorArgCommand) Exec() (object.Object, *utils.Message) {
	return object.CreateGtString(constants.ServerErrorArg), nil
}

func (c *ErrorArgCommand) ExecWithOutRedoLog() (object.Object, *utils.Message) {
	return object.CreateGtString(constants.ServerErrorArg), nil
}

func CreateErrorArgCommand() *ErrorArgCommand {
	return &ErrorArgCommand{}
}
