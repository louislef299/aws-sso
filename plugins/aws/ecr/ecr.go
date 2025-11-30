package ecr

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/docker/cli/cli/config"
	"github.com/docker/docker/registry"
	lconfig "github.com/louislef299/knot/internal/config"
	lregion "github.com/louislef299/knot/internal/region"
	laws "github.com/louislef299/knot/pkg/aws"
	"github.com/louislef299/knot/pkg/dlogin"
	ldocker "github.com/louislef299/knot/pkg/docker"
	"github.com/louislef299/knot/pkg/os"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const ECR_DISABLE_ECR_LOGIN = "ecr.disableECRLogin"

type ECRLogin struct {
	Username string
	Config   *aws.Config
}

func init() {
	dlogin.Register("ecr", &ECRLogin{})
}

func (e *ECRLogin) Init(cmd *cobra.Command) error {
	err := dlogin.Activate("ecr")
	if err != nil {
		return err
	}

	cmd.Flags().Bool("disableECRLogin", false, "Disables automatic detection and login for ECR")
	lconfig.AddConfigValue(ECR_DISABLE_ECR_LOGIN, "Disables automatic detection and login for ECR")

	return viper.BindPFlag(ECR_DISABLE_ECR_LOGIN, cmd.Flags().Lookup("disableECRLogin"))
}

func (a *ECRLogin) Login(ctx context.Context, config any) error {
	cfg, ok := config.(*ECRLogin)
	if !ok {
		return fmt.Errorf("expected ECRLogin, got %T", config)
	}

	if viper.GetBool(ECR_DISABLE_ECR_LOGIN) {
		log.Println("ECR Plugin is disabled, skipping...")
		return nil
	}

	ecrToken, ecrEndpoint, err := GetECRToken(ctx, cfg.Config)
	if err != nil {
		return fmt.Errorf("couldn't gather ecr token: %v", err)
	}

	return ldocker.Login("AWS", ecrToken, ecrEndpoint)
}

func (a *ECRLogin) Logout(ctx context.Context, config any) error {
	cfg, ok := config.(*ECRLogin)
	if !ok {
		return fmt.Errorf("expected ECRLogin, got %T", config)
	}

	if viper.GetBool(ECR_DISABLE_ECR_LOGIN) {
		log.Println("ECR Plugin is disabled, skipping...")
		return nil
	}

	// clean docker configs
	registry, err := GetECRRegistryName(ctx, cfg.Config)
	if err != nil {
		return fmt.Errorf("couldn't logout of docker: %v", err)
	} else {
		err = DockerLogout(registry)
		if err != nil {
			return fmt.Errorf("could not logout of ECR registry: %v", err)
		}
	}
	return nil
}

func DockerLogout(registryname string) error {
	registryname = registry.ConvertToHostname(registryname)
	dcfg, err := config.Load(config.Dir())
	if err != nil {
		return fmt.Errorf("loading config file failed: %v", err)
	}

	// check if we're logged in based on the records in the config file
	// which means it couldn't have user/pass cause they may be in the creds store
	if _, loggedIn := dcfg.AuthConfigs[registryname]; loggedIn {
		if err := dcfg.GetCredentialsStore(registryname).Erase(registryname); err != nil {
			return fmt.Errorf("could not erase credentials: %v", err)
		}
		log.Println("erased", registryname)
	} else {
		log.Println("wasn't logged into", registryname)
	}
	return nil
}

// Returns the name of the ECR registry for the AWS environment
func GetECRRegistryName(ctx context.Context, cfg *aws.Config) (string, error) {
	callerid, err := laws.GetCallerIdentity(ctx, cfg)
	if err != nil {
		return "", fmt.Errorf("could not get caller identity: %v", err)
	}

	r, err := lregion.GetRegion(lregion.ECR)
	if err != nil {
		return "", err
	}
	url, err := laws.GetURL()
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
