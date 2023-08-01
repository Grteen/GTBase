package analyzer

import (
	"GtBase/pkg/constants"
	"GtBase/src/analyzer"
	"GtBase/utils"
	"fmt"
	"testing"
)

func TestAssign(t *testing.T) {
	cmd := utils.EncodeFieldsToGtBasePacket([][]byte{[]byte("Set"), []byte("Key"), []byte("Val")})
	ok := analyzer.CreateCommandAssign(cmd[:len(cmd)-2], -1, nil).Assign().Analyze().Exec().ToString()
	if ok != constants.ServerOkReturn {
		t.Errorf("Exec should get %v but got %v", constants.ServerOkReturn, ok)
	}
	cmd = utils.EncodeFieldsToGtBasePacket([][]byte{[]byte("Get"), []byte("Key")})
	res := analyzer.CreateCommandAssign(cmd[:len(cmd)-2], -1, nil).Assign().Analyze().Exec().ToString()
	if res != "Val" {
		t.Errorf("Exec should get %v but got %v", "Val", res)
	}
	cmd = utils.EncodeFieldsToGtBasePacket([][]byte{[]byte("Get"), []byte("Impossible")})
	res = analyzer.CreateCommandAssign(cmd[:len(cmd)-2], -1, nil).Assign().Analyze().Exec().ToString()
	if res != constants.ServerGetNilReturn {
		t.Errorf("Exec should get %v but got %v", constants.ServerGetNilReturn, res)
	}
	cmd = utils.EncodeFieldsToGtBasePacket([][]byte{[]byte("Del"), []byte("Key")})
	ok = analyzer.CreateCommandAssign(cmd[:len(cmd)-2], -1, nil).Assign().Analyze().Exec().ToString()
	if ok != constants.ServerOkReturn {
		t.Errorf("Exec should get %v but got %v", constants.ServerOkReturn, ok)
	}
	cmd = utils.EncodeFieldsToGtBasePacket([][]byte{[]byte("Get"), []byte("Key")})
	res = analyzer.CreateCommandAssign(cmd[:len(cmd)-2], -1, nil).Assign().Analyze().Exec().ToString()
	if res != constants.ServerGetNilReturn {
		t.Errorf("Exec should get %v but got %v", constants.ServerGetNilReturn, res)
	}
	cmd = utils.EncodeFieldsToGtBasePacket([][]byte{[]byte("assad"), []byte("Impossible")})
	res = analyzer.CreateCommandAssign(cmd[:len(cmd)-2], -1, nil).Assign().Analyze().Exec().ToString()
	result := fmt.Sprintf(constants.ServerUnknownCommandFormat, "assad")
	if res != result {
		t.Errorf("Exec should get %v but got %v", result, res)
	}
	cmd = utils.EncodeFieldsToGtBasePacket([][]byte{[]byte("Get")})
	res = analyzer.CreateCommandAssign(cmd[:len(cmd)-2], -1, nil).Assign().Analyze().Exec().ToString()
	if res != constants.ServerErrorArg {
		t.Errorf("Exec should get %v but got %v", constants.ServerErrorArg, res)
	}
	cmd = utils.EncodeFieldsToGtBasePacket([][]byte{[]byte("Set"), []byte("First")})
	res = analyzer.CreateCommandAssign(cmd[:len(cmd)-2], -1, nil).Assign().Analyze().Exec().ToString()
	if res != constants.ServerErrorArg {
		t.Errorf("Exec should get %v but got %v", constants.ServerErrorArg, res)
	}
	cmd = utils.EncodeFieldsToGtBasePacket([][]byte{[]byte("Del")})
	res = analyzer.CreateCommandAssign(cmd[:len(cmd)-2], -1, nil).Assign().Analyze().Exec().ToString()
	if res != constants.ServerErrorArg {
		t.Errorf("Exec should get %v but got %v", constants.ServerErrorArg, res)
	}
}
