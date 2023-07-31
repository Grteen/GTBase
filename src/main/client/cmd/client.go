package cmd

import (
	"GtBase/pkg/constants"
	"GtBase/src/object"
	"GtBase/utils"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type GtBaseClient struct {
	fd   int
	host string
	port int
	dict map[string]func([]string, *GtBaseClient) object.Object
}

func (c *GtBaseClient) readFromReader(rd *bufio.Reader) (string, error) {
	fmt.Printf("> ")
	// ToDo Prompt
	text, err := rd.ReadString('\n')
	if err != nil {
		return "", err
	}

	text = strings.TrimSuffix(text, "\n")
	return text, nil
}

func (c *GtBaseClient) writeToGtBase(fd int, cmd []byte) error {
	_, err := utils.WriteFd(fd, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *GtBaseClient) readFromGtBase(fd int) ([]byte, error) {
	result, err := utils.ReadFd(fd)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *GtBaseClient) parseCmdToRun(cmd string) object.Object {
	parts := splitCmdToParts(cmd)
	if len(cmd) < 1 {
		return object.CreateGtString(constants.ServerUnknownCommand)
	}

	f, ok := c.dict[parts[0]]
	if !ok {
		resp := fmt.Sprintf(constants.ServerUnknownCommandFormat, parts[0])
		return object.CreateGtString(resp)
	}

	return f(parts, c)
}

func splitCmdToParts(cmd string) []string {
	return strings.Split(cmd, " ")
}

func CreateGtBaseClient(host string, port int) *GtBaseClient {
	return &GtBaseClient{host: host, port: port, dict: map[string]func([]string, *GtBaseClient) object.Object{
		"Get":         Get,
		"Set":         Set,
		"Del":         Del,
		"Quit":        QuitClient,
		"BecomeSlave": BecomeSlave,
	}}
}

func (c *GtBaseClient) Run() error {
	fd, err := utils.Dial(c.host, c.port)
	if err != nil {
		return err
	}
	defer utils.CloseFd(fd)
	c.fd = fd

	reader := bufio.NewReader(os.Stdin)
	for {
		cmd, err := c.readFromReader(reader)
		if err != nil {
			log.Println(err)
		}

		resp := c.parseCmdToRun(cmd)
		fmt.Println(resp.ToString())
	}
}
