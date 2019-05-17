package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/hekonsek/awsom"
	"github.com/hekonsek/osexit"
	"github.com/spf13/cobra"
)

var envCreateName string

func init() {
	envCreateCommand.Flags().StringVarP(&envCreateName, "name", "", "", "")

	envCommand.AddCommand(envCreateCommand)
}

var envCreateCommand = &cobra.Command{
	Use: "create",
	Run: func(cmd *cobra.Command, args []string) {
		if envCreateName == "" {
			fmt.Println("Environment " + color.GreenString("name") + " cannot be empty.")
			return
		}

		osexit.ExitOnError(awsom.NewVpcBuilder(envCreateName).Create())

		fmt.Println("Environment " + color.GreenString(envCreateName) + " created.")
	},
}
