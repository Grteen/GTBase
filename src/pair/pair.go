package pair

import "GtBase/src/object"

// Pair is used as the record
type Pair struct {
	key      object.GtString
	value    object.GtString
	flag     int8
	overFlow OverFlow
}

func (p *Pair) Key() object.GtString {
	return p.key
}

func (p *Pair) SetKey(key object.GtString) {
	p.key = key
}

func (p *Pair) Value() object.GtString {
	return p.value
}

func (p *Pair) SetValue(value object.GtString) {
	p.key = value
}

type OverFlow struct {
	overFlowIndex  int32
	overFlowOffset int32
}

func (of *OverFlow) OverFlowInfo() (int32, int32) {
	return of.overFlowIndex, of.overFlowOffset
}

func (of *OverFlow) SetOverFlowInfo(idx, off int32) {
	of.overFlowIndex = idx
	of.overFlowOffset = off
}
