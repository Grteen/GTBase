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
	uuidSelf string
	rs       *replic.ReplicState
}

func (a *BecomeSlaveAnalyzer) Analyze() Command {
	cmd := CreateBecomeSlaveCommand(a.hostSelf, a.portSelf, a.uuidSelf, a.rs)
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

func createBecomeSlaveAnalyzer(parts [][]byte, hostSelf string, portSelf int, uuidSelf string, rs *replic.ReplicState) Analyzer {
	return &BecomeSlaveAnalyzer{parts: parts, rs: rs, hostSelf: hostSelf, portSelf: portSelf, uuidSelf: uuidSelf}
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

	uuidItf, ok := args[constants.AssignArgUUIDSelf]
	if !ok {
		return nil
	}
	uuid, ok := uuidItf.(string)
	if !ok {
		return nil
	}
	return createBecomeSlaveAnalyzer(parts, host, port, uuid, rs)
}

type BecomeSlaveCommand struct {
	host string
	port int

	hostSelf string
	portSelf int
	uuidSelf string
	rs       *replic.ReplicState
}

func (c *BecomeSlaveCommand) Exec() (object.Object, *utils.Message) {
	return c.ExecWithOutRedoLog()
}

func (c *BecomeSlaveCommand) ExecWithOutRedoLog() (object.Object, *utils.Message) {
	err := command.BecomeSlave(c.host, c.hostSelf, c.port, c.portSelf, c.uuidSelf, c.rs)
	if err != nil {
		return object.CreateGtString(err.Error()), nil
	}

	return object.CreateGtString(constants.ServerOkReturn), nil
}

func CreateBecomeSlaveCommand(hostSelf string, portSelf int, uuidSelf string, rs *replic.ReplicState) *BecomeSlaveCommand {
	return &BecomeSlaveCommand{rs: rs, hostSelf: hostSelf, portSelf: portSelf, uuidSelf: uuidSelf}
}
