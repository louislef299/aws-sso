package account

import (
	"fmt"
	"log"
	"slices"
	"strings"

	lenv "github.com/louislef299/knot/internal/envs"
	lregion "github.com/louislef299/knot/internal/region"
	laws "github.com/louislef299/knot/pkg/aws"
	los "github.com/louislef299/knot/pkg/os"
	"github.com/spf13/viper"
)

const (
	ACCOUNT_GROUP  = "account"
	ACCOUNT_REGION = "region"
	ACCOUNT_ID     = "id"
	ACCOUNT_URL    = "url"
)

// Represents an AWS
// Account(https://docs.aws.amazon.com/accounts/latest/reference/accounts-welcome.html)
type Account struct {
	ID      string
	Private bool
	Region  string
	Token   string
	URL     string
}

func AddAccount(name string, account *Account) error {
	ViperSetAccount(name, *account)
	return viper.WriteConfig()
}

// Sets account value in Viper. Does not write to config
func ViperSetAccount(name string, acct Account) {
	viper.Set(fmt.Sprintf("account.%s", name), acct)
}

func GetAccountID(profile string) string {
	id := viper.GetString(fmt.Sprintf("%s.%s.%s", ACCOUNT_GROUP, profile, ACCOUNT_ID))
	if id == "" {
		log.Printf("couldn't find an account ID matching profile %s, using empty default...\n", profile)
	}
	return id
}

func GetAccountURL(profile string) string {
	url := viper.GetString(fmt.Sprintf("%s.%s.%s", ACCOUNT_GROUP, profile, ACCOUNT_URL))
	if url == "" {
		log.Printf("couldn't find an account URL matching profile %s, using core default...\n", url)
	} else {
		return url
	}
	return viper.GetString(lenv.CORE_URL)
}

func GetAccountRegion(profile string) string {
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

func GetAccountToken(profile string) string {
	t := viper.GetString(fmt.Sprintf("%s.%s.%s", ACCOUNT_GROUP, profile, "token"))
	if t == "" {
		return "default"
	}
	return t
}

func GetAccountPrivate(profile string) bool {
	return viper.GetBool(fmt.Sprintf("%s.%s.%s", ACCOUNT_GROUP, profile, "private"))
}

func ListAccounts(all bool) {
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

func trimSuffixes(s string) string {
	suffixes := []string{"id", "region"}
	for _, suff := range suffixes {
		s = strings.TrimSuffix(s, fmt.Sprintf(".%s", suff))
	}
	return s
}
