package cmd

import (
	"fmt"
	"github.com/hekonsek/awsom"
	"github.com/spf13/cobra"
)

func init() {
	appCommand.AddCommand(appListCommand)
}

var appListCommand = &cobra.Command{
	Use: "list",
	Run: func(cmd *cobra.Command, args []string) {
		applications, err := awsom.ListApplications()
		if err != nil {
			fmt.Printf("Something went wrong: %s\n", err.Error())
			return
		}
		fmt.Printf("Name\n")
		fmt.Printf("----\n")
		for _, application := range applications {
			fmt.Printf("%s\n", application.Name)
		}
	},
}
