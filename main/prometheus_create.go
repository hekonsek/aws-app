package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/hekonsek/awsom"
	"github.com/hekonsek/osexit"
	"github.com/spf13/cobra"
)

var prometheusCreateName string

func init() {
	prometheusCreateCommand.Flags().StringVarP(&prometheusCreateName, "name", "", "prometheus", "")

	prometheusCommand.AddCommand(prometheusCreateCommand)
}

var prometheusCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Provisions new Prometheus service for monitoring applications and infrastructure.",
	Run: func(cmd *cobra.Command, args []string) {
		if prometheusCreateName == "" {
			fmt.Println("Prometheus service " + color.GreenString("name") + " cannot be empty.")
			return
		}

		osexit.ExitOnError(awsom.NewPrometheusBuilder().WithName(prometheusCreateName).Create())
		fmt.Println("Prometheus service with name " + color.GreenString(prometheusCreateName) + " created.")
	},
}
