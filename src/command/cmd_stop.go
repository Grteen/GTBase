package command

import "GtBase/src/pair"

type stopFlag int32

type stopFunction func(*pair.Pair) (stopFlag, bool)

const (
	nextIsNil     stopFlag = 1
	nowKeyIsEqual stopFlag = 2
)

func stopWhenNextIsNil(p *pair.Pair) (stopFlag, bool) {
	if p.OverFlow().IsNil() {
		return nextIsNil, true
	}

	return nextIsNil, false
}
