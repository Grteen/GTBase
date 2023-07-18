package pair

import "GtBase/src/object"

// Pair is used as the record
type Pair struct {
	key      object.Object
	value    object.Object
	flag     int8
	overFlow OverFlow
}

func (p *Pair) Key() object.Object {
	return p.key
}

func (p *Pair) Value() object.Object {
	return p.value
}

func CreatePair(key, value object.Object, flag int8, of OverFlow) *Pair {
	return &Pair{key: key, value: value, flag: flag, overFlow: of}
}

type OverFlow struct {
	overFlowIndex  int32
	overFlowOffset int32
}

func (of *OverFlow) OverFlowInfo() (int32, int32) {
	return of.overFlowIndex, of.overFlowOffset
}

func CreateOverFlow(idx, off int32) OverFlow {
	return OverFlow{overFlowIndex: idx, overFlowOffset: off}
}
