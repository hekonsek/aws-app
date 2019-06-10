package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/hekonsek/awsom/aws"
	"github.com/hekonsek/osexit"
	"github.com/spf13/cobra"
)

var awsEcsDeleteName string

func init() {
	awsEcsDeleteCommand.Flags().StringVarP(&awsEcsDeleteName, "name", "", "", "")

	awsEcsCommand.AddCommand(awsEcsDeleteCommand)
}

var awsEcsDeleteCommand = &cobra.Command{
	Use: "delete",
	Run: func(cmd *cobra.Command, args []string) {
		if awsEcsDeleteName == "" {
			fmt.Println("ECS cluster " + color.GreenString("name") + " cannot be empty.")
			return
		}

		osexit.ExitOnError(aws.DeleteEcsCluster(awsEcsDeleteName))

		fmt.Println("ECS cluster " + color.GreenString(awsEcsDeleteName) + " deleted.")
	},
}
