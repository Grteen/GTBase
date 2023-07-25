package pair

import (
	"GtBase/pkg/constants"
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

func (p *Pair) Flag() int8 {
	return p.flag
}

func (p *Pair) Delete() {
	p.flag |= 1
}

func (p *Pair) OverFlow() *OverFlow {
	overFlow := CreateOverFlow(p.overFlow.OverFlowInfo())
	return &overFlow
}

func (p *Pair) SetOverFlow(of OverFlow) {
	p.overFlow = of
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

func (p *Pair) WriteInPageLock(idx, off int32) {
	page.WriteBytesToPageMemoryLock(idx, off, p.ToByte())
}

func (p *Pair) WriteInPageInMid(idx, off int32) {
	page.WriteBytesToPageMemory(idx, off-p.GetMidOffsetNotInBasic(), p.ToByte())
}

func (p *Pair) WriteInPageInMidLock(idx, off int32) {
	page.WriteBytesToPageMemoryLock(idx, off-p.GetMidOffsetNotInBasic(), p.ToByte())
}

// MidOffset points to place between flag and overFlowIndex
func (p *Pair) CalMidOffset(basicOff int32) int32 {
	return basicOff + int32(len(p.value.ToByte())) + constants.PairValLengthSize + int32(len(p.key.ToByte())) +
		constants.PairKeyLengthSize + constants.PairFlagSize
}

func (p *Pair) GetMidOffsetNotInBasic() int32 {
	return int32(len(p.value.ToByte())) + constants.PairValLengthSize + int32(len(p.key.ToByte())) +
		constants.PairKeyLengthSize + constants.PairFlagSize
}

func CreatePair(key, value object.Object, flag int8, of OverFlow) *Pair {
	return &Pair{key: key, value: value, flag: flag, overFlow: of}
}

func ReadPair(pg *page.Page, midOff int32) *Pair {
	temp := midOff

	flag := pg.ReadFlag(temp)
	temp -= 1

	keyLen := pg.ReadKeyLength(temp)
	temp -= 4

	key := pg.ReadKey(temp, keyLen)
	temp -= keyLen

	valLen := pg.ReadValLength(temp)
	temp -= 4

	val := pg.ReadVal(temp, valLen)

	temp = midOff

	overflowIdx, overflowOff := pg.ReadOverFlow(temp)
	return CreatePair(key, val, flag, CreateOverFlow(overflowIdx, overflowOff))
}

func IsDelete(flag int8) bool {
	flag &= 1
	if flag == 1 {
		return true
	}

	return false
}

type OverFlow struct {
	overFlowIndex  int32
	overFlowOffset int32
}

func (of *OverFlow) OverFlowInfo() (int32, int32) {
	return of.overFlowIndex, of.overFlowOffset
}

func (of *OverFlow) IsNil() bool {
	return of.overFlowIndex == 0 && of.overFlowOffset == 0
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

func CreateNullOverFlow() OverFlow {
	return OverFlow{0, 0}
}
