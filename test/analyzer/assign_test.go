package analyzer

import (
	"GtBase/pkg/constants"
	"GtBase/src/analyzer"
	"fmt"
	"testing"
)

func TestAssign(t *testing.T) {
	ok := analyzer.CreateCommandAssign([]byte("Set key val")).Assign().Analyze().Exec().ToString()
	if ok != constants.ServerOkReturn {
		t.Errorf("Exec should get %v but got %v", constants.ServerOkReturn, ok)
	}

	res := analyzer.CreateCommandAssign([]byte("Get key")).Assign().Analyze().Exec().ToString()
	if res != "val" {
		t.Errorf("Exec should get %v but got %v", "val", res)
	}

	res = analyzer.CreateCommandAssign([]byte("Get Impossible")).Assign().Analyze().Exec().ToString()
	if res != constants.ServerGetNilReturn {
		t.Errorf("Exec should get %v but got %v", constants.ServerGetNilReturn, res)
	}

	ok = analyzer.CreateCommandAssign([]byte("Del key")).Assign().Analyze().Exec().ToString()
	if ok != constants.ServerOkReturn {
		t.Errorf("Exec should get %v but got %v", constants.ServerOkReturn, ok)
	}

	res = analyzer.CreateCommandAssign([]byte("Get key")).Assign().Analyze().Exec().ToString()
	if res != constants.ServerGetNilReturn {
		t.Errorf("Exec should get %v but got %v", constants.ServerGetNilReturn, res)
	}

	res = analyzer.CreateCommandAssign([]byte("asdasd key")).Assign().Analyze().Exec().ToString()
	result := fmt.Sprintf(constants.ServerUnknownCommandFormat, "asdasd")
	if res != result {
		t.Errorf("Exec should get %v but got %v", result, res)
	}

	res = analyzer.CreateCommandAssign([]byte("Get")).Assign().Analyze().Exec().ToString()
	if res != constants.ServerErrorArg {
		t.Errorf("Exec should get %v but got %v", constants.ServerErrorArg, res)
	}

	res = analyzer.CreateCommandAssign([]byte("Set First")).Assign().Analyze().Exec().ToString()
	if res != constants.ServerErrorArg {
		t.Errorf("Exec should get %v but got %v", constants.ServerErrorArg, res)
	}

	res = analyzer.CreateCommandAssign([]byte("Del")).Assign().Analyze().Exec().ToString()
	if res != constants.ServerErrorArg {
		t.Errorf("Exec should get %v but got %v", constants.ServerErrorArg, res)
	}
}
