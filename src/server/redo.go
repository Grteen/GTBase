package server

import (
	"GtBase/pkg/constants"
	"GtBase/pkg/glog"
	"GtBase/src/analyzer"
	"GtBase/src/page"
	"GtBase/src/redo"
	"os"
)

func redoCmd(redo *redo.Redo) {
	analyzer.CreateCommandAssign(redo.GetCmd(), redo.GetCMN()).Assign().Analyze().ExecWithOutRedoLog()
}

// if redo all of this page's command it will return a error constants.ReadNextRedoPageError
// if it return nil that means redo is over
func redoCmdInPage(idx, checkPoint int32) error {
	pg, err := page.ReadRedoPage(idx)
	if err != nil {
		return err
	}

	var off int32 = 0

	for off < int32(constants.PageSize) {
		r, errr := redo.ReadRedo(pg, off)
		if errr != nil {
			return errr
		}

		if r == nil {
			break
		}

		if r.GetCMN() >= checkPoint {
			redoCmd(r)
		}

		off += r.GetCmdLen() + constants.RedoLogCMNSize + constants.RedoLogCmdLenSize
	}

	return nil
}

func redoLogTotalLen() (int32, error) {
	fileInfo, err := os.Stat(constants.RedoLogToDo)
	if err != nil {
		return -1, glog.Error("InitNextWrite can't Stat file %v becasuse %v", constants.RedoLogToDo, err)
	}

	fileSize := fileInfo.Size()
	return int32(fileSize) / int32(constants.PageSize), nil
}

func findFirstRedoPageToRedo(checkPoint, totalLen int32) (int32, error) {
	if totalLen == 0 {
		return 0, nil
	}

	l, r := 0, int(totalLen-1)

	for l < r {
		mid := l + (r-l)/2
		cmn, err := getRedoPageFirstLogCMN(int32(mid))
		if err != nil {
			return -1, err
		}

		if cmn >= checkPoint {
			r = mid - 1
		} else {
			l = mid + 1
		}
	}

	return int32(r), nil
}

func getRedoPageFirstLogCMN(idx int32) (int32, error) {
	pg, err := page.ReadRedoPage(idx)
	if err != nil {
		return -1, err
	}

	redo, errr := redo.ReadRedo(pg, 0)
	if err != nil {
		return -1, errr
	}

	return redo.GetCMN(), nil
}

func RedoLog() error {
	checkPoint, err := page.ReadCheckPoint()
	if err != nil {
		return err
	}

	totalLen, errr := redoLogTotalLen()
	if errr != nil {
		return errr
	}

	start, errf := findFirstRedoPageToRedo(checkPoint, totalLen)
	if errf != nil {
		return errf
	}

	for start < totalLen {
		err := redoCmdInPage(start, checkPoint)
		if err != nil {
			if err.Error() == constants.ReadNextRedoPageError {
				start++
				continue
			}

			return err
		}

		return nil
	}

	return nil
}
