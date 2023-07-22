package command

import (
	"GtBase/src/bucket"
	"GtBase/src/object"
)

func Del(key object.Object) error {
	firstIdx, firstOff, err := bucket.FindFirstRecord(key)
	if err != nil {
		return err
	}

	p, loc, errf := FindSameKey(firstIdx, firstOff, key.ToString())
	if errf != nil {
		return errf
	}

	p.Delete()

	p.WriteInPageInMid(loc.GetIdx(), loc.GetOff())
	return nil
}
