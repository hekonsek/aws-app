package main

import (
	"fmt"
	"github.com/hekonsek/awsom/aws"
	"github.com/hekonsek/osexit"
	"github.com/spf13/cobra"
)

func init() {
	awsVpcCommand.AddCommand(awsVpcListCommand)
}

var awsVpcListCommand = &cobra.Command{
	Use: "list",
	Run: func(cmd *cobra.Command, args []string) {
		vpcs, err := aws.ListVpc()
		osexit.ExitOnError(err)

		for _, vpc := range vpcs {
			fmt.Println(vpc)
		}
	},
}
