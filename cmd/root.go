package cmd

import (
	"context"
	"fmt"
	"html/template"
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
	debug              bool
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
	Long: `An AWS SSO helper CLI to streamline authentication.

more information at: https://aws-sso.netlify.app/`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if debug {
			log.SetFlags(log.Ltime | log.Ldate | log.Llongfile)
		}

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

	rootCmd.PersistentFlags().StringVar(&cmdTimeout, "commandTimeout", "1m", "timeout for network commands executed")
	var err error
	commandTimeout, err = time.ParseDuration(cmdTimeout)
	if err != nil {
		log.Fatal("could not parse commandTimeout: ", err)
	}
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "add debug message headers to logger")

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
		defer f.Close()
		tmpl, err := template.New("config").Parse(getConfigTemplate())
		if err != nil {
			panic(err)
		}
		err = tmpl.Execute(f, "")
		if err != nil {
			panic(err)
		}
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

func getConfigTemplate() string {
	return `# Account represents the AWS account alias. This will then be added to the
# aws-sso account list command and allows for aws-sso login <account> to work
# properly.
[account]
[account.dev]
id = '000000000000'
private = false
region = 'us-east-2'
token = 'default'

# Core represents all configurations that can be used across Accounts and
# Plugins. These are useful to aws-sso functioning on your local system.
[core]
browser = 'chrome'
defaultregion = 'us-east-1'
plugins = ['oidc', 'eks', 'ecr']
ssoregion = 'us-east-1'
url = 'https://docs.aws.amazon.com/signin/latest/userguide/sign-in-urls-defined.html'

# The Session and Token sections are managed by the aws-sso CLI tool. You
# typically shouldn't have to mess with these unless there are some low-level
# errors happening on your machine. To get rid of your current session
# altogether, feel free to run aws-sso logout --clean.
[session]

[token]
`
}
