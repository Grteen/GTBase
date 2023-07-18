package pair

import (
	"GtBase/src/object"
	"encoding/binary"
)

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

func (p *Pair) ToByte() []byte {
	keyByte := p.key.ToByte()
	valByte := p.value.ToByte()
	ofByte := p.overFlow.ToByte()

	result := make([]byte, len(keyByte)+len(valByte)+len(ofByte)+1)
	result = append(result, keyByte...)
	result = append(result, valByte...)
	result = append(result, ofByte...)

	return result
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

func (of *OverFlow) ToByte() []byte {
	result := make([]byte, 0, 8)
	idxByte := make([]byte, 4)
	offByte := make([]byte, 4)

	binary.PutVarint(idxByte, int64(of.overFlowIndex))
	binary.PutVarint(offByte, int64(of.overFlowOffset))
	result = append(result, idxByte...)
	result = append(result, offByte...)

	return result
}

func CreateOverFlow(idx, off int32) OverFlow {
	return OverFlow{overFlowIndex: idx, overFlowOffset: off}
}
