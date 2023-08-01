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
	ok, _ := analyzer.CreateCommandAssign(cmd[:len(cmd)-2], -1, nil).Assign().Analyze().Exec()
	if ok.ToString() != constants.ServerOkReturn {
		t.Errorf("Exec should get %v but got %v", constants.ServerOkReturn, ok)
	}
	cmd = utils.EncodeFieldsToGtBasePacket([][]byte{[]byte("Get"), []byte("Key")})
	res, _ := analyzer.CreateCommandAssign(cmd[:len(cmd)-2], -1, nil).Assign().Analyze().Exec()
	if res.ToString() != "Val" {
		t.Errorf("Exec should get %v but got %v", "Val", res)
	}
	cmd = utils.EncodeFieldsToGtBasePacket([][]byte{[]byte("Get"), []byte("Impossible")})
	res, _ = analyzer.CreateCommandAssign(cmd[:len(cmd)-2], -1, nil).Assign().Analyze().Exec()
	if res.ToString() != constants.ServerGetNilReturn {
		t.Errorf("Exec should get %v but got %v", constants.ServerGetNilReturn, res)
	}
	cmd = utils.EncodeFieldsToGtBasePacket([][]byte{[]byte("Del"), []byte("Key")})
	ok, _ = analyzer.CreateCommandAssign(cmd[:len(cmd)-2], -1, nil).Assign().Analyze().Exec()
	if ok.ToString() != constants.ServerOkReturn {
		t.Errorf("Exec should get %v but got %v", constants.ServerOkReturn, ok)
	}
	cmd = utils.EncodeFieldsToGtBasePacket([][]byte{[]byte("Get"), []byte("Key")})
	res, _ = analyzer.CreateCommandAssign(cmd[:len(cmd)-2], -1, nil).Assign().Analyze().Exec()
	if res.ToString() != constants.ServerGetNilReturn {
		t.Errorf("Exec should get %v but got %v", constants.ServerGetNilReturn, res)
	}
	cmd = utils.EncodeFieldsToGtBasePacket([][]byte{[]byte("assad"), []byte("Impossible")})
	res, _ = analyzer.CreateCommandAssign(cmd[:len(cmd)-2], -1, nil).Assign().Analyze().Exec()
	result := fmt.Sprintf(constants.ServerUnknownCommandFormat, "assad")
	if res.ToString() != result {
		t.Errorf("Exec should get %v but got %v", result, res)
	}
	cmd = utils.EncodeFieldsToGtBasePacket([][]byte{[]byte("Get")})
	res, _ = analyzer.CreateCommandAssign(cmd[:len(cmd)-2], -1, nil).Assign().Analyze().Exec()
	if res.ToString() != constants.ServerErrorArg {
		t.Errorf("Exec should get %v but got %v", constants.ServerErrorArg, res)
	}
	cmd = utils.EncodeFieldsToGtBasePacket([][]byte{[]byte("Set"), []byte("First")})
	res, _ = analyzer.CreateCommandAssign(cmd[:len(cmd)-2], -1, nil).Assign().Analyze().Exec()
	if res.ToString() != constants.ServerErrorArg {
		t.Errorf("Exec should get %v but got %v", constants.ServerErrorArg, res)
	}
	cmd = utils.EncodeFieldsToGtBasePacket([][]byte{[]byte("Del")})
	res, _ = analyzer.CreateCommandAssign(cmd[:len(cmd)-2], -1, nil).Assign().Analyze().Exec()
	if res.ToString() != constants.ServerErrorArg {
		t.Errorf("Exec should get %v but got %v", constants.ServerErrorArg, res)
	}
}
