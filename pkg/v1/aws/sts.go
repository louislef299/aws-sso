package aws

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sso/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	. "github.com/louislef299/aws-sso/internal/envs"
	los "github.com/louislef299/aws-sso/pkg/v1/os"
	"github.com/spf13/viper"
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

// Returns the region in precedence of environment
// region, config region and finally default region.
func GetRegion() (string, error) {
	var r string
	checkOutput := func(r string) error {
		if r != "" && validRegion(r) {
			return nil
		}
		if r == "" {
			return ErrEmptyResponse
		}
		if !validRegion(r) {
			return ErrRegionInvalid
		}
		return ErrRegionNotFound
	}

	r = os.Getenv("AWS_REGION")
	if err := checkOutput(r); err == nil {
		return r, nil
	} else if err == ErrRegionInvalid {
		return "", err
	}

	r = os.Getenv("AWS_DEFAULT_REGION")
	if err := checkOutput(r); err == nil {
		return r, nil
	} else if err == ErrRegionInvalid {
		return "", err
	}

	r = viper.GetString(SESSION_REGION)
	if err := checkOutput(r); err == nil {
		return r, nil
	} else if err == ErrRegionInvalid {
		return "", err
	}

	r = viper.GetString(CORE_REGION)
	if err := checkOutput(r); err == nil {
		return r, nil
	} else if err == ErrRegionInvalid {
		return "", err
	}

	return "", ErrRegionNotFound
}

func GetURL() (string, error) {
	r, err := GetRegion()
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

	target := homeDir + "/.aws/sso/cache/last-usage.json"
	usageInformation := LastUsageInformation{
		AccountId:   *accountInfo.AccountId,
		AccountName: *accountInfo.AccountName,
		Role:        *roleInfo.RoleName,
	}
	log.Printf("saving data to %s", target)
	return los.WriteStructToFile(usageInformation, target)
}

func validRegion(region string) bool {
	for _, r := range AwsRegions {
		if strings.Compare(r, strings.TrimSpace(region)) == 0 {
			return true
		}
	}
	return false
}
