/*
Copyright © 2022 Louis Lefebvre <louislefebvre1999@gmail.com>
*/
package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/louislef299/aws-sso/internal/browser"
	. "github.com/louislef299/aws-sso/internal/envs"
	lregion "github.com/louislef299/aws-sso/internal/region"
	laws "github.com/louislef299/aws-sso/pkg/v1/aws"
	ldocker "github.com/louislef299/aws-sso/pkg/v1/docker"
	lk8s "github.com/louislef299/aws-sso/pkg/v1/kube"
	los "github.com/louislef299/aws-sso/pkg/v1/os"
	"github.com/louislef299/aws-sso/pkg/v1/prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	role, startUrl, output, token    string
	clusterName                      string
	disableEKSLogin, disableECRLogin bool
	private, refresh, skipDefaults   bool

	ErrKeyDoesNotExist     = errors.New("the provided key doesn't exist")
	ErrClusterDoesNotExist = errors.New("the provided cluster does not exist in this environment")
	ErrClustersDoNotExist  = errors.New("no clusters in this environment")
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

		region, err = syncAccountRegionSession(requestProfile, region)
		if err != nil {
			log.Fatal("could not sync account session:", err)
		}

		// Configure cluster options before gathering token
		options := []lk8s.ClusterOptionsFunc{
			lk8s.WithProfile(los.GetProfile(requestProfile)),
			lk8s.WithRegion(region),
		}
		log.Println("Set profile to", los.GetProfile(requestProfile))
		imp, err := cmd.Flags().GetString("as")
		if err != nil {
			log.Println(err)
		}
		impG, err := cmd.Flags().GetStringArray("as-group")
		if err != nil {
			log.Println(err)
		}
		// check for impersonation flags(kube or aws)
		if imp != "" && len(impG) > 0 {
			log.Printf("impersonating user %s in group %s\n", imp, impG)
			options = append(options, lk8s.WithImpersonation(imp, impG))
		} else if imp != "" && len(impG) <= 0 {
			log.Fatal("when impersonating, must provide both a Username and a Group(use --as & --as-group)")
		}

		cfg, err := getAWSConfig(cmd.Context(), requestProfile, newAuth)
		if err != nil {
			log.Fatal("could not generate AWS config: ", err)
		}

		wg := sync.WaitGroup{}
		if !viper.GetBool(CORE_DISABLE_ECR_LOGIN) && !disableECRLogin {
			// configure docker credentials
			wg.Add(1)
			go func() {
				defer wg.Done()
				log.Println("configuring local docker credentials with ECR token")
				ctx, cancel := context.WithTimeout(cmd.Context(), commandTimeout)
				defer cancel()

				ecrToken, ecrEndpoint, err := laws.GetECRToken(ctx, cfg)
				if err != nil {
					log.Fatal("couldn't gather ecr token:", err)
				}

				err = ldocker.Login("AWS", ecrToken, ecrEndpoint)
				if err != nil {
					log.Fatalf("could not log docker into ecr endpoint %s: %v", ecrEndpoint, err)
				}
			}()
		}

		if !viper.GetBool(CORE_DISABLE_EKS_LOGIN) && !disableEKSLogin {
			wg.Add(1)
			go func() {
				defer wg.Done()

				cluster, err := getClusterName(cmd.Context(), cfg, false)
				if err == ErrClustersDoNotExist {
					log.Println("there were no clusters found in this environment. skipping Kubernetes EKS and Docker ECR configuration")
					return
				} else if err == ErrClusterDoesNotExist {
					log.Printf("cluster %s wasn't found, please select one of the provided clusters:", cluster)
					clusters, err := laws.GetClusters(cmd.Context(), cfg)
					if err != nil {
						log.Fatal("couldn't gather clusters: ", err)
					}
					cluster = fuzzyCluster(clusters)
				} else if err != nil {
					log.Printf("could not gather cluster information: %v\n", err)
					return
				}
				log.Println("using cluster", cluster)

				loginEKS(cmd.Context(), *cfg, cluster, options...)
			}()
		}

		wg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringVarP(&region, "region", "r", "", "The region you would like to use at login")
	BindConfigValue(SESSION_REGION, loginCmd.Flags().Lookup("region"))

	loginCmd.Flags().StringVarP(&startUrl, "url", "u", "", "The AWS SSO start url")
	BindConfigValue(SESSION_URL, loginCmd.Flags().Lookup("url"))

	loginCmd.Flags().StringVar(&role, "role", "", "The IAM role to use when logging in")
	BindConfigValue(SESSION_ROLE, loginCmd.Flags().Lookup("role"))

	loginCmd.Flags().BoolVar(&disableEKSLogin, "disableEKSLogin", false, "Disables automatic detection and login for EKS")
	loginCmd.Flags().BoolVar(&disableECRLogin, "disableECRLogin", true, "Disables automatic detection and login for ECR")
	loginCmd.Flags().BoolVarP(&private, "private", "p", false, "Open a private browser when gathering/refreshing token")
	loginCmd.Flags().BoolVar(&refresh, "refresh", false, "Whether to manually refresh your local authentication token")
	loginCmd.Flags().BoolVar(&skipDefaults, "skipDefaults", false, "Skip the default login values and use prompt selection")
	loginCmd.Flags().StringVarP(&token, "token", "t", "", "The token to use when logging in. To be used when managing multiple session tokens at once (shorthand '-' for default token)")

	loginCmd.Flags().StringVarP(&clusterName, "cluster", "c", "", "The cluster you would like to target when logging in")
	loginCmd.Flags().StringVarP(&output, "output", "o", "json", "The output format for sso")
}

