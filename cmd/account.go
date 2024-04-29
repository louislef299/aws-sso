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

	lregion "github.com/louislef299/aws-sso/internal/region"
	laws "github.com/louislef299/aws-sso/pkg/v1/aws"
	los "github.com/louislef299/aws-sso/pkg/v1/os"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	ACCOUNT_GROUP = "account"
	ROLE_ARN_KEY  = "role_arn"
)

var (
	accountNumber string
	accountName   string
	accountRegion string

	ErrNoAccountFound = errors.New("no account found")
	ErrAccountsToFix  = errors.New("accounts need to be fixed")
)

type Account struct {
	ID     string
	Region string
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

		err = addAccount(accountName, accountNumber, region)
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
		listAccounts()
	},
}

// accountFixupCmd represents the account command
var accountFixupCmd = &cobra.Command{
	Use:   "fixup",
	Short: "Reformats any account profiles using deprecated pattern.",
	Long: `Fixup reformats any account profiles using deprecated pattern.
The reformatting will move from simply associating a profile to 
an account ID to a more structured account with an ID and 
region.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := fixAccounts(false)
		if err != nil {
			log.Fatal("couldn't fix accounts:", err)
		}
		log.Println("accounts have been reformatted!")
	},
}

// accountPluralCmd represents the account command
var accountPluralCmd = &cobra.Command{
	Use:     "accounts",
	Short:   "List the AWS account mappings.",
	Aliases: []string{"accts"},
	Hidden:  true,
	Run: func(cmd *cobra.Command, args []string) {
		listAccounts()
	},
}

func init() {
	rootCmd.AddCommand(accountCmd)
	rootCmd.AddCommand(accountPluralCmd)

	accountCmd.AddCommand(accountListCmd)
	accountCmd.AddCommand(accountAddCmd)
	accountCmd.AddCommand(accountFixupCmd)

	accountAddCmd.Flags().StringVarP(&accountRegion, "region", "r", "", "The default region to associate to the account")
	accountAddCmd.Flags().StringVar(&accountName, "name", "", "The logical name of the account being added")
	if err := accountAddCmd.MarkFlagRequired("name"); err != nil {
		log.Fatal("couldn't mark flag as required:", err)
	}
	accountAddCmd.Flags().StringVar(&accountNumber, "number", "", "The account number of the account associated to the account name")
	if err := accountAddCmd.MarkFlagRequired("number"); err != nil {
		log.Fatal("couldn't mark flag as required:", err)
	}
}

func addAccount(name, id, region string) error {
	viperAddAccount(name, Account{ID: id, Region: region})
	return viper.WriteConfig()
}

// Sets account value in Viper. Does not write to config
func viperAddAccount(name string, acct Account) {
	viper.Set(fmt.Sprintf("account.%s", name), acct)
}

func getAccountID(profile string) string {
	id := viper.GetString(fmt.Sprintf("%s.%s", ACCOUNT_GROUP, profile))
	if id == "" {
		log.Printf("couldn't find an account ID matching profile %s, using empty default...\n", profile)
	}
	return id
}

func listAccounts() {
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

	var account Account
	for _, a := range acctList {
		err := accts.UnmarshalKey(a, &account)
		if err == nil {
			fmt.Printf("%s:\n  ID: %s\n  Region: %s\n", a, account.ID, account.Region)
		}
	}

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

func trimSuffixes(s string) string {
	suffixes := []string{"id", "region"}
	for _, suff := range suffixes {
		s = strings.TrimSuffix(s, fmt.Sprintf(".%s", suff))
	}
	return s
}

// fixAccounts is a temporary function to fix the account structures of the
// existing name: id string configuration.
func fixAccounts(check bool) error {
	accts := viper.Sub(ACCOUNT_GROUP)
	if accts == nil {
		return nil
	}

	acctKeys := accts.AllKeys()
	var acctList []string
	for _, k := range acctKeys {
		acctList = append(acctList, trimSuffixes(k))
	}
	acctList = slices.Compact(acctList)

	slices.Sort(acctList)
	count := 0
	defaultRegion, err := lregion.GetRegion(lregion.EKS)
	if err != nil {
		return err
	}

	var account Account
	for _, a := range acctList {
		err := accts.UnmarshalKey(a, &account)
		if err != nil {
			if !check {
				fmt.Printf("reformatting profile %s to new structure...\n", a)
				viperAddAccount(a, Account{ID: accts.GetString(a), Region: defaultRegion})
			}
			count++
		}
	}
	if count == 0 {
		fmt.Println("no accounts to reformat!")
		return nil
	} else if check {
		return ErrAccountsToFix
	}

	return viper.WriteConfig()
}
