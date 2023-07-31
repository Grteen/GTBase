package main

import (
	"GtBase/src/server"
	"flag"
	"log"
)

func main() {
	var port int

	flag.IntVar(&port, "p", 9877, "port")
	flag.Parse()

	s := server.CreateGtBaseServer("127.0.0.1", port)
	err := s.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
