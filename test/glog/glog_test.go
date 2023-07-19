package glog

import (
	"GtBase/pkg/glog"
	"testing"
)

func TestLog(t *testing.T) {
	glog.Log("this is a test")
	var str string = "good"
	glog.Log("this is a %v test", str)
}
