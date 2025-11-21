package cmd

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"strings"
	"text/tabwriter"

	"github.com/louislef299/knot/internal/envs"
	"github.com/louislef299/knot/pkg/config"
	"github.com/louislef299/knot/pkg/dlogin"
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
	Short:   "List your core local configuration values.",
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

		tgt := args[0]
		params := strings.Join(args[1:], " ")

		// validateInput may set the value directly for special cases
		validateInput(tgt, params)

		// Only set if validateInput didn't handle it
		if tgt != envs.CORE_PLUGINS {
			viper.Set(tgt, params)
		}

		if err := viper.WriteConfig(); err != nil {
			log.Fatal("couldn't write to config:", err)
		}

		cmd.Printf("set %s to %v\n", tgt, viper.Get(tgt))
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
	configListCmd.Flags().BoolVar(&allConfigValues, "all", false, "List all configuration values, including tool internal values")

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

func validateInput(target string, params string) {
	switch target {
	case envs.CORE_PLUGINS:
		// Parse and validate plugin list
		plugins := parsePlugins(params)
		if len(plugins) == 0 {
			log.Fatal("core.plugins cannot be empty. Available plugins: ", registeredPluginDrivers())
		}

		// Validate all plugins exist
		available := registeredPluginDrivers()
		for _, p := range plugins {
			if !slices.Contains(available, p) {
				log.Fatalf("unknown plugin %q. Available plugins: %v", p, available)
			}
		}

		// Set as proper string slice instead of string
		viper.Set(target, plugins)
	}
}

// parsePlugins handles various plugin config formats and returns a clean string slice
// Accepts:
// - TOML arrays: []string{"oidc", "eks", "ecr"} or []interface{}{"oidc", "eks", "ecr"}
// - Comma-separated: "oidc,eks,ecr" or "[oidc,eks]"
// - Space-separated: "oidc eks ecr"
// - Single plugin: "oidc"
func parsePlugins(value interface{}) []string {
	// Handle nil or empty
	if value == nil {
		return []string{}
	}

	// If it's already a string slice, check if it needs parsing
	if slice, ok := value.([]string); ok {
		if len(slice) == 0 {
			return []string{}
		}
		// If single element, might be a concatenated string
		if len(slice) == 1 {
			parsed := parsePluginString(slice[0])
			// Only return parsed if we got multiple items, or if the single item changed
			if len(parsed) > 1 || (len(parsed) == 1 && parsed[0] != slice[0]) {
				return parsed
			}
		}
		return slice
	}

	// If it's an interface slice (common when viper reads TOML arrays)
	if slice, ok := value.([]interface{}); ok {
		if len(slice) == 0 {
			return []string{}
		}
		// Convert []interface{} to []string
		var result []string
		for _, v := range slice {
			if str, ok := v.(string); ok {
				result = append(result, str)
			}
		}
		return result
	}

	// If it's a string, parse it
	if str, ok := value.(string); ok {
		return parsePluginString(str)
	}

	return []string{}
}

// parsePluginString parses a plugin string in various formats
func parsePluginString(input string) []string {
	input = strings.TrimSpace(input)
	if input == "" {
		return []string{}
	}

	// Remove surrounding brackets
	if len(input) > 0 && input[0] == '[' {
		input = input[1:]
	}
	if len(input) > 0 && input[len(input)-1] == ']' {
		input = input[:len(input)-1]
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return []string{}
	}

	// Detect separator type: comma takes precedence over space
	var parts []string
	if strings.Contains(input, ",") {
		// Comma-separated
		for _, part := range strings.Split(input, ",") {
			cleaned := cleanPluginName(part)
			if cleaned != "" {
				parts = append(parts, cleaned)
			}
		}
	} else if strings.Contains(input, " ") {
		// Space-separated
		for _, part := range strings.Split(input, " ") {
			cleaned := cleanPluginName(part)
			if cleaned != "" {
				parts = append(parts, cleaned)
			}
		}
	} else {
		// Single plugin
		cleaned := cleanPluginName(input)
		if cleaned != "" {
			parts = []string{cleaned}
		}
	}

	return parts
}

// cleanPluginName removes quotes and whitespace from plugin names
func cleanPluginName(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, "\"'")
	return strings.TrimSpace(s)
}

// registeredPluginDrivers returns list of all registered plugin drivers
// that are available to be activated. This is used for validation when
// users set core.plugins configuration.
// Note: This must be called after plugins are imported in root.go
func registeredPluginDrivers() []string {
	// Since plugins are registered via init() in imported packages,
	// we can safely call this after imports
	return dlogin.Drivers()
}
