package command

import (
	"GtBase/pkg/glog"
	"GtBase/src/pair"
)

type stopFlag int32

type stopFunction func(*pair.Pair, []string) (stopFlag, bool, error)

type stopStruct struct {
	f   stopFunction
	arg []string
}

const (
	notTrigger    stopFlag = 0
	nextIsNil     stopFlag = 1
	nowKeyIsEqual stopFlag = 2
)

// arg is nil
func stopWhenNextIsNil(p *pair.Pair, arg []string) (stopFlag, bool, error) {
	if p.OverFlow().IsNil() {
		return nextIsNil, true, nil
	}

	return notTrigger, false, nil
}

// arg[0] is key to compare
func stopWhenKeyEqual(p *pair.Pair, arg []string) (stopFlag, bool, error) {
	if len(arg) != 1 {
		return notTrigger, false, glog.Error("argument's length should be %v but got %v", 1, len(arg))
	}

	if pair.IsDelete(p.Flag()) {
		return notTrigger, false, nil
	}

	if p.Key().ToString() == arg[0] {
		return nowKeyIsEqual, true, nil
	}

	return notTrigger, false, nil
}
