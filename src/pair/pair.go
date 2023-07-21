package pair

import (
	"GtBase/src/object"
	"GtBase/src/page"
	"GtBase/utils"
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

// value value-length  key key-length flag overflowIndex overflowOffset
func (p *Pair) ToByte() []byte {
	keyByte := p.key.ToByte()
	valByte := p.value.ToByte()
	ofByte := p.overFlow.ToByte()

	totalLength := len(valByte) + 4 + len(keyByte) + 4 + 1 + len(ofByte)

	result := make([]byte, 0, totalLength)

	result = append(result, valByte...)
	result = append(result, utils.Encodeint32ToBytesSmallEnd(int32(len(valByte)))...)

	result = append(result, keyByte...)
	result = append(result, utils.Encodeint32ToBytesSmallEnd(int32(len(keyByte)))...)

	result = append(result, byte(p.flag))

	result = append(result, ofByte...)

	return result
}

func (p *Pair) WriteInPage(idx, off int32) {
	page.WriteBytesToPageMemory(idx, off, p.ToByte())
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

	result = append(result, utils.Encodeint32ToBytesSmallEnd(of.overFlowIndex)...)
	result = append(result, utils.Encodeint32ToBytesSmallEnd(of.overFlowOffset)...)

	return result
}

func CreateOverFlow(idx, off int32) OverFlow {
	return OverFlow{overFlowIndex: idx, overFlowOffset: off}
}
