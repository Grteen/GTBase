package analyzer

import (
	"GtBase/pkg/constants"
	"GtBase/src/client"
	"GtBase/src/command"
	"GtBase/src/object"
	"GtBase/src/replic"
	"GtBase/utils"
)

type GetHeartAnalyzer struct {
	parts [][]byte

	c  *client.GtBaseClient
	rs *replic.ReplicState
}

func (a *GetHeartAnalyzer) Analyze() Command {
	cmd := CreateGetHeartCommand(a.c, a.rs)
	return a.getHeartSeq(0, cmd)
}

func (a *GetHeartAnalyzer) getHeartSeq(nowIdx int32, c *GetHeartCommand) Command {
	if len(a.parts) <= int(nowIdx) {
		return CreateErrorArgCommand()
	}
	c.heartSeq = utils.EncodeBytesSmallEndToint32(a.parts[nowIdx])
	return a.getLogIdx(nowIdx+1, c)
}

func (a *GetHeartAnalyzer) getLogIdx(nowIdx int32, c *GetHeartCommand) Command {
	if len(a.parts) <= int(nowIdx) {
		return CreateErrorArgCommand()
	}
	c.logIdx = utils.EncodeBytesSmallEndToint32(a.parts[nowIdx])
	return a.getLogOff(nowIdx+1, c)
}

func (a *GetHeartAnalyzer) getLogOff(nowIdx int32, c *GetHeartCommand) Command {
	if len(a.parts) <= int(nowIdx) {
		return CreateErrorArgCommand()
	}
	c.logOff = utils.EncodeBytesSmallEndToint32(a.parts[nowIdx])
	return a.getSeq(nowIdx+1, c)
}

func (a *GetHeartAnalyzer) getSeq(nowIdx int32, c *GetHeartCommand) Command {
	if len(a.parts) <= int(nowIdx) {
		return CreateErrorArgCommand()
	}
	c.seq = utils.EncodeBytesSmallEndToint32(a.parts[nowIdx])
	return c
}

func createGetHeartAnalyzer(parts [][]byte, c *client.GtBaseClient, rs *replic.ReplicState) Analyzer {
	return &GetHeartAnalyzer{parts: parts, c: c, rs: rs}
}

func CreateGetHeartAnalyzer(parts [][]byte, cmd []byte, cmn int32, args map[string]interface{}) Analyzer {
	clientItf, ok := args[constants.AssignArgClient]
	if !ok {
		return nil
	}
	client, ok := clientItf.(*client.GtBaseClient)
	if !ok {
		return nil
	}

	rsItf, ok := args[constants.AssignArgReplicState]
	if !ok {
		return nil
	}
	rs, ok := rsItf.(*replic.ReplicState)
	if !ok {
		return nil
	}

	return createGetHeartAnalyzer(parts, client, rs)
}

type GetHeartCommand struct {
	logIdx   int32
	logOff   int32
	seq      int32
	heartSeq int32

	c  *client.GtBaseClient
	rs *replic.ReplicState
}

func (c *GetHeartCommand) Exec() (object.Object, *utils.Message) {
	command.GetHeart(c.logIdx, c.logOff, c.seq, c.heartSeq, c.c, c.rs)
	return nil, nil
}

func (c *GetHeartCommand) ExecWithOutRedoLog() (object.Object, *utils.Message) {
	command.GetHeart(c.logIdx, c.logOff, c.seq, c.heartSeq, c.c, c.rs)
	return nil, nil
}

func CreateGetHeartCommand(c *client.GtBaseClient, rs *replic.ReplicState) *GetHeartCommand {
	return &GetHeartCommand{c: c, rs: rs}
}
