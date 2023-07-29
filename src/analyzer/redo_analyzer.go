package analyzer

import (
	"GtBase/pkg/constants"
	"GtBase/src/command"
	"GtBase/src/object"
	"GtBase/src/replic"
	"GtBase/utils"
	"log"
)

type RedoAnalyzer struct {
	parts [][]byte

	rs *replic.ReplicState
}

func (a *RedoAnalyzer) Analyze() Command {
	cmd := CreateRedoCommand(a.rs)
	return a.getSeq(0, cmd)
}

func (a *RedoAnalyzer) getSeq(nowIdx int32, c *RedoCommand) Command {
	if len(a.parts) <= int(nowIdx) {
		return CreateErrorArgCommand()
	}
	c.seq = utils.EncodeBytesSmallEndToint32(a.parts[nowIdx])
	return a.getRedoLog(nowIdx+1, c)
}

func (a *RedoAnalyzer) getRedoLog(nowIdx int32, c *RedoCommand) Command {
	if len(a.parts) <= int(nowIdx) {
		return CreateErrorArgCommand()
	}
	c.redoLog = a.parts[nowIdx]
	return c
}

func createRedoAnalyzer(parts [][]byte, rs *replic.ReplicState) Analyzer {
	return &RedoAnalyzer{parts: parts, rs: rs}
}

func CreateRedoAnalyzer(parts [][]byte, cmd []byte, cmn int32, args map[string]interface{}) Analyzer {
	rsItf, ok := args[constants.AssignArgReplicState]
	if !ok {
		return nil
	}
	rs, ok := rsItf.(*replic.ReplicState)
	if !ok {
		return nil
	}

	return createRedoAnalyzer(parts, rs)
}

type RedoCommand struct {
	seq     int32
	redoLog []byte

	rs *replic.ReplicState
}

func (c *RedoCommand) Exec() object.Object {
	return c.ExecWithOutRedoLog()
}

func (c *RedoCommand) ExecWithOutRedoLog() object.Object {
	err := command.Redo(c.seq, c.redoLog, c.rs)
	if err != nil {
		log.Println(err)
	}
	return nil
}

func CreateRedoCommand(rs *replic.ReplicState) *RedoCommand {
	return &RedoCommand{rs: rs}
}
