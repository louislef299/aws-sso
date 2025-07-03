/*
Copyright © 2023 Louis Lefebvre <louislefebvre1999@gmail.com>
*/
package cmd

import (
	"context"
	"log"
	"os"
	"path"
	"time"

	los "github.com/louislef299/aws-sso/pkg/os"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	region, cmdTimeout string
	commandTimeout     time.Duration
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "aws-sso",
	Short: "AWS Auth",
	Long:  `An AWS login helper to make authentication easier`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(ctx context.Context) {
	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	initConfig()
	rootCmd.PersistentFlags().StringVar(&cmdTimeout, "commandTimeout", "3s", "the default timeout for network commands executed")
	var err error
	commandTimeout, err = time.ParseDuration(cmdTimeout)
	if err != nil {
		log.Fatal("could not parse commandTimeout: ", err)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigType("toml")
	if c := os.Getenv("AWS_AUTH_CONFIG"); c != "" {
		viper.AddConfigPath(path.Dir(c))
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".aws-sso" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".aws-sso")
		file := path.Join(home, ".aws-sso")

		if exists, err := los.IsFileOrFolderExisting(file); err != nil {
			panic(err)
		} else if !exists {
			f, err := os.Create(file)
			if err != nil {
				panic(err)
			}
			f.Close()
		}
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}
