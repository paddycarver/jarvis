package main

import (
	"github.com/paddyforan/jarvis/cli"
	"os"
	"strings"
)

func main() {
	err := cli.Run(strings.Join(os.Args[1:], " "))
	if err != nil {
		panic(err)
	}
}
