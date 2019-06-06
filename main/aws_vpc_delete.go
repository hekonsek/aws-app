package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/hekonsek/awsom/aws"
	"github.com/hekonsek/osexit"
	"github.com/spf13/cobra"
)

var awsVpcDeleteName string

func init() {
	awsVpcDeleteCommand.Flags().StringVarP(&awsVpcDeleteName, "name", "", "", "")

	awsVpcCommand.AddCommand(awsVpcDeleteCommand)
}

var awsVpcDeleteCommand = &cobra.Command{
	Use: "delete",
	Run: func(cmd *cobra.Command, args []string) {
		if awsVpcDeleteName == "" {
			fmt.Println("VPC " + color.GreenString("name") + " cannot be empty.")
			return
		}

		osexit.ExitOnError(aws.DeleteVpc(awsVpcDeleteName))

		fmt.Println("VPC " + color.GreenString(awsVpcDeleteName) + " deleted.")
	},
}
