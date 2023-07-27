package main

import (
	"GtBase/src/server"
	"log"
)

func main() {
	s := server.CreateGtBaseServer()
	err := s.Run(1234)
	if err != nil {
		log.Fatalln(err)
	}
}
