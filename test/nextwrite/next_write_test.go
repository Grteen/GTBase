package nextwrite

import (
	"GtBase/src/nextwrite"
	"os"
	"testing"
)

func TestInitPageFile(t *testing.T) {
	nextwrite.InitCMNFile()

	if _, err := os.Stat(nextwrite.CMNPathToDo); os.IsNotExist(err) {
		t.Errorf("InitCMNFile() should create the %s but it didn't", nextwrite.CMNPathToDo)
	}
}

func TestFactorySingleton(t *testing.T) {
	fa1 := nextwrite.GetNextWriteFactory()
	fa2 := nextwrite.GetNextWriteFactory()
	if fa1 != fa2 {
		t.Errorf("GetNextWriteFactory's address should be same but got %p and %p", fa1, fa2)
	}
}

func TestReadWriteCMN(t *testing.T) {
	data := []int32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	for _, d := range data {
		result, err := nextwrite.GetCMN()
		if err != nil {
			t.Errorf(err.Error())
		}

		if result != d {
			t.Errorf("GetCMN should get %v but got %v", d, result)
		}
	}
}
