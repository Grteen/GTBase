package server

import "GtBase/src/object"

type Command interface {
	Exec() object.Object
}
