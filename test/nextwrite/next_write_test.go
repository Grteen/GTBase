package nextwrite

import (
	"GtBase/src/nextwrite"
	"testing"
)

func TestFactorySingleton(t *testing.T) {
	fa1 := nextwrite.GetNextWriteFactory()
	fa2 := nextwrite.GetNextWriteFactory()
	if fa1 != fa2 {
		t.Errorf("GetNextWriteFactory's address should be same but got %p and %p", fa1, fa2)
	}
}

func TestReadWriteCMN(t *testing.T) {
	data := []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for _, d := range data {

	}
}
