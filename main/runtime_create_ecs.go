package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/hekonsek/awsom"
	"github.com/hekonsek/osexit"
	"github.com/spf13/cobra"
)

var runtimeCreateEcsName string

func init() {
	runtimeCreateEcsCommand.Flags().StringVarP(&runtimeCreateEcsName, "name", "", "", "")

	runtimeCreateCommand.AddCommand(runtimeCreateEcsCommand)
}

var runtimeCreateEcsCommand = &cobra.Command{
	Use: "ecs",
	Run: func(cmd *cobra.Command, args []string) {
		if runtimeCreateEcsName == "" {
			fmt.Println("Runtime " + color.GreenString("name") + " cannot be empty.")
			return
		}

		osexit.ExitOnError(awsom.NewEcsClusterBuilder(runtimeCreateEcsName).Create())
		fmt.Println("Runtime " + color.GreenString(runtimeCreateEcsName) + " created.")
	},
}
