/*
Copyright © 2022 Louis Lefebvre <louislefebvre1999@gmail.com>
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/louislef299/aws-sso/pkg/v1/version"
	"github.com/spf13/cobra"
)

var short bool

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"ver", "vers"},
	Short:   "Print the version for aws-sso",
	Long:    `Print the version for aws-sso`,
	Run: func(cmd *cobra.Command, args []string) {
		if short {
			fmt.Println(version.String())
		} else {
			err := version.PrintVersion(os.Stdout, rootCmd)
			if err != nil {
				log.Fatal("couldn't print version:", err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVar(&short, "shorthand", false, "print out just the aws-sso version number")
}
