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
	cmn      int32
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

func (p *Pair) GetCMN() int32 {
	return p.cmn
}

func (p *Pair) SetCMN(cmn int32) {
	p.cmn = cmn
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
	page.WriteBytesToPairPageMemory(idx, off, p.ToByte(), p.cmn)
}

func (p *Pair) WriteInPageLock(idx, off int32) {
	page.WriteBytesToPairPageMemoryLock(idx, off, p.ToByte(), p.cmn)
}

func (p *Pair) WriteInPageInMid(idx, off int32) {
	page.WriteBytesToPairPageMemory(idx, off-p.GetMidOffsetNotInBasic(), p.ToByte(), p.cmn)
}

func (p *Pair) WriteInPageInMidLock(idx, off int32) {
	page.WriteBytesToPairPageMemoryLock(idx, off-p.GetMidOffsetNotInBasic(), p.ToByte(), p.cmn)
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

func CreatePair(key, value object.Object, flag int8, of OverFlow, cmn int32) *Pair {
	return &Pair{key: key, value: value, flag: flag, overFlow: of, cmn: cmn}
}

func ReadPair(pg *page.PairPage, midOff int32) *Pair {
	temp := midOff

	flag := pg.ReadFlag(temp)
	temp -= constants.PairFlagSize

	keyLen := pg.ReadKeyLength(temp)
	temp -= constants.PairKeyLengthSize

	key := pg.ReadKey(temp, keyLen)
	temp -= keyLen

	valLen := pg.ReadValLength(temp)
	temp -= constants.PairValLengthSize

	val := pg.ReadVal(temp, valLen)

	temp = midOff

	overflowIdx, overflowOff := pg.ReadOverFlow(temp)
	return CreatePair(key, val, flag, CreateOverFlow(overflowIdx, overflowOff), -1)
}

func IsDelete(flag int8) bool {
	flag &= 1
	return flag == 1
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