func fuzzyCluster(clusters []string) string {
	indexChoice, _ := prompt.Select("Select your cluster", clusters, prompt.FuzzySearchWithPrefixAnchor(clusters))
	log.Printf("Selected cluster %s", clusters[indexChoice])
	return clusters[indexChoice]
}

func getAWSConfig(ctx context.Context, profile string, newProfile bool) (*aws.Config, error) {
	// check if referencing a local profile
	lc, err := laws.IsLocalConfig(profile)
	if err != nil {
		log.Println("couldn't find predefined AWS configurations:", err)
	}

	if lc {
		log.Println("using existing configuration profile", profile)
		cfg, err := config.LoadDefaultConfig(ctx,
			config.WithRegion(region),
			config.WithSharedConfigProfile(profile),
		)
		if err != nil {
			return nil, err
		}
		deepSet(SESSION_PROFILE, profile)

		return &cfg, nil
	}

	ssoRegion, err := lregion.GetRegion(lregion.STS)
	if err != nil {
		log.Fatal("could not gather sso region:", err)
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
	deepSet(SESSION_PROFILE, p)

	log.Println("loading up new config", p, "with region", region)
	// Start up new config with newly configured profile
	cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(region), config.WithSharedConfigProfile(p))
	if err != nil {
		log.Fatal("couldn't load new config:", err)
	}

	return &cfg, nil
}

func getClusterName(ctx context.Context, cfg *aws.Config, skipFuzzy bool) (string, error) {
	clusters, err := laws.GetClusters(ctx, cfg)
	if err != nil {
		return "", err
	}

	if len(clusters) == 0 {
		return "", ErrClustersDoNotExist
	}

	if skipDefaults {
		return fuzzyCluster(clusters), nil
	}

	var cluster string
	if clusterName != "" {
		cluster = clusterName
	} else if defaultCluster := viper.GetString(CORE_DEFAULT_CLUSTER); defaultCluster != "" {
		r, err := regexp.Compile(defaultCluster)
		if err == nil {
			log.Printf("looking for cluster matching expression: %s\n", defaultCluster)
			for _, c := range clusters {
				if r.MatchString(c) {
					cluster = c
					break
				}
			}
		} else {
			log.Printf("assuming default static name due to regex failure: %v\n", err)
			cluster = defaultCluster
		}
	} else if len(clusters) == 1 {
		cluster = clusters[0]
	} else if len(clusters) == 0 {
		return "", ErrClustersDoNotExist
	} else if c := viper.GetString(SESSION_CLUSTER); c != "" {
		cluster = c
	} else if !skipFuzzy {
		cluster = fuzzyCluster(clusters)
	}

	if !slices.Contains(clusters, cluster) {
		return cluster, ErrClusterDoesNotExist
	}
	return cluster, nil
}

func getBrowser() browser.Browser {
	if private {
		log.Println("browser set to open incognito(no cookies)")
	} else {
		log.Println("browser set to default(use cookies)")
	}
	return browser.GetBrowser(viper.GetString(CORE_BROWSER), private)
}

func getURL() string {
	// validate there is a url endpoint
	url := viper.GetString(CORE_URL)
	if url != "" {
		return url
	}

	url = viper.GetString(SESSION_URL)
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
	u := getURL()
	deepSet(SESSION_URL, u)

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
	deepSet(SESSION_ROLE, *role.RoleName)

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

func loginEKS(ctx context.Context, cfg aws.Config, cluster string, optFns ...lk8s.ClusterOptionsFunc) {
	// configure kubernetes credentials
	clusterInfo, err := laws.GetClusterInfo(ctx, &cfg, cluster)
	if err != nil {
		log.Fatal("couldn't gather cluster information:", err)
	}

	log.Printf("configuring kubernetes configuration cluster access for %s\n", cluster)
	err = lk8s.ConfigureCluster(ctx, clusterInfo.Cluster, optFns...)
	if err != nil {
		log.Fatal("could not update kubeconfig: ", err)
	}
	viper.Set(SESSION_CLUSTER, *clusterInfo.Cluster.Name)
	err = viper.WriteConfig()
	if err != nil {
		log.Fatal("could not write cluster information to config:", err)
	}
}
