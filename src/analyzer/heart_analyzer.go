package analyzer

import (
	"GtBase/pkg/constants"
	"GtBase/src/command"
	"GtBase/src/object"
	"GtBase/src/replic"
	"GtBase/utils"
	"log"
)

type HeartAnalyzer struct {
	parts [][]byte
	rs    *replic.ReplicState
}

func (a *HeartAnalyzer) Analyze() Command {
	cmd := CreateHeartCommand(a.rs)
	return a.getHeartSeq(0, cmd)
}

func (a *HeartAnalyzer) getHeartSeq(nowIdx int32, c *HeartCommand) Command {
	if len(a.parts) <= int(nowIdx) {
		return CreateErrorArgCommand()
	}
	c.heartSeq = utils.EncodeBytesSmallEndToint32(a.parts[nowIdx])
	return c
}

func createHeartAnalyzer(rs *replic.ReplicState) Analyzer {
	return &HeartAnalyzer{rs: rs}
}

func CreateHeartAnalyzer(parts [][]byte, cmd []byte, cmn int32, args map[string]interface{}) Analyzer {
	rsItf, ok := args[constants.AssignArgReplicState]
	if !ok {
		return nil
	}
	rs, ok := rsItf.(*replic.ReplicState)
	if !ok {
		return nil
	}

	return createHeartAnalyzer(rs)
}

type HeartCommand struct {
	heartSeq int32
	rs       *replic.ReplicState
}

func (c *HeartCommand) Exec() object.Object {
	return c.ExecWithOutRedoLog()
}

func (c *HeartCommand) ExecWithOutRedoLog() object.Object {
	err := command.Heart(c.heartSeq, c.rs)
	if err != nil {
		log.Println(err)
	}

	return nil
}

func CreateHeartCommand(rs *replic.ReplicState) *HeartCommand {
	return &HeartCommand{rs: rs}
}