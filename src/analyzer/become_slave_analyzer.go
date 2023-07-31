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

	hostSelf string
	portSelf int
	rs       *replic.ReplicState
}

func (a *BecomeSlaveAnalyzer) Analyze() Command {
	cmd := CreateBecomeSlaveCommand(a.hostSelf, a.portSelf, a.rs)
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

func createBecomeSlaveAnalyzer(parts [][]byte, hostSelf string, portSelf int, rs *replic.ReplicState) Analyzer {
	return &BecomeSlaveAnalyzer{parts: parts, rs: rs, hostSelf: hostSelf, portSelf: portSelf}
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

	hostItf, ok := args[constants.AssignArgHostSelf]
	if !ok {
		return nil
	}
	host, ok := hostItf.(string)
	if !ok {
		return nil
	}

	portItf, ok := args[constants.AssignArgPortSelf]
	if !ok {
		return nil
	}
	port, ok := portItf.(int)
	if !ok {
		return nil
	}

	return createBecomeSlaveAnalyzer(parts, host, port, rs)
}

type BecomeSlaveCommand struct {
	host     string
	port     int
	hostSelf string
	portSelf int

	rs *replic.ReplicState
}

func (c *BecomeSlaveCommand) Exec() object.Object {
	return c.ExecWithOutRedoLog()
}

func (c *BecomeSlaveCommand) ExecWithOutRedoLog() object.Object {
	err := command.BecomeSlave(c.host, c.hostSelf, c.port, c.portSelf, c.rs)
	if err != nil {
		return object.CreateGtString(err.Error())
	}

	return object.CreateGtString(constants.ServerOkReturn)
}

func CreateBecomeSlaveCommand(hostSelf string, portSelf int, rs *replic.ReplicState) *BecomeSlaveCommand {
	return &BecomeSlaveCommand{rs: rs, hostSelf: hostSelf, portSelf: portSelf}
}
