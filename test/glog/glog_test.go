package glog

import (
	"GtBase/pkg/glog"
	"errors"
	"testing"
)

func TestLog(t *testing.T) {
	glog.Log("this is a test")
	var str string = "good"
	glog.Log("this is a %v test", str)
}

func TestError(t *testing.T) {
	glog.Log(glog.Error("this is a error").Error())
	err := errors.New("New error")
	glog.Log(glog.Error("error because of %v", err).Error())
}
