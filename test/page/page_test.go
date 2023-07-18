package page

import (
	"GtBase/src/page"
	"os"
	"strings"
	"testing"
)

func TestInitPageFile(t *testing.T) {
	page.InitPageFile()

	var filePath strings.Builder
	filePath.WriteString(page.FilePath)
	filePath.WriteString(page.TempFileNameToDo)

	if _, err := os.Stat(filePath.String()); os.IsNotExist(err) {
		if err != nil {
			t.Errorf("InitPageFile() should create the %s but it didn't", filePath.String())
		}
	}
}
