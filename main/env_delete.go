package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/hekonsek/awsom"
	"github.com/hekonsek/osexit"
	"github.com/spf13/cobra"
)

var envDeleteName string

func init() {
	envDeleteCommand.Flags().StringVarP(&envDeleteName, "name", "", "", "")

	envCommand.AddCommand(envDeleteCommand)
}

var envDeleteCommand = &cobra.Command{
	Use: "delete",
	Run: func(cmd *cobra.Command, args []string) {
		if envDeleteName == "" {
			fmt.Println("Environment " + color.GreenString("name") + " cannot be empty.")
			return
		}

		osexit.ExitOnError(awsom.DeleteVpc(envDeleteName))

		fmt.Println("Environment " + color.GreenString(envDeleteName) + " deleted.")
	},
}
