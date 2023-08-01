package analyzer

import (
	"GtBase/pkg/constants"
	"GtBase/src/client"
	"GtBase/src/command"
	"GtBase/src/object"
	"GtBase/src/replic"
	"GtBase/utils"
)

// Redo [LogIdx] [LogOff] [Seq]
type GetRedoAnalyzer struct {
	parts [][]byte

	c  *client.GtBaseClient
	rs *replic.ReplicState
}

func (a *GetRedoAnalyzer) Analyze() Command {
	cmd := CreateGetRedoCommand(a.c, a.rs)
	return a.getLogIdx(0, cmd)
}

func (a *GetRedoAnalyzer) getLogIdx(nowIdx int32, c *GetRedoCommand) Command {
	if len(a.parts) <= int(nowIdx) {
		return CreateErrorArgCommand()
	}
	c.logIdx = utils.EncodeBytesSmallEndToint32(a.parts[nowIdx])
	return a.getLogOff(nowIdx+1, c)
}

func (a *GetRedoAnalyzer) getLogOff(nowIdx int32, c *GetRedoCommand) Command {
	if len(a.parts) <= int(nowIdx) {
		return CreateErrorArgCommand()
	}
	c.logOff = utils.EncodeBytesSmallEndToint32(a.parts[nowIdx])
	return a.getSeq(nowIdx+1, c)
}

func (a *GetRedoAnalyzer) getSeq(nowIdx int32, c *GetRedoCommand) Command {
	if len(a.parts) <= int(nowIdx) {
		return CreateErrorArgCommand()
	}
	c.seq = utils.EncodeBytesSmallEndToint32(a.parts[nowIdx])
	return c
}

func createGetRedoAnalyzer(parts [][]byte, c *client.GtBaseClient, rs *replic.ReplicState) Analyzer {
	return &GetRedoAnalyzer{parts: parts, c: c, rs: rs}
}

func CreateGetRedoAnalyzer(parts [][]byte, cmd []byte, cmn int32, args map[string]interface{}) Analyzer {
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

	return createGetRedoAnalyzer(parts, client, rs)
}

type GetRedoCommand struct {
	logIdx int32
	logOff int32
	seq    int32

	c  *client.GtBaseClient
	rs *replic.ReplicState
}

func (c *GetRedoCommand) Exec() (object.Object, *utils.Message) {
	command.GetRedo(c.logIdx, c.logOff, c.seq, c.c, c.rs)
	return nil, nil
}

func (c *GetRedoCommand) ExecWithOutRedoLog() (object.Object, *utils.Message) {
	command.GetRedo(c.logIdx, c.logOff, c.seq, c.c, c.rs)
	return nil, nil
}

func CreateGetRedoCommand(c *client.GtBaseClient, rs *replic.ReplicState) *GetRedoCommand {
	return &GetRedoCommand{c: c, rs: rs}
}
