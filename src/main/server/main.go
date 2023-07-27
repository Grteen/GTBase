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

	s := server.CreateGtBaseServer()
	err := s.Run(port)
	if err != nil {
		log.Fatalln(err)
	}
}
