package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/louislef299/aws-sso/internal/browser"
	"github.com/louislef299/aws-sso/internal/envs"
	lregion "github.com/louislef299/aws-sso/internal/region"
	laws "github.com/louislef299/aws-sso/pkg/aws"
	lconfig "github.com/louislef299/aws-sso/pkg/config"
	"github.com/louislef299/aws-sso/pkg/dlogin"
	los "github.com/louislef299/aws-sso/pkg/os"
	pecr "github.com/louislef299/aws-sso/plugins/aws/ecr"
	peks "github.com/louislef299/aws-sso/plugins/aws/eks"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	role, startUrl, output, token  string
	clusterName                    string
	disableEKSLogin                bool
	private, refresh, skipDefaults bool

	ErrKeyDoesNotExist = errors.New("the provided key doesn't exist")
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:     "login",
	Example: "  aws-sso login env1",
	Short:   "Retrieve short-lived credentials via AWS SSO & SSOOIDC",
	Args:    cobra.MaximumNArgs(1),
	Long: `Creates and returns an access token for the authorized client. 
The access token issued will be used to fetch short-term 
credentials for the assigned roles in the AWS account.

If the account has an EKS cluster, authenticates with
the cluster and logs you into you ECR in your account.
EKS and ECR auth can be disabled with configuration
updates.`,
	Run: func(cmd *cobra.Command, args []string) {
		var requestProfile string
		var err error
		// find out if an account profile is being requested
		if len(args) == 1 {
			requestProfile = args[0]
		}

		newAuth := false
		if requestProfile == "" {
			fmt.Printf("please enter a prefix alias for this context(ex: env1): ")
			reader := bufio.NewReader(os.Stdin)
			alias, err := reader.ReadString('\n')
			if err != nil {
				log.Fatal("An error occurred while reading input: ", err)
			}
			requestProfile = strings.TrimSuffix(alias, "\n")
			newAuth = true
		}

		// check if referencing a local profile
		lc, err := laws.IsLocalConfig(requestProfile)
		if err != nil {
			log.Println("couldn't find predefined AWS configurations:", err)
		}

		var profileToSet string
		// if not a local profile, check for account information
		if !lc {
			if !private {
				private = getAccountPrivate(requestProfile)
			}

			if token != "" {
				setToken(token)
			} else if t := getAccountToken(requestProfile); t != "" {
				setToken(t)
			} else {
				checkToken()
			}
			log.Println("using token", getCurrentToken())

			profileToSet = los.GetProfile(requestProfile)
		} else {
			profileToSet = requestProfile
		}

		region, err = syncAccountRegionSession(requestProfile, region)
		if err != nil {
			log.Fatal("could not sync account session:", err)
		}

		cfg, err := getAWSConfig(cmd.Context(), requestProfile, region, newAuth)
		if err != nil {
			log.Fatal("could not generate AWS config: ", err)
		}
		region = cfg.Region

		wg := sync.WaitGroup{}
		// configure docker credentials
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Println("configuring local docker credentials with ECR token")

			ctx, cancel := context.WithTimeout(cmd.Context(), commandTimeout)
			defer cancel()
			if err = dlogin.DLogin(ctx, "ecr", &pecr.ECRLogin{
				Username: "AWS",
				Config:   cfg,
			}); err != nil {
				panic(err)
			}
		}()

		if !viper.GetBool(envs.CORE_DISABLE_EKS_LOGIN) && !disableEKSLogin {
			wg.Add(1)
			go func() {
				defer wg.Done()

				ctx, cancel := context.WithTimeout(cmd.Context(), commandTimeout)
				defer cancel()
				if err = dlogin.DLogin(ctx, "eks", &peks.EKSLogin{
					Cluster:      clusterName,
					Profile:      profileToSet,
					Region:       region,
					Command:      cmd,
					Config:       cfg,
					SkipDefaults: skipDefaults,
				}); err != nil {
					panic(err)
				}
			}()
		}

		wg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringVarP(&region, "region", "r", "", "The region you would like to use at login")
	lconfig.BindConfigValue(envs.SESSION_REGION, loginCmd.Flags().Lookup("region"))

	loginCmd.Flags().StringVarP(&startUrl, "url", "u", "", "The AWS SSO start url")
	lconfig.BindConfigValue(envs.SESSION_URL, loginCmd.Flags().Lookup("url"))

	loginCmd.Flags().StringVar(&role, "role", "", "The IAM role to use when logging in")
	lconfig.BindConfigValue(envs.SESSION_ROLE, loginCmd.Flags().Lookup("role"))

	loginCmd.Flags().BoolVarP(&private, "private", "p", false, "Open a private browser when gathering/refreshing token")
	loginCmd.Flags().BoolVar(&refresh, "refresh", false, "Whether to manually refresh your local authentication token")
	loginCmd.Flags().BoolVar(&skipDefaults, "skipDefaults", false, "Skip the default login values and use prompt selection")
	loginCmd.Flags().StringVarP(&token, "token", "t", "", "The token to use when logging in. To be used when managing multiple session tokens at once (shorthand '-' for default token)")

	loginCmd.Flags().StringVarP(&clusterName, "cluster", "c", "", "The cluster you would like to target when logging in")
	loginCmd.Flags().StringVarP(&output, "output", "o", "json", "The output format for sso")
}

