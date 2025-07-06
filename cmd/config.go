package cmd

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"strings"
	"text/tabwriter"

	"github.com/louislef299/aws-sso/internal/envs"
	"github.com/louislef299/aws-sso/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	acctGroupRegex    = `^account\.*`
	sessionGroupRegex = `^session\.*`

	allConfigValues bool
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"conf"},
	Short:   "Local configuration used for aws-sso",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := rootCmd.PersistentPreRunE(cmd, args); err != nil {
			log.Fatal(err)
		}

		cmd.Println("Using config file", viper.ConfigFileUsed())
	},
}

// configGetCmd represents the get command
var configGetCmd = &cobra.Command{
	Use:     "get",
	Short:   "Get a configuration value.",
	Example: "  aws-sso config get name",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Println("must provide at least one configuration value retrieve")
			os.Exit(1)
		}
		for _, arg := range args {
			cmd.Printf("value of %s: %v\n", arg, viper.Get(arg))
		}
	},
}

// configListCmd represents the list command
var configListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List your local configuration values.",
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("Current config values:")
		keys := viper.AllKeys()
		acctRegex, sessRegex := regexp.MustCompile(acctGroupRegex), regexp.MustCompile(sessionGroupRegex)
		slices.Sort(keys)
		for _, k := range keys {
			if (acctRegex.MatchString(k) || sessRegex.MatchString(k)) && !allConfigValues {
				continue
			}

			value := viper.Get(k)
			if value == "" && !allConfigValues {
				continue
			}
			cmd.Printf("%s=%v\n", k, value)
		}
	},
}

// configSetCmd represents the set command
var configSetCmd = &cobra.Command{
	Use:     "set",
	Short:   "Set a local configuration value.",
	Example: "  aws-sso config set name Louis Lefebvre",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			log.Println("usage: aws-sso config set <key> <value>")
			os.Exit(1)
		}

		viper.Set(args[0], strings.Join(args[1:], " "))
		if err := viper.WriteConfig(); err != nil {
			log.Fatal("couldn't write to config:", err)
		}

		cmd.Printf("set %s to %v\n", args[0], viper.Get(args[0]))
	},
}

// configUnsetCmd represents the unset command
var configUnsetCmd = &cobra.Command{
	Use:   "unset",
	Short: "Unset your config settings",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal("must provide one key to unset. example: aws-sso config unset core.profile")
		}

		viper.Set(args[0], "")
		if err := viper.WriteConfig(); err != nil {
			log.Fatal("could not write to config file:", err)
		}
		cmd.Println("successfully unset", args[0])
	},
}

// configValuesCmd represents the values command
var configValuesCmd = &cobra.Command{
	Use:     "values",
	Aliases: []string{"vals"},
	Short:   "Get the possible configuration values.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("The following values are available for configuration:")
		w := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight)
		for _, c := range config.GetCurrentConfigValues() {
			fmt.Fprintln(w, c.String())
		}
		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.AddCommand(configListCmd)
	configListCmd.Flags().BoolVarP(&allConfigValues, "all", "a", false, "List all configuration values, including tool internal values")

	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configUnsetCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configValuesCmd)

	// Space saved for config values not bound to flags
	config.AddConfigValue(envs.CORE_DEFAULT_CLUSTER, "The default cluster to target when logging in, supports go regex expressions(golang.org/s/re2syntax)")
	config.AddConfigValue(envs.CORE_DEFAULT_ROLE, "The default iam role to use when logging in, supports go regex expressions(golang.org/s/re2syntax)")
	config.AddConfigValue(envs.CORE_DEFAULT_REGION, "The default region used when a region is not found in your environment or set with flags")
	config.AddConfigValue(envs.CORE_SSO_REGION, "The region to use for the AWS SSO authentication")
	config.AddConfigValue(envs.CORE_URL, "The default sso start url used when logging in")
	config.AddConfigValue(envs.CORE_BROWSER, "Default browser is required for advanced features like opening in a private browser")
}
