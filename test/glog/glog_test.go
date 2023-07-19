package glog

import (
	"GtBase/pkg/glog"
	"testing"
)

func TestGlogSingleton(t *testing.T) {
	glog1 := glog.Glog()
	glog2 := glog.Glog()
	if glog1 != glog2 {
		t.Errorf("glog.Glog()'s address should be same but got %p and %p", glog1, glog2)
	}
}
