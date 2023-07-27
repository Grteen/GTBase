package server

import (
	"GtBase/src/nextwrite"
	"GtBase/src/page"
)

func initFile() {
	page.InitBucketPageFile()
	page.InitPageFile()
	page.InitCheckPointFile()
	page.InitRedoLog()
	nextwrite.InitCMNFile()
	nextwrite.InitNextWrite()
}
