/*
Copyright © 2025 Louis LeFebvre
*/

package cmd

import (
	"log"

	"github.com/louislef299/knot/internal/config"
	"github.com/spf13/cobra"
)

// profileCmd represents the profiles command
var profileCmd = &cobra.Command{
	Use:   "profiles",
	Short: "List the profiles currently configured",
	Run: func(cmd *cobra.Command, args []string) {
		err := config.ListProfiles(cmd.OutOrStdout())
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(profileCmd)
}
