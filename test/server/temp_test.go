package server

import (
	"GtBase/src/page"
	"fmt"
	"testing"
)

func TestFind(t *testing.T) {
	pg, err := page.ReadPairPage(1)
	if err != nil {
		t.Errorf(err.Error())
	}

	fmt.Println(pg.SrcSliceLength(0, 100))
}
