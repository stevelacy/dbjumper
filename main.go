package main

import (
	"github.com/stevelacy/dbjumper/cli"
)

func main() {
	err := cli.Init()
	if err != nil {
		panic(err)
	}
}
