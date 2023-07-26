package analyzer

import (
	"GtBase/pkg/constants"
	"GtBase/src/analyzer"
	"GtBase/src/page"
	"testing"
)

func TestAnalyzer(t *testing.T) {
	page.DeleteBucketPageFile()
	page.InitBucketPageFile()
	data := []struct {
		key string
		val string
	}{
		{"Key", "Val"},
		{"Hello", "World"},
	}

	for _, d := range data {
		cmd := make([][]byte, 0)
		cmd = append(cmd, []byte(d.key))
		cmd = append(cmd, []byte(d.val))

		a := analyzer.CreateSetAnalyzer(cmd, []byte(""), -1)
		res := a.Analyze().Exec().ToString()
		if res != constants.ServerOkReturn {
			t.Errorf("Exec should get %v but got %v", constants.ServerOkReturn, res)
		}
	}

	for _, d := range data {
		cmd := make([][]byte, 0)
		cmd = append(cmd, []byte(d.key))

		a := analyzer.CreateGetAnalyzer(cmd, []byte(""), -1)
		res := a.Analyze().Exec().ToString()
		if res != d.val {
			t.Errorf("Exec should get %v but got %v", d.val, res)
		}
	}

	cmd := make([][]byte, 0)
	cmd = append(cmd, []byte(data[1].key))

	a := analyzer.CreateDelAnalyzer(cmd, []byte(""), -1)
	res := a.Analyze().Exec().ToString()
	if res != constants.ServerOkReturn {
		t.Errorf("Exec should get %v but got %v", constants.ServerOkReturn, res)
	}

	a2 := analyzer.CreateGetAnalyzer(cmd, []byte(""), -1)
	res = a2.Analyze().Exec().ToString()
	if res != constants.ServerGetNilReturn {
		t.Errorf("Exec should get %v but got %v", constants.ServerGetNilReturn, res)
	}
}