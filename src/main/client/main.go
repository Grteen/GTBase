package main

import (
	"GtBase/src/main/client/cmd"
	"flag"
)

func main() {
	var host string
	var port int

	flag.StringVar(&host, "h", "127.0.0.1", "host")
	flag.IntVar(&port, "p", 9877, "port")
	flag.Parse()

	c := cmd.CreateGtBaseClient(host, port)
	c.Run()
}
