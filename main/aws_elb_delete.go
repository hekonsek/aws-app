package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/hekonsek/awsom/aws"
	"github.com/hekonsek/osexit"
	"github.com/spf13/cobra"
)

var awsElbDeleteName string

func init() {
	awsElbDeleteCommand.Flags().StringVarP(&awsElbDeleteName, "name", "", "", "")

	awsElbCommand.AddCommand(awsElbDeleteCommand)
}

var awsElbDeleteCommand = &cobra.Command{
	Use: "delete",
	Run: func(cmd *cobra.Command, args []string) {
		if awsElbDeleteName == "" {
			fmt.Println("ELB " + color.GreenString("name") + " cannot be empty.")
			return
		}

		osexit.ExitOnError(aws.DeleteElasticLoadBalancer(awsElbDeleteName))

		fmt.Println("ELB " + color.GreenString(awsElbDeleteName) + " deleted.")
	},
}
