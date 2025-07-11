package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/louislef299/aws-sso/internal/envs"
	"github.com/louislef299/aws-sso/pkg/dlogin"
	los "github.com/louislef299/aws-sso/pkg/os"
	_ "github.com/louislef299/aws-sso/plugins/aws/ecr"
	_ "github.com/louislef299/aws-sso/plugins/aws/eks"
	_ "github.com/louislef299/aws-sso/plugins/aws/oidc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	region, cmdTimeout string
	commandTimeout     time.Duration
)

const (
	AO_CONFIG_NAME = ".aws-sso"
	AO_ENV_PREFIX  = "AWS_SSO"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "aws-sso",
	Short: "AWS Auth",
	Long:  `An AWS login helper to make authentication easier`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

		// Force parse flags manually before using viper-bound values
		if err := cmd.Flags().Parse(os.Args[1:]); err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(cmd.Context(), commandTimeout)
		go func() {
			<-ctx.Done()
			cancel()
		}()

		cmd.SetContext(ctx)
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(ctx context.Context) {
	// we need to always initConfig due to plugin flags needing to get
	// registered with help and usage commands
	initConfig()
	initPlugins()

	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cmdTimeout, "commandTimeout", "1m", "the default timeout for network commands executed")
	var err error
	commandTimeout, err = time.ParseDuration(cmdTimeout)
	if err != nil {
		log.Fatal("could not parse commandTimeout: ", err)
	}

	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	file := path.Join(home, AO_CONFIG_NAME)
	if exists, err := los.IsFileOrFolderExisting(file); err != nil {
		panic(err)
	} else if !exists {
		f, err := os.Create(file)
		if err != nil {
			panic(err)
		}
		f.Close()
	}

	rootCmd.PersistentFlags().String("config", home, "Configuration file to use during execution")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigType("toml")
	viper.SetEnvPrefix(AO_ENV_PREFIX)
	viper.AutomaticEnv()

	configFile := viper.GetString("config")
	info, err := os.Stat(configFile)
	if err == nil && info.IsDir() {
		viper.AddConfigPath(configFile)
		viper.SetConfigName(AO_CONFIG_NAME)
	} else {
		viper.SetConfigFile(configFile)
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Sprintf("no configuration file found: %v", err))
		} else {
			panic(err)
		}
	}
}

// Initialize all the plugins with the loginCmd
func initPlugins() {
	viper.SetDefault(envs.CORE_PLUGINS, []string{"oidc", "eks", "ecr"})
	plugins := viper.GetStringSlice(envs.CORE_PLUGINS)
	for _, p := range plugins {
		err := dlogin.Init(p, loginCmd)
		if err != nil {
			panic(err)
		}
	}
}
