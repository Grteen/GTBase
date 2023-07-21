package command

import (
	"GtBase/src/bucket"
	"GtBase/src/nextwrite"
	"GtBase/src/pair"
	"GtBase/utils"
)

func FirstSetInThisBucket(p *pair.Pair) error {
	nw, err := nextwrite.GetNextWriteAndIncreaseIt(int32(len(p.ToByte())))
	if err != nil {
		return err
	}

	idx, off := nw.NextWriteInfo()
	p.WriteInPage(idx, off)

	UpdateBucket(p, idx, off)
	return nil
}

func UpdateBucket(p *pair.Pair, idx, off int32) {
	firstHash := utils.FirstHash(p.Key().ToByte())
	secondHash := utils.SecondHash(firstHash)

	b := bucket.CreateBucket(bucket.CreateBucketHeader(firstHash, secondHash), idx, off)

	b.WriteInPage()
}
