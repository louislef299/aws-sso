package cmd

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/louislef299/aws-sso/internal/envs"
	lconfig "github.com/louislef299/aws-sso/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	ACTIVE_TOKEN_ID   = "active"
	INACTIVE_TOKEN_ID = "inactive"
)

// tokenCmd represents the token command
var tokenCmd = &cobra.Command{
	Use:     "token",
	Aliases: []string{"tok", "to", "ken"},
	Short:   "Manage multiple SSO access tokens",
	Long: `Manages multiple cached SSO tokens for reuse. Beneficial
when dealing with multiple AWS Organizations.`,
	PersistentPreRun: checkTokenV,
}

// tokensCmd represents the list command
var tokensCmd = &cobra.Command{
	Use:              "tokens",
	Hidden:           true,
	Short:            "List your tokens.",
	PersistentPreRun: checkTokenV,
	Run: func(cmd *cobra.Command, args []string) {
		listTokens()
	},
}

// tokenListCmd represents the list command
var tokenListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List your tokens",
	Run: func(cmd *cobra.Command, args []string) {
		listTokens()
	},
}

// tokenAddCmd represents the add command
var tokenAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a token.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		addToken(args[0])
		log.Printf("successfully added token %s!\n", args[0])
	},
}

// tokenCurrentCmd represents the current command
var tokenCurrentCmd = &cobra.Command{
	Use:     "current",
	Aliases: []string{"cur", "curr"},
	Short:   "List the current token.",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("the current token is", getCurrentToken())
	},
}

// tokenRemoveCmd represents the remove command
var tokenRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm"},
	Short:   "Remove a token.",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		t := getToken(args[0])
		if t == "" {
			log.Printf("token '%s' was not found\n", args[0])
			return
		} else if strings.Compare(t, ACTIVE_TOKEN_ID) == 0 {
			log.Printf("token '%s' is currently active, removing an active token is not allowed\n", args[0])
			return
		}
		removeToken(args[0])
		log.Printf("successfully removed token %s!\n", args[0])
	},
}

// tokenUseCmd represents the use command
var tokenUseCmd = &cobra.Command{
	Use:   "use",
	Short: "Use a token.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !doesTokenExist(args[0]) {
			log.Printf("token '%s' was not found\n", args[0])
			return
		}
		setToken(args[0])
		log.Printf("using token %s\n", args[0])
	},
}

func init() {
	rootCmd.AddCommand(tokenCmd)
	rootCmd.AddCommand(tokensCmd)

	tokenCmd.AddCommand(tokenListCmd)
	tokenCmd.AddCommand(tokenAddCmd)
	tokenCmd.AddCommand(tokenRemoveCmd)
	tokenCmd.AddCommand(tokenCurrentCmd)
	tokenCmd.AddCommand(tokenUseCmd)
}

func setToken(name string) {
	if c := getCurrentToken(); c != "" {
		viper.Set(fmt.Sprintf("%s.%s", envs.TOKEN_HEADER, c), INACTIVE_TOKEN_ID)
	}
	if name == "-" {
		name = envs.DEFAULT_TOKEN_NAME
	}
	viper.Set(fmt.Sprintf("%s.%s", envs.TOKEN_HEADER, name), ACTIVE_TOKEN_ID)
	viper.Set(envs.SESSION_TOKEN, name)
	if err := viper.WriteConfig(); err != nil {
		log.Fatal("couldn't write config:", err)
	}
}

// Viper wrapper
func checkTokenV(cmd *cobra.Command, args []string) { checkToken() }

// Quick check to make sure the session token is set
func checkToken() {
	if !viper.IsSet(envs.SESSION_TOKEN) || getCurrentToken() == "" {
		setToken(envs.DEFAULT_TOKEN_NAME)
	}
}

// Returns the current session token
func addToken(name string) {
	lconfig.DeepSet(fmt.Sprintf("%s.%s", envs.TOKEN_HEADER, name), ACTIVE_TOKEN_ID)
	setToken(name)
}

func getCurrentToken() string         { return viper.GetString(envs.SESSION_TOKEN) }
func doesTokenExist(name string) bool { return (getToken(name) != "") }
func getToken(name string) string {
	return viper.GetString(fmt.Sprintf("%s.%s", envs.TOKEN_HEADER, name))
}

func listTokens() {
	tokens := viper.Sub(envs.TOKEN_HEADER)
	fmt.Println("Local Tokens:")

	if tokens == nil {
		// The default token will always exist
		setToken(envs.DEFAULT_TOKEN_NAME)
		fmt.Printf("* %s\n", envs.DEFAULT_TOKEN_NAME)
		return
	}

	tokenList := tokens.AllKeys()
	slices.Sort(tokenList)
	for _, t := range tokenList {
		if strings.Compare(tokens.GetString(t), ACTIVE_TOKEN_ID) == 0 {
			fmt.Printf("* %s\n", t)
		} else if strings.Compare(tokens.GetString(t), INACTIVE_TOKEN_ID) == 0 {
			fmt.Printf("  %s\n", t)
		}
	}
}

func removeToken(name string) {
	lconfig.DeepSet(fmt.Sprintf("%s.%s", envs.TOKEN_HEADER, name), "")
}
