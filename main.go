package main

import (
	"GtBase/src/nextwrite"
	"GtBase/src/page"
	"GtBase/src/server"
	"context"
	"log"
)

func initFile() {
	page.InitBucketPageFile()
	page.InitPageFile()
	page.InitCheckPointFile()
	page.InitRedoLog()
	nextwrite.InitCMNFile()
	nextwrite.InitNextWrite()
}

func main() {
	initFile()

	s := server.CreateGtBaseServer()
	ctx, cancel := context.WithCancel(context.Background())
	go page.FlushDirtyList(ctx)
	err := s.Run(1234)
	if err != nil {
		log.Fatalln(err)
	}

	cancel()
}
