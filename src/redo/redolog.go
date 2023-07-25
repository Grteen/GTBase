package redo

import (
	"GtBase/pkg/constants"
	"GtBase/src/nextwrite"
	"GtBase/src/page"
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

func (r *Redo) WriteInPage(idx, off int32) {
	page.WriteBytesToRedoPageMemory(idx, off, r.ToByte())
}

func CreateRedo(cmn, cmdlen int32, cmd []byte) *Redo {
	return &Redo{cmn: cmn, cmdLen: cmdlen, cmd: cmd}
}

func WriteRedoLog(cmn int32, cmd []byte) error {
	redo := CreateRedo(cmn, int32(len(cmd)), cmd)
	nw, err := nextwrite.GetRedoNextWriteAndIncreaseIt(int32(len(redo.ToByte())))
	if err != nil {
		return err
	}

	redo.WriteInPage(nw.NextWriteInfo())
	return nil
}
