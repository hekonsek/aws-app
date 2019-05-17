package main

import (
	"fmt"
	"github.com/GeertJohan/go.rice"
	"github.com/fatih/color"
	"github.com/hekonsek/awsom"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

func init() {
	appCommand.AddCommand(appConfigureCommand)
}

var appConfigureCommand = &cobra.Command{
	Use: "configure",
	Run: func(cmd *cobra.Command, args []string) {
		box, err := rice.FindBox("../rice")
		awsom.ExitOnCliError(err)

		for _, configFile := range []string{"Dockerfile", "buildspec-version.yml", "buildspec-build.yml", "buildspec-dockerize.yml"} {
			if _, err := os.Stat(configFile); os.IsNotExist(err) {
				mavenDockerfile, err := box.String(configFile)
				awsom.ExitOnCliError(err)

				err = ioutil.WriteFile(configFile, []byte(mavenDockerfile), 0600)
				awsom.ExitOnCliError(err)

				fmt.Println("Generated " + color.GreenString(configFile) + " file.")
			}
		}
	},
}
