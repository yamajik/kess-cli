package main

import (
	"github.com/yamajik/kess/cmd"
)

var (
	version = "0.1.0"
)

func main() {
	cmd.Execute(version)
}
