package main

import (
	"GtBase/pkg/constants"
	"GtBase/src/main/client/cmd"
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type GtBaseClient struct {
	host string
	port int
	dict map[string]func()
}

func (c *GtBaseClient) readFromReader(rd *bufio.Reader) (string, error) {
	// ToDo Prompt
	text, err := rd.ReadString('\n')
	if err != nil {
		return "", err
	}

	text = strings.TrimSuffix(text, "\n")
	return text, nil
}

func (c *GtBaseClient) writeToGtBase(conn net.Conn, cmd string) error {
	_, err := conn.Write([]byte(cmd + constants.CommandSep))
	if err != nil {
		return err
	}
	return nil
}

func (c *GtBaseClient) readFromGtBase(conn net.Conn) ([]byte, error) {
	result := make([]byte, 1024)
	n, err := conn.Read(result)
	if err != nil {
		return nil, err
	}

	return result[:n], nil
}

func (c *GtBaseClient) parseCmdToRun(cmd string) bool {
	f, ok := c.dict[cmd]
	if !ok {
		return false
	}

	f()
	return true
}

func CreateGtBaseClient(host string, port int) *GtBaseClient {
	return &GtBaseClient{host: host, port: port, dict: map[string]func(){
		"Quit": cmd.QuitClient,
	}}
}

func (c *GtBaseClient) Run() error {
	conn, err := net.Dial("tcp", c.host+":"+strconv.Itoa(c.port))
	if err != nil {
		return err
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		cmd, err := c.readFromReader(reader)
		if err != nil {
			log.Println(err)
		}

		if c.parseCmdToRun(cmd) {
			continue
		}

		errw := c.writeToGtBase(conn, cmd)
		if errw != nil {
			log.Println(errw)
		}

		result, errr := c.readFromGtBase(conn)
		if errr != nil {
			log.Println(errr)
		}

		fmt.Println(string(result))
	}
}

func main() {
	var host string
	var port int

	flag.StringVar(&host, "h", "127.0.0.1", "host")
	flag.IntVar(&port, "p", 9877, "port")
	flag.Parse()

	c := CreateGtBaseClient(host, port)
	c.Run()
}
