package cmd

import (
	"github.com/hekonsek/awsom"
	"github.com/spf13/cobra"
)
import (
	"fmt"
	"os"
)

var rootCmd = &cobra.Command{
	Use:  "awsom",
	Long: "Awsom - toolkit making AWS application delivery simple.",
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			fmt.Println(err)
			os.Exit(awsom.UnixExitCodeGeneralError)
		}
	},
}

func ExecuteRootCmd() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(awsom.UnixExitCodeGeneralError)
	}
}
