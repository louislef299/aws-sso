package oidc

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConf "github.com/aws/aws-sdk-go-v2/config"
	lacct "github.com/louislef299/aws-sso/internal/account"
	"github.com/louislef299/aws-sso/internal/browser"
	"github.com/louislef299/aws-sso/internal/envs"
	lregion "github.com/louislef299/aws-sso/internal/region"
	laws "github.com/louislef299/aws-sso/pkg/aws"
	lconfig "github.com/louislef299/aws-sso/pkg/config"
	"github.com/louislef299/aws-sso/pkg/dlogin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type OIDCLogin struct {
	Config       *aws.Config
	Profile      string
	Region       string
	NewProfile   bool
	Private      bool
	Refresh      bool
	SkipDefaults bool
}

func init() {
	dlogin.Register("oidc", &OIDCLogin{})
}

func (e *OIDCLogin) Init(cmd *cobra.Command) error {
	return nil
}

func (a *OIDCLogin) Login(ctx context.Context, config any) error {
	cfg, ok := config.(*OIDCLogin)
	if !ok {
		return fmt.Errorf("expected OIDCLogin, got %T", config)
	}

	// check if referencing a local profile
	lc, err := laws.IsLocalConfig(cfg.Profile)
	if err != nil {
		log.Println("couldn't find predefined AWS configurations:", err)
	}

	if lc {
		log.Println("using existing configuration profile", cfg.Profile)
		awsCfg, err := awsConf.LoadDefaultConfig(ctx,
			awsConf.WithSharedConfigProfile(cfg.Profile),
		)
		if err != nil {
			return err
		}
		lconfig.DeepSet(envs.SESSION_PROFILE, cfg.Profile)
		cfg.Config = &awsCfg
		return nil
	}

	ssoRegion, err := lregion.GetRegion(lregion.STS)
	if err != nil {
		log.Fatal("could not gather sso region:", err)
	}

	// gross anti-pattern, but too lazy to reprogram the existing logic for now
	if cfg.Region == "cn-north-1" || cfg.Region == "cn-northwest-1" {
		ssoRegion = "cn-north-1"
	}
	log.Println("using sso region", ssoRegion, "to login")

	acctID := lacct.GetAccountID(cfg.Profile)
	awsCfg, err := awsConf.LoadDefaultConfig(ctx, awsConf.WithRegion(ssoRegion))
	if err != nil {
		return err
	}

	p, err := loginAWS(ctx, awsCfg, acctID, cfg)
	if err != nil {
		log.Fatal("couldn't log into AWS: ", err)
	}
	lconfig.DeepSet(envs.SESSION_PROFILE, p)

	log.Println("loading up new config", p, "with region", cfg.Region)
	// Start up new config with newly configured profile
	awsCfg, err = awsConf.LoadDefaultConfig(ctx, awsConf.WithRegion(cfg.Region), awsConf.WithSharedConfigProfile(p))
	if err != nil {
		log.Fatal("couldn't load new config:", err)
	}
	cfg.Config = &awsCfg

	return nil
}

func (a *OIDCLogin) Logout(ctx context.Context, config any) error {
	return nil
}

func loginAWS(ctx context.Context, awsCfg aws.Config, acctID string, cfg *OIDCLogin) (string, error) {
	u := getURL(cfg.Profile)
	lconfig.DeepSet(envs.SESSION_URL, u)

	clientInfo, err := laws.GatherClientInformation(ctx, &awsCfg, u, getBrowser(cfg), cfg.Refresh)
	if err != nil {
		return "", err
	}

	account, err := laws.RetrieveAccountInformation(ctx, &awsCfg, &clientInfo.AccessToken, &acctID)
	if err != nil {
		return "", err
	}
	acctID = *account.AccountId

	role, err := laws.RetrieveRoleInfo(ctx, &awsCfg, *account.AccountId, clientInfo.AccessToken, cfg.SkipDefaults)
	if err != nil {
		return "", err
	}
	log.Println("using aws role", *role.RoleName)
	lconfig.DeepSet(envs.SESSION_ROLE, *role.RoleName)

	err = laws.SaveUsageInformation(account, &role)
	if err != nil {
		return "", err
	}

	// set the new profile in account config
	if cfg.NewProfile {
		err = lacct.AddAccount(cfg.Profile, &lacct.Account{
			ID:      acctID,
			Region:  awsCfg.Region,
			Private: cfg.Private,
			Token:   viper.GetString(envs.SESSION_TOKEN),
		})
		if err != nil {
			log.Println("WARNING: couldn't write to configuration file:", err)
		}
	}
	return laws.GetAndSaveRoleCredentials(ctx, &awsCfg, account.AccountId, role.RoleName, &clientInfo.AccessToken, cfg.Profile, awsCfg.Region)
}

func getBrowser(cfg *OIDCLogin) browser.Browser {
	if cfg.Private {
		log.Println("browser set to open incognito(no cookies)")
	} else {
		log.Println("browser set to default(use cookies)")
	}
	return browser.GetBrowser(viper.GetString(envs.CORE_BROWSER), cfg.Private)
}

func getURL(profile string) string {
	url := lacct.GetAccountURL(profile)
	if url != "" {
		return url
	}

	url = viper.GetString(envs.SESSION_URL)
	if url != "" {
		return url
	}

	fmt.Printf("enter your AWS access portal URL: ")
	reader := bufio.NewReader(os.Stdin)
	url, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("An error occurred while reading input: ", err)
	}
	url = strings.TrimSuffix(url, "\n")
	viper.Set("session.url", url)

	if err := viper.WriteConfig(); err != nil {
		log.Fatal("could not write to config file:", err)
	}
	return url
}
