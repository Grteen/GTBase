package server

import (
	"GtBase/pkg/constants"
	"GtBase/pkg/glog"
	"GtBase/utils"
)

var typeDict map[int]func(interface{}) error = map[int]func(interface{}) error{
	constants.MessageNeedRedo: needRedoMessage,
}

func DoMessage(msg *utils.Message) error {
	f, ok := typeDict[msg.Type]
	if !ok {
		return nil
	}

	go f(msg.Msg)
	return nil
}

func needRedoMessage(msg interface{}) error {
	redoLog, ok := msg.([]byte)
	if !ok {
		return glog.Error("Invalid Type")
	}

	return RedoCmdInReplic(redoLog)
}
