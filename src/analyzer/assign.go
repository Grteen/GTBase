package analyzer

import (
	"GtBase/pkg/constants"
	"GtBase/src/client"
	"GtBase/src/replic"
	"GtBase/utils"
)

type CommandAssign struct {
	bts  []byte
	dict map[string]func([][]byte, []byte, int32, map[string]interface{}) Analyzer
	cmn  int32
	args map[string]interface{}
}

func (c *CommandAssign) InitDict() {
	c.dict = map[string]func([][]byte, []byte, int32, map[string]interface{}) Analyzer{
		constants.SetCommand:      CreateSetAnalyzer,
		constants.GetCommand:      CreateGetAnalyzer,
		constants.DelCommand:      CreateDelAnalyzer,
		constants.SlaveCommand:    CreateSlaveAnalyzer,
		constants.GetRedoCommand:  CreateGetRedoAnalyzer,
		constants.RedoCommand:     CreateRedoAnalyzer,
		constants.GetHeartCommand: CreateGetHeartAnalyzer,
		constants.HeartCommand:    CreateHeartAnalyzer,
		constants.BecomeSlave:     CreateBecomeSlaveAnalyzer,
	}
}

func (c *CommandAssign) Assign() Analyzer {
	fileds := utils.DecodeGtBasePacket(c.bts)
	cmd := fileds[0]

	f, ok := c.dict[string(cmd)]
	if !ok {
		return CreateUnknownCommandAnalyzer(cmd)
	}

	return f(fileds[1:], c.bts, c.cmn, c.args)
}

func CreateCommandAssignArgs(c *client.GtBaseClient, rs *replic.ReplicState, hostSelf string, portSelf int, uuidSelf string) map[string]interface{} {
	return map[string]interface{}{
		constants.AssignArgClient:      c,
		constants.AssignArgReplicState: rs,
		constants.AssignArgHostSelf:    hostSelf,
		constants.AssignArgPortSelf:    portSelf,
		constants.AssignArgUUIDSelf:    uuidSelf,
	}
}

func CreateCommandAssign(bts []byte, cmn int32, args map[string]interface{}) *CommandAssign {
	result := &CommandAssign{bts: bts, cmn: cmn, args: args}
	result.InitDict()
	return result
}
