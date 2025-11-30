/*
Copyright © 2025 Louis LeFebvre
*/

package cmd

import (
	"log"

	"github.com/louislef299/knot/internal/config"
	"github.com/spf13/cobra"
)

// profileCmd represents the profile command
var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Profile management",
	Run: func(cmd *cobra.Command, args []string) {
		err := config.ListProfiles(cmd.OutOrStdout())
		if err != nil {
			log.Fatal(err)
		}
	},
}

// profileAddCmd represents the profile command
var profileAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Profile management",
	Run: func(cmd *cobra.Command, args []string) {
		// SaveProfile(name string, cfg *ProfileConfig) error
	},
}

func init() {
	rootCmd.AddCommand(profileCmd)
	profileCmd.AddCommand(profileAddCmd)
}
