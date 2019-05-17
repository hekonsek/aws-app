package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/hekonsek/awsom"
	"github.com/spf13/cobra"
)

var appCreateName string
var appCreateGitUrl string

func init() {
	appCreateCommand.Flags().StringVarP(&appCreateName, "name", "", "", "")
	appCreateCommand.Flags().StringVarP(&appCreateGitUrl, "git-url", "", "", "")

	appCommand.AddCommand(appCreateCommand)
}

var appCreateCommand = &cobra.Command{
	Use: "create",
	Run: func(cmd *cobra.Command, args []string) {
		if appCreateName == "" {
			fmt.Println("Application " + color.GreenString("name") + " cannot be empty.")
			return
		}
		if appCreateGitUrl == "" {
			fmt.Println("Application " + color.GreenString("git URL") + " cannot be empty.")
			return
		}

		err := (&awsom.Application{
			Name:   appCreateName,
			GitUrl: appCreateGitUrl,
		}).CreateOrUpdate()
		if err != nil {
			if err.Error() == awsom.ErrorApplicationNameTooShort {
				fmt.Printf("Application name %s is less than 3 characters long.\n", color.GreenString(appCreateName))
				return
			} else {
				fmt.Printf("Something went wrong: %s\n", err.Error())
				return
			}
		}
		fmt.Println("Application " + color.GreenString(appCreateName) + " created.")
		println()
		fmt.Printf("See: %s\n", color.GreenString(fmt.Sprintf("https://console.aws.amazon.com/codesuite/codepipeline/pipelines/%s/view", appCreateName)))
	},
}
