package aws

import (
	"context"
	"errors"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	"github.com/aws/aws-sdk-go-v2/service/sso/types"
	"github.com/louislef299/knot/pkg/prompt"
)

var ErrAccountNotFound = errors.New("account provided couldn't be found")

func RetrieveAccountInformation(ctx context.Context, cfg *aws.Config, accessToken, acctID *string) (*types.AccountInfo, error) {
	client := sso.NewFromConfig(*cfg)
	accounts, err := client.ListAccounts(ctx, &sso.ListAccountsInput{
		AccessToken: accessToken,
		MaxResults:  aws.Int32(100),
	})
	if err != nil {
		return nil, err
	}

	if *acctID == "" {
		// account ID not provided, fuzzy search user input
		sortedAccounts := sortAccounts(accounts.AccountList)

		var accountsToSelect []string
		for i, info := range sortedAccounts {
			accountsToSelect = append(accountsToSelect, prompt.LINEPREFIX+strconv.Itoa(i)+" "+*info.AccountName+" "+*info.AccountId) //+" "+accounts.GetAccountMap()[*info.AccountId])
		}

		label := "Select your account - To choose one account directly just enter #{Int}"
		indexChoice, _ := prompt.Select(label, accountsToSelect, prompt.FuzzySearchWithPrefixAnchor(accountsToSelect))
		accountInfo := sortedAccounts[indexChoice]
		log.Printf("Selected account: %s - %s", *accountInfo.AccountName, *accountInfo.AccountId)
		return &accountInfo, nil
	}
	for _, acct := range accounts.AccountList {
		if strings.Compare(*acct.AccountId, *acctID) == 0 {
			return &acct, nil
		}
	}
	return nil, ErrAccountNotFound
}

func sortAccounts(accountList []types.AccountInfo) []types.AccountInfo {
	var sortedAccounts []types.AccountInfo
	sortedAccounts = append(sortedAccounts, accountList...)

	sort.Slice(sortedAccounts, func(i, j int) bool {
		return *sortedAccounts[i].AccountName < *sortedAccounts[j].AccountName
	})
	return sortedAccounts
}
