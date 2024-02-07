/*
Copyright © 2023 Louis Lefebvre <louislefebvre1999@gmail.com>
*/
package cmd

import (
	"errors"
	"fmt"
	"log"
	"slices"

	laws "github.com/louislef299/aws-sso/pkg/v1/aws"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	ACCOUNT_GROUP = "account"
	ROLE_ARN_KEY  = "role_arn"

	accountIDRegex = `\d{12}`
)

var (
	accountNumber string
	accountName   string

	ErrNoAccountFound = errors.New("no account found")
)

// accountCmd represents the account command
var accountCmd = &cobra.Command{
	Use:     "account",
	Aliases: []string{"acct"},
	Short:   "Manage AWS account aliases.",
	Long: `You can associate AWS account IDs to an alias that is
used by the login command. These values are stored
in your config file.`,
}

// accountAddCmd represents the account command
var accountAddCmd = &cobra.Command{
	Use:     "add",
	Short:   "Associate an alias to an AWS account ID.",
	Example: "  aws-sso account add --name env1 --number 000000000",
	Run: func(cmd *cobra.Command, args []string) {
		viperAddAccount(accountName, accountNumber)
		err := viper.WriteConfig()
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

// accountPluralCmd represents the account command
var accountPluralCmd = &cobra.Command{
	Use:    "accounts",
	Short:  "List the AWS account mappings.",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		listAccounts()
	},
}

func init() {
	rootCmd.AddCommand(accountCmd)
	rootCmd.AddCommand(accountPluralCmd)

	accountCmd.AddCommand(accountListCmd)
	accountCmd.AddCommand(accountAddCmd)

	accountAddCmd.Flags().StringVar(&accountName, "name", "", "The logical name of the account being added")
	if err := accountAddCmd.MarkFlagRequired("name"); err != nil {
		log.Fatal("couldn't mark flag as required:", err)
	}
	accountAddCmd.Flags().StringVar(&accountNumber, "number", "", "The account number of the account associated to the account name")
	if err := accountAddCmd.MarkFlagRequired("number"); err != nil {
		log.Fatal("couldn't mark flag as required:", err)
	}
}

// Sets account value in Viper. Does not write to config
func viperAddAccount(name, id string) {
	viper.Set(fmt.Sprintf("account.%s", name), id)
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
	acctList := accts.AllKeys()
	slices.Sort(acctList)
	for _, a := range acctList {
		fmt.Printf("%s: %s\n", a, accts.GetString(a))
	}

	fmt.Printf("\nAWS Configs:\n")

	sections, err := laws.GetAWSProfiles()
	if err != nil {
		log.Fatal(err)
	}
	for _, s := range sections {
		fmt.Println(s)
	}
}
