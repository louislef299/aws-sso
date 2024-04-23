package aws

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sso/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	lregion "github.com/louislef299/aws-sso/internal/region"
	los "github.com/louislef299/aws-sso/pkg/v1/os"
)

var (
	AwsRegions = []string{
		"us-east-2",
		"us-east-1",
		"us-west-1",
		"us-west-2",
		"af-south-1",
		"ap-east-1",
		"ap-south-1",
		"ap-northeast-3",
		"ap-northeast-2",
		"ap-southeast-1",
		"ap-southeast-2",
		"ap-northeast-1",
		"ca-central-1",
		"eu-central-1",
		"eu-west-1",
		"eu-west-2",
		"eu-south-1",
		"eu-west-3",
		"eu-north-1",
		"me-south-1",
		"sa-east-1",
		"us-gov-east-1",
		"us-gov-west-1",
		"cn-north-1",
		"cn-northwest-1",
	}

	ErrEmptyResponse  = errors.New("an empty response was returned")
	ErrRegionInvalid  = errors.New("the provided region is invalid")
	ErrRegionNotFound = errors.New("could not find a region in the system")
)

const LastUsageLocation = "/.aws/sso/cache/last-usage.json"

type LastUsageInformation struct {
	AccountId   string `json:"account_id"`
	AccountName string `json:"account_name"`
	Role        string `json:"role"`
}

// Gather sts caller identity
func GetCallerIdentity(ctx context.Context, cfg *aws.Config) (*sts.GetCallerIdentityOutput, error) {
	svc := sts.NewFromConfig(*cfg)
	return svc.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
}

func GetURL() (string, error) {
	r, err := lregion.GetRegion(lregion.STS)
	if err != nil {
		return "", err
	}

	switch r {
	case "cn-north-1", "cn-northwest-1":
		return "amazonaws.com.cn", nil
	default:
		return "amazonaws.com", nil
	}
}

func SaveUsageInformation(accountInfo *types.AccountInfo, roleInfo *types.RoleInfo) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	target := homeDir + LastUsageLocation
	usageInformation := LastUsageInformation{
		AccountId:   *accountInfo.AccountId,
		AccountName: *accountInfo.AccountName,
		Role:        *roleInfo.RoleName,
	}
	log.Printf("saving data to %s", target)
	return los.WriteStructToFile(usageInformation, target)
}
