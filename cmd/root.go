package cmd

import "github.com/spf13/cobra"
import (
	"fmt"
	"os"
)

var rootCmd = &cobra.Command{
	Use: "awsom",
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func ExecuteRootCmd() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
