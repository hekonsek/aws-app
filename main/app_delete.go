package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/hekonsek/awsom"
	"github.com/spf13/cobra"
)

var appDeleteName string

func init() {
	appDeleteCommand.Flags().StringVarP(&appDeleteName, "name", "", "", "")

	appCommand.AddCommand(appDeleteCommand)
}

var appDeleteCommand = &cobra.Command{
	Use: "delete",
	Run: func(cmd *cobra.Command, args []string) {
		if appDeleteName == "" {
			fmt.Println("Application " + color.GreenString("name") + " cannot be empty.")
			return
		}

		err := awsom.DeleteApplication(appDeleteName)
		if err != nil {
			fmt.Printf("Something went wrong: %s\n", err.Error())
			return
		}
		fmt.Println("Application " + color.GreenString(appDeleteName) + " deleted.")
	},
}
