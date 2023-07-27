package main

import (
	"GtBase/pkg/constants"
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func readFromReader(rd *bufio.Reader) (string, error) {
	// ToDo Prompt
	text, err := rd.ReadString('\n')
	if err != nil {
		return "", err
	}

	text = strings.TrimSuffix(text, "\n")
	return text, nil
}

func writeToGtBase(conn net.Conn, cmd string) error {
	_, err := conn.Write([]byte(cmd + constants.CommandSep))
	if err != nil {
		return err
	}
	return nil
}

func readFromGtBase(conn net.Conn) ([]byte, error) {
	result := make([]byte, 1024)
	n, err := conn.Read(result)
	if err != nil {
		return nil, err
	}

	return result[:n], nil
}

func main() {
	var host string
	var port int

	flag.StringVar(&host, "h", "127.0.0.1", "host")
	flag.IntVar(&port, "p", 9877, "port")
	flag.Parse()

	conn, err := net.Dial("tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		cmd, err := readFromReader(reader)
		if err != nil {
			log.Println(err)
		}

		errw := writeToGtBase(conn, cmd)
		if errw != nil {
			log.Println(errw)
		}

		result, errr := readFromGtBase(conn)
		if errr != nil {
			log.Println(errr)
		}

		fmt.Println(string(result))
	}
}
