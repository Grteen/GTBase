package command

import (
	"GtBase/src/bucket"
	"GtBase/src/object"
)

func Del(key object.Object, cmn int32) error {
	firstIdx, firstOff, err := bucket.FindFirstRecordRLock(key)
	if firstIdx == 0 && firstOff == 0 {
		return nil
	}
	if err != nil {
		return err
	}

	p, loc, errf := FindSameKey(firstIdx, firstOff, key.ToString())
	if errf != nil {
		return errf
	}
	if p == nil {
		return nil
	}

	p.SetCMN(cmn)
	p.Delete()

	p.WriteInPageInMidLock(loc.GetIdx(), loc.GetOff())
	return nil
}
