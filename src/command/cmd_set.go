package command

import (
	"GtBase/pkg/glog"
	"GtBase/src/bucket"
	"GtBase/src/nextwrite"
	"GtBase/src/pair"
	"GtBase/utils"
)

// func Set(p *pair.Pair) error {
// 	firstRecordIdx, firstRecordOff, err := bucket.FindFirstRecord(p.Key())
// 	if err != nil {
// 		return err
// 	}

// 	if bucket.IsNilFirstRecord(firstRecordIdx, firstRecordOff) {
// 		return FirstSetInThisBucket(p)
// 	}

// 	p, loc, errf := FindFinalRecord(firstRecordIdx, firstRecordOff)
// 	if errf != nil {
// 		return errf
// 	}

// 	return nil
// }

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
	p, loc, flag, err := TraverseList(firstRecordIdx, firstRecordOff, []stopFunction{stopWhenNextIsNil})
	if err != nil {
		return nil, nil, err
	}

	if flag == nextIsNil {
		return p, loc, nil
	}

	return nil, nil, glog.Error("Flag %v not equal to any condition", flag)
}

// func WriteRecordAndUpdatePrevRecord(newp, prevp *pair.Pair) error {

// }

// func UpdatePrevRecord(prevp *pair.Pair, of *pair.OverFlow) {
// 	prevp.SetOverFlow(*of)
// 	prevp.WriteInPage()
// }
