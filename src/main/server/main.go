package main

import (
	"GtBase/src/option"
	"GtBase/src/server"
	"flag"
	"log"
)

func main() {
	var port int
	var storeWay int
	var needRedo int

	flag.IntVar(&port, "p", 9877, "port")
	flag.IntVar(&storeWay, "s", 0, "server is Cache or Persistence 0 is Persistence and 1 is Cache")
	flag.IntVar(&needRedo, "r", 1, "server need redo log or not 0 means don't redo and 1 means redo")
	flag.Parse()

	option.SetGlobalOpt(&option.GlobalOpt{
		StoreWay: storeWay,
		NeedRedo: needRedo,
	})

	s := server.CreateGtBaseServer("127.0.0.1", port)
	err := s.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
