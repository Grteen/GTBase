package analyzer

import (
	"GtBase/pkg/constants"
	"GtBase/src/command"
	"GtBase/src/object"
	"GtBase/src/replic"
	"log"
)

type HeartAnalyzer struct {
	// parts [][]byte
	rs *replic.ReplicState
}

func (a *HeartAnalyzer) Analyze() Command {
	cmd := CreateHeartCommand(a.rs)
	return cmd
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
	rs *replic.ReplicState
}

func (c *HeartCommand) Exec() object.Object {
	return c.ExecWithOutRedoLog()
}

func (c *HeartCommand) ExecWithOutRedoLog() object.Object {
	err := command.Heart(c.rs)
	if err != nil {
		log.Println(err)
	}

	return nil
}

func CreateHeartCommand(rs *replic.ReplicState) *HeartCommand {
	return &HeartCommand{rs: rs}
}
