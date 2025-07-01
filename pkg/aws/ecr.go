package aws

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	lregion "github.com/louislef299/aws-sso/internal/region"
	"github.com/louislef299/aws-sso/pkg/os"
)

// Returns the name of the ECR registry for the AWS environment
func GetECRRegistryName(ctx context.Context, cfg *aws.Config) (string, error) {
	callerid, err := GetCallerIdentity(ctx, cfg)
	if err != nil {
		return "", fmt.Errorf("could not get caller identity: %v", err)
	}

	r, err := lregion.GetRegion(lregion.ECR)
	if err != nil {
		return "", err
	}
	url, err := GetURL()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s.dkr.ecr.%s.%s", *callerid.Account, r, url), nil
}

// Gather an ECR authentication token and return token, proxy endpoint
func GetECRToken(ctx context.Context, cfg *aws.Config) (string, string, error) {
	svc := ecr.NewFromConfig(*cfg)
	resp, err := svc.GetAuthorizationToken(ctx, &ecr.GetAuthorizationTokenInput{})
	if err != nil {
		return "", "", fmt.Errorf("could not get ECR authorization token: %v", err)
	}

	if len(resp.AuthorizationData) == 0 {
		return "", "", errors.New("no authorization data returned")
	}
	r := resp.AuthorizationData[len(resp.AuthorizationData)-1]
	log.Println("Docker credentials expire at", r.ExpiresAt.Format(time.RFC822))
	d, err := os.Decode(*r.AuthorizationToken)
	if err != nil {
		return "", "", err
	}
	return strings.TrimPrefix(d, "AWS:"), *r.ProxyEndpoint, nil
}
