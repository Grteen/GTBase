package redo

import (
	"GtBase/pkg/constants"
	"GtBase/src/nextwrite"
	"GtBase/src/page"
	"GtBase/utils"
	"log"
	"os"
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
	page.WriteBytesToRedoPageMemory(idx, off, r.ToByte(), r.cmn)
}

func CreateRedo(cmn, cmdlen int32, cmd []byte) *Redo {
	return &Redo{cmn: cmn, cmdLen: cmdlen, cmd: cmd}
}

func ReadRedo(pg *page.RedoPage, off int32) *Redo {
	temp := off

	cmn := pg.ReadCMN(temp)
	temp += constants.RedoLogCMNSize

	cmdLen := pg.ReadCmdLen(temp)
	temp += constants.RedoLogCmdLenSize

	cmd := pg.ReadCmd(temp, cmdLen)

	return CreateRedo(cmn, cmdLen, cmd)
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

func InitRedoLog() {
	if _, err := os.Stat(constants.RedoLogToDo); os.IsNotExist(err) {
		_, errc := os.Create(constants.RedoLogToDo)
		if errc != nil {
			log.Fatalf("InitRedoLog can't create the RedoLog because %s\n", err)
		}

		errm := os.Chmod(constants.RedoLogToDo, 0777)
		if errm != nil {
			log.Fatalf("InitRedoLog can't chmod because of %s\n", errm)
		}
	}
}

func DeleteRedoLog() {
	if _, err := os.Stat(constants.RedoLogToDo); os.IsNotExist(err) {
		return
	}

	errr := os.Remove(constants.RedoLogToDo)
	if errr != nil {
		log.Fatalf("DeleteRedoLog can't remove the RedoLog because %s\n", errr)
	}
}
