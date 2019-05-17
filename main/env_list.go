package main

import (
	"fmt"
	"github.com/hekonsek/awsom"
	"github.com/hekonsek/osexit"
	"github.com/spf13/cobra"
)

func init() {
	envCommand.AddCommand(envListCommand)
}

var envListCommand = &cobra.Command{
	Use: "list",
	Run: func(cmd *cobra.Command, args []string) {
		vpcs, err := awsom.ListVpc()
		osexit.ExitOnError(err)

		for _, vpc := range vpcs {
			fmt.Println(vpc)
		}
	},
}
