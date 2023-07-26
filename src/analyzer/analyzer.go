package analyzer

import "GtBase/src/object"

type Command interface {
	Exec() object.Object
	ExecWithOutRedoLog() object.Object
}

type Analyzer interface {
	Analyze() Command
}
