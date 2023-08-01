package server

import (
	"GtBase/pkg/constants"
	"GtBase/utils"
	"time"
)

var typeDict map[int]func() error = map[int]func() error{
	constants.MessageNeedRedo: needRedoMessage,
}

func DoMessage(msg *utils.Message) error {
	f, ok := typeDict[msg.Type]
	if !ok {
		return nil
	}

	go f()
	return nil
}

func needRedoMessage() error {
	old, err := redoLogTotalSize()
	if err != nil {
		return err
	}

	for {
		off, err := redoLogTotalSize()
		if err != nil {
			return err
		}

		if off > old {
			errr := RedoLog()
			if errr != nil {
				return errr
			}
			old = off
		}
		time.Sleep(500 * time.Millisecond)
	}
}
