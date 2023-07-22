package command

import (
	"GtBase/src/page"
	"GtBase/src/pair"
)

// TraverseList returns the current pair when the stop function returns true
func TraverseList(recordIdx, recordOff int32, stop []stopFunction) (*pair.Pair, stopFlag, error) {
	pg, err := page.ReadPage(recordIdx)
	if err != nil {
		return nil, 0, err
	}

	p := pair.ReadPair(pg, recordOff)

	for _, s := range stop {
		flag, ok := s(p)
		if ok {
			return p, flag, nil
		}
	}

	nextIdx, nextOff := p.OverFlow().OverFlowInfo()

	return TraverseList(nextIdx, nextOff, stop)
}
