/*
Copyright © 2023 Louis Lefebvre <louislefebvre1999@gmail.com>
*/
package cmd

import (
	"errors"
	"fmt"
	"log"
	"slices"
	"strings"

	lenv "github.com/louislef299/aws-sso/internal/envs"
	lregion "github.com/louislef299/aws-sso/internal/region"
	laws "github.com/louislef299/aws-sso/pkg/v1/aws"
	los "github.com/louislef299/aws-sso/pkg/v1/os"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	ACCOUNT_GROUP  = "account"
	ACCOUNT_REGION = "region"
	ACCOUNT_ID     = "id"
	ROLE_ARN_KEY   = "role_arn"
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

type Account struct {
	ID      string
	Private bool
	Region  string
	Token   string
	URL     string
}

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
		if region == "" {
			region, err = lregion.GetRegion(lregion.EKS)
			if err != nil {
				log.Fatal("couldn't get default region:", err)
			}
		}

		err = addAccount(accountName, &Account{
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
		listAccounts(allAccounts)
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
			accountNumber = getAccountID(args[0])
		}
		if accountRegion == "" {
			accountRegion = getAccountRegion(args[0])
		}

		viperSetAccount(args[0], Account{
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
		listAccounts(allAccounts)
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
	accountAddCmd.Flags().StringVar(&accountURL, "url", viper.GetString(lenv.CORE_URL), "The SSO URL to use for the account")
	accountAddCmd.Flags().BoolVarP(&accountPrivate, "private", "p", false, "The account is a private account")
	accountAddCmd.Flags().StringVar(&accountName, "name", "", "The logical name of the account being added")
	if err := accountAddCmd.MarkFlagRequired("name"); err != nil {
		log.Fatal("couldn't mark flag as required:", err)
	}
	accountAddCmd.Flags().StringVar(&accountNumber, "number", "", "The account number of the account associated to the account name")
	if err := accountAddCmd.MarkFlagRequired("number"); err != nil {
		log.Fatal("couldn't mark flag as required:", err)
	}

	accountSetCmd.Flags().StringVarP(&accountRegion, "region", "r", "", "The default region to associate to the account")
	accountSetCmd.Flags().StringVar(&accountNumber, "number", "", "The account number of the account associated to the account name")
	accountSetCmd.Flags().StringVarP(&accountToken, "token", "t", "default", "The token to use for the account")
	accountSetCmd.Flags().StringVar(&accountURL, "url", viper.GetString(lenv.CORE_URL), "The SSO URL to use for the account")
	accountSetCmd.Flags().BoolVarP(&accountPrivate, "private", "p", false, "The account is a private account")
}

func addAccount(name string, account *Account) error {
	viperSetAccount(name, *account)
	return viper.WriteConfig()
}

// Sets account value in Viper. Does not write to config
func viperSetAccount(name string, acct Account) {
	viper.Set(fmt.Sprintf("account.%s", name), acct)
}

func getAccountID(profile string) string {
	id := viper.GetString(fmt.Sprintf("%s.%s.%s", ACCOUNT_GROUP, profile, ACCOUNT_ID))
	if id == "" {
		log.Printf("couldn't find an account ID matching profile %s, using empty default...\n", profile)
	}
	return id
}

func getAccountRegion(profile string) string {
	var err error
	r := viper.GetString(fmt.Sprintf("%s.%s.%s", ACCOUNT_GROUP, profile, ACCOUNT_REGION))
	if r == "" {
		log.Printf("couldn't find an account region matching profile %s, using local default...\n", profile)
		r, err = lregion.GetRegion(lregion.EKS)
		if err != nil {
			log.Fatal("couldn't get default region:", err)
		}
	}
	return r
}

func getAccountToken(profile string) string {
	t := viper.GetString(fmt.Sprintf("%s.%s.%s", ACCOUNT_GROUP, profile, "token"))
	if t == "" {
		return "default"
	}
	return t
}

func getAccountPrivate(profile string) bool {
	return viper.GetBool(fmt.Sprintf("%s.%s.%s", ACCOUNT_GROUP, profile, "private"))
}

func listAccounts(all bool) {
	accts := viper.Sub(ACCOUNT_GROUP)
	if accts == nil {
		log.Println("no accounts configured! use the 'account add' command to create new mappings")
		return
	}

	fmt.Println("Account mapping:")
	acctKeys := accts.AllKeys()
	var acctList []string
	for _, k := range acctKeys {
		acctList = append(acctList, trimSuffixes(k))
	}
	slices.Sort(acctList)
	acctList = slices.Compact(acctList)

	for _, a := range acctList {
		var account Account
		err := accts.UnmarshalKey(a, &account)
		var t string
		if account.Token == "" {
			t = "default"
		} else {
			t = account.Token
		}

		var url string
		if account.URL == "" {
			url = "(default) " + viper.GetString(lenv.CORE_URL)
		} else {
			url = account.URL
		}

		if err == nil {
			fmt.Printf("%s:\n  ID: %s\n  Region: %s\n  Private: %t\n  Token: %s\n  SSO URL: %s\n",
				a, account.ID, account.Region, account.Private, t, url)
		}
	}

	if all {
		fmt.Printf("\nAWS Configs:\n")
		sections, err := laws.GetAWSProfiles()
		if err != nil {
			log.Fatal(err)
		}
		for _, s := range sections {
			if !los.IsManagedProfile(s) {
				fmt.Println(s)
			}
		}
	}
}

func syncAccountRegionSession(profile, region string) (string, error) {
	if region == "" {
		region = getAccountRegion(profile)
	}
	viper.Set(lenv.SESSION_REGION, region)
	err := viper.WriteConfig()
	return region, err
}

func trimSuffixes(s string) string {
	suffixes := []string{"id", "region"}
	for _, suff := range suffixes {
		s = strings.TrimSuffix(s, fmt.Sprintf(".%s", suff))
	}
	return s
}
