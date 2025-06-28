package main

import (
	"os"

	"github.com/rokeller/aqc/cmd"
)

func main() {
	cmd.Execute(func(err error) { os.Exit(1) })
}
