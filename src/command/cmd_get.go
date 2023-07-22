package command

import (
	"GtBase/src/bucket"
	"GtBase/src/object"
	"GtBase/src/page"
	"GtBase/src/pair"
)

type pairLoc struct {
	idx int32
	off int32
}

func (l *pairLoc) GetIdx() int32 {
	return l.idx
}

func (l *pairLoc) GetOff() int32 {
	return l.off
}

func CreatePairLoc(idx, off int32) *pairLoc {
	return &pairLoc{idx, off}
}

func Get(key object.Object) (object.Object, error) {
	firstIdx, firstOff, err := bucket.FindFirstRecord(key)
	if err != nil {
		return nil, err
	}

	p, _, errf := FindSameKey(firstIdx, firstOff, key.ToString())
	if errf != nil {
		return nil, errf
	}

	if p == nil {
		return nil, nil
	}

	return p.Value(), nil
}

func FindSameKey(firstRecordIdx, firstRecordOff int32, key string) (*pair.Pair, *pairLoc, error) {
	p, loc, flag, errt := TraverseList(firstRecordIdx, firstRecordOff, []stopStruct{{stopWhenKeyEqual, []string{key}}})
	if errt != nil {
		return nil, nil, errt
	}

	if flag == nowKeyIsEqual {
		return p, loc, nil
	}

	return nil, nil, nil
}

// TraverseList returns the current pair when the stop function returns true
// second return value it the index and offset of current pair's middle offset
func TraverseList(recordIdx, recordOff int32, stop []stopStruct) (*pair.Pair, *pairLoc, stopFlag, error) {
	pg, err := page.ReadPage(recordIdx)
	if err != nil {
		return nil, nil, 0, err
	}

	p := pair.ReadPair(pg, recordOff)

	for _, s := range stop {
		flag, ok, err := s.f(p, s.arg)
		if err != nil {
			return nil, nil, notTrigger, err
		}
		if ok {
			return p, CreatePairLoc(recordIdx, recordOff), flag, nil
		}
	}

	nextIdx, nextOff := p.OverFlow().OverFlowInfo()
	if nextIdx == 0 && nextOff == 0 {
		return nil, nil, notTrigger, nil
	}

	return TraverseList(nextIdx, nextOff, stop)
}
