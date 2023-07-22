package command

import (
	"GtBase/pkg/glog"
	"GtBase/src/bucket"
	"GtBase/src/nextwrite"
	"GtBase/src/object"
	"GtBase/src/pair"
	"GtBase/utils"
)

func Set(key object.Object, val object.Object) error {
	p := pair.CreatePair(key, val, 0, pair.CreateNullOverFlow())
	return set(p)
}

func set(p *pair.Pair) error {
	firstRecordIdx, firstRecordOff, err := bucket.FindFirstRecord(p.Key())
	if err != nil {
		return err
	}

	if bucket.IsNilFirstRecord(firstRecordIdx, firstRecordOff) {
		return FirstSetInThisBucket(p)
	}

	prevp, prevLoc, errf := FindFinalRecordDelSameKey(firstRecordIdx, firstRecordOff, p.Key().ToString())
	if errf != nil {
		return errf
	}

	errw := WriteRecordAndUpdatePrevRecord(p, prevp, prevLoc)
	if errw != nil {
		return errw
	}

	return nil
}

func FirstSetInThisBucket(p *pair.Pair) error {
	nw, err := nextwrite.GetNextWriteAndIncreaseIt(int32(len(p.ToByte())))
	if err != nil {
		return err
	}

	idx, off := nw.NextWriteInfo()
	p.WriteInPage(idx, off)

	UpdateBucket(p, idx, off)
	return nil
}

func UpdateBucket(p *pair.Pair, idx, off int32) {
	firstHash := utils.FirstHash(p.Key().ToByte())
	secondHash := utils.SecondHash(firstHash)

	b := bucket.CreateBucket(bucket.CreateBucketHeader(firstHash, secondHash), idx, p.CalMidOffset(off))

	b.WriteInPage()
}

func FindFinalRecord(firstRecordIdx, firstRecordOff int32) (*pair.Pair, *pairLoc, error) {
	p, loc, flag, err := TraverseList(firstRecordIdx, firstRecordOff, []stopStruct{{stopWhenNextIsNil, nil}})
	if err != nil {
		return nil, nil, err
	}

	if HasFlag(flag, nextIsNil) {
		return p, loc, nil
	}

	return nil, nil, glog.Error("Flag %v not equal to any condition", flag)
}

func FindFinalRecordDelSameKey(firstRecordIdx, firstRecordOff int32, key string) (*pair.Pair, *pairLoc, error) {
	p, loc, flag, err := TraverseList(firstRecordIdx, firstRecordOff, []stopStruct{{stopWhenNextIsNil, nil}, {stopWhenKeyEqual, []string{key}}})
	if err != nil {
		return nil, nil, err
	}

	if HasFlag(flag, nowKeyIsEqual) {
		p.Delete()
		p.WriteInPageInMid(loc.GetIdx(), loc.GetOff())
		return FindFinalRecordDelSameKey(loc.GetIdx(), loc.GetOff(), key)
	}

	if HasFlag(flag, nextIsNil) {
		return p, loc, nil
	}

	return nil, nil, glog.Error("Flag %v not equal to any condition", flag)
}

func WriteRecordAndUpdatePrevRecord(newp, prevp *pair.Pair, prevLoc *pairLoc) error {
	nw, err := nextwrite.GetNextWriteAndIncreaseIt(int32(len(newp.ToByte())))
	if err != nil {
		return err
	}

	newp.WriteInPage(nw.NextWriteInfo())

	idx, off := nw.NextWriteInfo()

	of := pair.CreateOverFlow(idx, newp.CalMidOffset(off))
	UpdatePrevRecord(prevp, prevLoc, &of)

	return nil
}

func UpdatePrevRecord(prevp *pair.Pair, prevLoc *pairLoc, of *pair.OverFlow) {
	prevp.SetOverFlow(*of)
	prevp.WriteInPageInMid(prevLoc.GetIdx(), prevLoc.GetOff())
}
