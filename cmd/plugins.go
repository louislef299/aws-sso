package cmd

import (
	"fmt"

	"github.com/louislef299/aws-sso/pkg/dlogin"
	"github.com/spf13/cobra"
)

var pluginsCmd = &cobra.Command{
	Use:   "plugins",
	Short: "List active plugins.",
	Run: func(cmd *cobra.Command, args []string) {
		drivers := dlogin.Drivers()
		if len(drivers) == 0 {
			fmt.Println("There aren't any active plugins!")
			return
		}

		fmt.Println("The following plugins are active:")
		for _, d := range drivers {
			fmt.Println(d)
		}
	},
}

func init() {
	rootCmd.AddCommand(pluginsCmd)
}
