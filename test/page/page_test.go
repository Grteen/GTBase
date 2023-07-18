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
	filePath.WriteString(page.PageFilePathToDo)
	filePath.WriteString(page.PageFilePathToDo)

	if _, err := os.Stat(filePath.String()); os.IsNotExist(err) {
		t.Errorf("InitPageFile() should create the %s but it didn't", filePath.String())
	}
}
