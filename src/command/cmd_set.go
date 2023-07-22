package command

import (
	"GtBase/pkg/glog"
	"GtBase/src/bucket"
	"GtBase/src/nextwrite"
	"GtBase/src/pair"
	"GtBase/utils"
)

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

func FindFinalRecord(firstRecordIdx, firstRecordOff int32) (*pair.Pair, error) {
	p, flag, err := TraverseList(firstRecordIdx, firstRecordOff, []stopFunction{stopWhenNextIsNil})
	if err != nil {
		return nil, err
	}

	if flag == nextIsNil {
		return p, nil
	}

	return nil, glog.Error("Flag %v not equal to any condition", flag)
}
