package redo

import (
	"GtBase/pkg/constants"
	"GtBase/utils"
)

type Redo struct {
	cmn    int32
	cmdLen int32
	cmd    []byte
}

func (r *Redo) GetCMN() int32 {
	return r.cmn
}

func (r *Redo) GetCmdLen() int32 {
	return r.cmdLen
}

func (r *Redo) GetCmd() []byte {
	return r.cmd
}

func (r *Redo) ToByte() []byte {
	result := make([]byte, 0, int(constants.RedoLogCMNSize+constants.RedoLogCmdLenSize+int32(len(r.cmd))))
	result = append(result, utils.Encodeint32ToBytesSmallEnd(r.cmn)...)
	result = append(result, utils.Encodeint32ToBytesSmallEnd(r.cmdLen)...)
	result = append(result, r.cmd...)

	return result
}

// func (r *Redo) WriteInPage() {
// 	page.WriteBytesToRedoPageMemoryLock(idx int32, off int32, bts []byte)
// }
