package analyzer

import (
	"GtBase/pkg/constants"
	"GtBase/src/client"
	"GtBase/src/command"
	"GtBase/src/object"
	"GtBase/src/replic"
	"GtBase/utils"
)

// Slave [LogIdx] [LogOff] [Seq]
type SlaveAnalyzer struct {
	parts [][]byte

	c  *client.GtBaseClient
	rs *replic.ReplicState
}

func (a *SlaveAnalyzer) Analyze() Command {
	cmd := CreateSlaveCommand(a.c, a.rs)
	return a.getLogIdx(0, cmd)
}

func (a *SlaveAnalyzer) getLogIdx(nowIdx int32, c *SlaveCommand) Command {
	if len(a.parts) <= int(nowIdx) {
		return CreateErrorArgCommand()
	}
	c.logIdx = utils.EncodeBytesSmallEndToint32(a.parts[nowIdx])
	return a.getLogOff(nowIdx+1, c)
}

func (a *SlaveAnalyzer) getLogOff(nowIdx int32, c *SlaveCommand) Command {
	if len(a.parts) <= int(nowIdx) {
		return CreateErrorArgCommand()
	}
	c.logOff = utils.EncodeBytesSmallEndToint32(a.parts[nowIdx])
	return a.getSeq(nowIdx+1, c)
}

func (a *SlaveAnalyzer) getSeq(nowIdx int32, c *SlaveCommand) Command {
	if len(a.parts) <= int(nowIdx) {
		return CreateErrorArgCommand()
	}
	c.seq = utils.EncodeBytesSmallEndToint32(a.parts[nowIdx])
	return c
}

func CreateSlaveAnalyzer(parts [][]byte, c *client.GtBaseClient, rs *replic.ReplicState) Analyzer {
	return &SlaveAnalyzer{parts: parts, c: c, rs: rs}
}

type SlaveCommand struct {
	logIdx int32
	logOff int32
	seq    int32

	c  *client.GtBaseClient
	rs *replic.ReplicState
}

func (c *SlaveCommand) Exec() object.Object {
	command.Slave(c.logIdx, c.logOff, c.seq, c.c, c.rs)
	return object.CreateGtString(constants.ServerOkReturn)
}

func (c *SlaveCommand) ExecWithOutRedoLog() object.Object {
	command.Slave(c.logIdx, c.logOff, c.seq, c.c, c.rs)
	return object.CreateGtString(constants.ServerOkReturn)
}

func CreateSlaveCommand(c *client.GtBaseClient, rs *replic.ReplicState) *SlaveCommand {
	return &SlaveCommand{c: c, rs: rs}
}
