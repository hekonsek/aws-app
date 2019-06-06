package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/hekonsek/awsom/aws"
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

		osexit.ExitOnError(aws.DeleteVpc(envDeleteName))

		fmt.Println("Environment " + color.GreenString(envDeleteName) + " deleted.")
	},
}
