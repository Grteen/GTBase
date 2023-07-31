package analyzer

import (
	"GtBase/pkg/constants"
	"GtBase/src/command"
	"GtBase/src/object"
	"GtBase/src/replic"
	"GtBase/utils"
)

// BecomeSlave [host] [port]
type BecomeSlaveAnalyzer struct {
	parts [][]byte

	rs *replic.ReplicState
}

func (a *BecomeSlaveAnalyzer) Analyze() Command {
	cmd := CreateBecomeCommand(a.rs)
	return a.getHost(0, cmd)
}

func (a *BecomeSlaveAnalyzer) getHost(nowIdx int32, c *BecomeSlaveCommand) Command {
	if len(a.parts) <= int(nowIdx) {
		return CreateErrorArgCommand()
	}
	c.host = string(a.parts[nowIdx])
	return a.getPort(nowIdx+1, c)
}

func (a *BecomeSlaveAnalyzer) getPort(nowIdx int32, c *BecomeSlaveCommand) Command {
	if len(a.parts) <= int(nowIdx) {
		return CreateErrorArgCommand()
	}
	c.port = int(utils.EncodeBytesSmallEndToint32(a.parts[nowIdx]))
	return c
}

func createBecomeSlaveAnalyzer(parts [][]byte, rs *replic.ReplicState) Analyzer {
	return &BecomeSlaveAnalyzer{parts: parts, rs: rs}
}

func CreateBecomeSlaveAnalyzer(parts [][]byte, cmd []byte, cmn int32, args map[string]interface{}) Analyzer {
	rsItf, ok := args[constants.AssignArgReplicState]
	if !ok {
		return nil
	}
	rs, ok := rsItf.(*replic.ReplicState)
	if !ok {
		return nil
	}

	return createBecomeSlaveAnalyzer(parts, rs)
}

type BecomeSlaveCommand struct {
	host string
	port int

	rs *replic.ReplicState
}

func (c *BecomeSlaveCommand) Exec() object.Object {
	return c.ExecWithOutRedoLog()
}

func (c *BecomeSlaveCommand) ExecWithOutRedoLog() object.Object {
	err := command.BecomeSlave(c.host, c.port, c.rs)
	if err != nil {
		return object.CreateGtString(err.Error())
	}

	return object.CreateGtString(constants.ServerOkReturn)
}

func CreateBecomeCommand(rs *replic.ReplicState) *BecomeSlaveCommand {
	return &BecomeSlaveCommand{rs: rs}
}