func getAWSConfig(ctx context.Context, profile, region string, newProfile bool) (*aws.Config, error) {
	// check if referencing a local profile
	lc, err := laws.IsLocalConfig(profile)
	if err != nil {
		log.Println("couldn't find predefined AWS configurations:", err)
	}

	if lc {
		log.Println("using existing configuration profile", profile)
		cfg, err := config.LoadDefaultConfig(ctx,
			config.WithSharedConfigProfile(profile),
		)
		if err != nil {
			return nil, err
		}
		deepSet(envs.SESSION_PROFILE, profile)

		return &cfg, nil
	}

	ssoRegion, err := lregion.GetRegion(lregion.STS)
	if err != nil {
		log.Fatal("could not gather sso region:", err)
	}

	// gross anti-pattern, but too lazy to reprogram the existing logic for now
	if region == "cn-north-1" || region == "cn-northwest-1" {
		ssoRegion = "cn-north-1"
	}
	log.Println("using sso region", ssoRegion, "to login")

	acctID := getAccountID(profile)
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(ssoRegion))
	if err != nil {
		return nil, err
	}

	p, err := loginAWS(ctx, cfg, acctID, profile, newProfile)
	if err != nil {
		log.Fatal("couldn't log into AWS: ", err)
	}
	deepSet(envs.SESSION_PROFILE, p)

	log.Println("loading up new config", p, "with region", region)
	// Start up new config with newly configured profile
	cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(region), config.WithSharedConfigProfile(p))
	if err != nil {
		log.Fatal("couldn't load new config:", err)
	}

	return &cfg, nil
}

func getBrowser() browser.Browser {
	if private {
		log.Println("browser set to open incognito(no cookies)")
	} else {
		log.Println("browser set to default(use cookies)")
	}
	return browser.GetBrowser(viper.GetString(envs.CORE_BROWSER), private)
}

func getURL(profile string) string {
	url := getAccountURL(profile)
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

func loginAWS(ctx context.Context, cfg aws.Config, acctID, profile string, newProfile bool) (string, error) {
	u := getURL(profile)
	deepSet(envs.SESSION_URL, u)

	clientInfo, err := laws.GatherClientInformation(ctx, &cfg, u, getBrowser(), refresh)
	if err != nil {
		return "", err
	}

	account, err := laws.RetrieveAccountInformation(ctx, &cfg, &clientInfo.AccessToken, &acctID)
	if err != nil {
		return "", err
	}
	acctID = *account.AccountId

	role, err := laws.RetrieveRoleInfo(ctx, &cfg, *account.AccountId, clientInfo.AccessToken, skipDefaults)
	if err != nil {
		return "", err
	}
	log.Println("using aws role", *role.RoleName)
	deepSet(envs.SESSION_ROLE, *role.RoleName)

	err = laws.SaveUsageInformation(account, &role)
	if err != nil {
		return "", err
	}

	// set the new profile in account config
	if newProfile {
		err = addAccount(profile, &Account{
			ID:      acctID,
			Region:  cfg.Region,
			Private: private,
			Token:   getCurrentToken(),
		})
		if err != nil {
			log.Println("WARNING: couldn't write to configuration file:", err)
		}
	}
	return laws.GetAndSaveRoleCredentials(ctx, &cfg, account.AccountId, role.RoleName, &clientInfo.AccessToken, profile, cfg.Region)
}
