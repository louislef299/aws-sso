package aws

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sso/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	lregion "github.com/louislef299/aws-sso/internal/region"
	los "github.com/louislef299/aws-sso/pkg/os"
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
