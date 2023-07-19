package page

import (
	"GtBase/src/page"
	"testing"
)

func TestPagePoolSingleton(t *testing.T) {
	pool := page.GetPagePool()
	pool2 := page.GetPagePool()
	if pool != pool2 {
		t.Errorf("GetPagePool's address should be same but got %p and %p", pool, pool2)
	}
}
