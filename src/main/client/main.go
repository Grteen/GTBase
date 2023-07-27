package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:1234")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSuffix(text, "\n")

		_, err := conn.Write([]byte(text + "\r\n"))
		if err != nil {
			fmt.Println("Error sending:", err)
			return
		}

		buf := make([]byte, 1024)
		_, err = conn.Read(buf)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("Response:", string(buf))
	}
}
