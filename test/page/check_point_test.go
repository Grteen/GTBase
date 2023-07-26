package page

import (
	"GtBase/src/page"
	"testing"
)

func TestCheckPoint(t *testing.T) {
	page.DeleteCheckPointFile()
	page.InitCheckPointFile()

	data := []int32{1, 2, 4, 5, 6, 15, 114514}

	for _, d := range data {
		err := page.WriteCheckPoint(d)
		if err != nil {
			t.Errorf(err.Error())
		}

		res, errr := page.ReadCheckPoint()
		if errr != nil {
			t.Errorf(err.Error())
		}

		if res != d {
			t.Errorf("ReadCheckPoint should get %v but got %v", d, res)
		}
	}
}
