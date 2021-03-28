package main

import (
	"os"

	"github.com/dikaeinstein/tally/cli"
)

func main() {
	os.Exit(cli.Run("tally", os.Stdout))
}
