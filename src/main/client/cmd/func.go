package cmd

import (
	"fmt"
	"os"
)

func QuitClient() {
	fmt.Println("bye")
	os.Exit(0)
}
