package main

import (
	"github.com/Madredix/clickadu/src/cmd"
)

func main() {
	cmd.RootCmd.Execute() //nolint:errcheck
}
