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
		plugs := dlogin.Plugins()
		if len(plugs) == 0 {
			fmt.Println("There aren't any active plugins!")
			return
		}

		fmt.Println("The following plugins are active:")
		for _, p := range plugs {
			fmt.Println(p)
		}
	},
}

func init() {
	rootCmd.AddCommand(pluginsCmd)
}
