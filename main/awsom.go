package main

import (
	"github.com/hekonsek/awsom"
	"github.com/hekonsek/awsom/cmd"
)

func main() {
	awsom.ExitOnCliError(cmd.RootCmd.Execute())
}
