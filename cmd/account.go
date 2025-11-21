package cmd

import (
	"errors"
	"fmt"
	"log"

	lacct "github.com/louislef299/knot/internal/account"
	lenv "github.com/louislef299/knot/internal/envs"
	lregion "github.com/louislef299/knot/internal/region"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	ROLE_ARN_KEY = "role_arn"
)

var (
	accountNumber  string
	accountName    string
	accountRegion  string
	accountToken   string
	accountPrivate bool
	accountURL     string
	allAccounts    bool

	ErrNoAccountFound = errors.New("no account found")
	ErrAccountsToFix  = errors.New("accounts need to be fixed")
)

// accountCmd represents the account command
var accountCmd = &cobra.Command{
	Use:     "account",
	Aliases: []string{"acct"},
	Short:   "Manage AWS account aliases",
	Long: `You can associate AWS account IDs to an alias that is
used by the login command. These values are stored
in your config file.`,
}

// accountAddCmd represents the account command
var accountAddCmd = &cobra.Command{
	Use:     "add",
	Short:   "Associate an alias to an AWS account ID and default region.",
	Example: "  aws-sso account add --name env1 --number 000000000 --region us-west-2",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		if accountRegion == "" {
			accountRegion, err = lregion.GetRegion(lregion.EKS)
			if err != nil {
				log.Fatal("couldn't get default region:", err)
			}
		}
		if accountURL == "" {
			accountURL = lacct.GetAccountURL(accountName)
		}

		err = lacct.AddAccount(accountName, &lacct.Account{
			ID:      accountNumber,
			Region:  accountRegion,
			Private: accountPrivate,
			Token:   accountToken,
			URL:     accountURL,
		})
		if err != nil {
			log.Fatal("couldn't write to configuration file:", err)
		}
		log.Printf("associated account %s to account number %s\n", accountName, accountNumber)
	},
}

// accountListCmd represents the account command
var accountListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List the local AWS account mappings.",
	Long: `Lists the local AWS account alias mapping and AWS name
profiles found.`,
	Run: func(cmd *cobra.Command, args []string) {
		lacct.ListAccounts(allAccounts)
	},
}

// accountSetCmd represents the account command
var accountSetCmd = &cobra.Command{
	Use:     "set",
	Short:   "Set the AWS account profile values.",
	Args:    cobra.ExactArgs(1),
	Example: `  aws-sso account set env1 --number 000000000 --region us-west-2 --token dev --private`,
	Run: func(cmd *cobra.Command, args []string) {
		if accountNumber == "" {
			accountNumber = lacct.GetAccountID(args[0])
		}
		if accountRegion == "" {
			accountRegion = lacct.GetAccountRegion(args[0])
		}
		if accountURL == "" {
			accountURL = lacct.GetAccountURL(args[0])
		}

		lacct.ViperSetAccount(args[0], lacct.Account{
			ID:      accountNumber,
			Region:  accountRegion,
			Private: accountPrivate,
			Token:   accountToken,
			URL:     accountURL,
		})
		err := viper.WriteConfig()
		if err != nil {
			log.Fatal("couldn't write to configuration file:", err)
		}
		fmt.Printf("account values have been set for %s\n", args[0])
	},
}

// accountPluralCmd represents the account command
var accountPluralCmd = &cobra.Command{
	Use:     "accounts",
	Short:   "List the AWS account mappings.",
	Aliases: []string{"accts"},
	Hidden:  true,
	Run: func(cmd *cobra.Command, args []string) {
		lacct.ListAccounts(allAccounts)
	},
}

func init() {
	rootCmd.AddCommand(accountCmd)
	rootCmd.AddCommand(accountPluralCmd)

	accountCmd.AddCommand(accountListCmd)
	accountCmd.AddCommand(accountAddCmd)
	accountCmd.AddCommand(accountSetCmd)

	accountPluralCmd.Flags().BoolVarP(&allAccounts, "all", "a", false, "List all accounts, including AWS config profiles")
	accountListCmd.Flags().BoolVarP(&allAccounts, "all", "a", false, "List all accounts, including AWS config profiles")

	accountAddCmd.Flags().StringVarP(&accountRegion, "region", "r", "", "The default region to associate to the account")
	accountAddCmd.Flags().StringVarP(&accountToken, "token", "t", "default", "The token to use for the account")
	accountAddCmd.Flags().StringVar(&accountURL, "url", "", "The SSO URL to use for the account")
	accountAddCmd.Flags().BoolVarP(&accountPrivate, "private", "p", false, "The account is a private account")
	accountAddCmd.Flags().StringVar(&accountName, "name", "", "The logical name of the account being added")
	if err := accountAddCmd.MarkFlagRequired("name"); err != nil {
		log.Fatal("couldn't mark flag as required:", err)
	}
	accountAddCmd.Flags().StringVar(&accountNumber, "id", "", "The account id of the account associated to the account name")
	if err := accountAddCmd.MarkFlagRequired("id"); err != nil {
		log.Fatal("couldn't mark flag as required:", err)
	}

	accountSetCmd.Flags().StringVarP(&accountRegion, "region", "r", "", "The default region to associate to the account")
	accountSetCmd.Flags().StringVar(&accountNumber, "id", "", "The account id of the account associated to the account name")
	accountSetCmd.Flags().StringVarP(&accountToken, "token", "t", "default", "The token to use for the account")
	accountSetCmd.Flags().StringVar(&accountURL, "url", "", "The SSO URL to use for the account")
	accountSetCmd.Flags().BoolVarP(&accountPrivate, "private", "p", false, "The account is a private account")
}

func syncAccountRegionSession(profile, region string) (string, error) {
	if region == "" {
		region = lacct.GetAccountRegion(profile)
	}
	viper.Set(lenv.SESSION_REGION, region)
	err := viper.WriteConfig()
	return region, err
}
