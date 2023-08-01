package analyzer

import (
	"GtBase/src/object"
	"GtBase/utils"
)

type Command interface {
	Exec() (object.Object, *utils.Message)
	ExecWithOutRedoLog() (object.Object, *utils.Message)
}

type Analyzer interface {
	Analyze() Command
}
