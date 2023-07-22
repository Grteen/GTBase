package command

import (
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

// TraverseList returns the current pair when the stop function returns true
func TraverseList(recordIdx, recordOff int32, stop []stopFunction) (*pair.Pair, *pairLoc, stopFlag, error) {
	pg, err := page.ReadPage(recordIdx)
	if err != nil {
		return nil, nil, 0, err
	}

	p := pair.ReadPair(pg, recordOff)

	for _, s := range stop {
		flag, ok := s(p)
		if ok {
			return p, CreatePairLoc(recordIdx, recordOff), flag, nil
		}
	}

	nextIdx, nextOff := p.OverFlow().OverFlowInfo()

	return TraverseList(nextIdx, nextOff, stop)
}
