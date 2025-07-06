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

	lacct "github.com/louislef299/aws-sso/internal/account"
	"github.com/louislef299/aws-sso/internal/envs"
	laws "github.com/louislef299/aws-sso/pkg/aws"
	lconfig "github.com/louislef299/aws-sso/pkg/config"
	"github.com/louislef299/aws-sso/pkg/dlogin"
	los "github.com/louislef299/aws-sso/pkg/os"
	pecr "github.com/louislef299/aws-sso/plugins/aws/ecr"
	peks "github.com/louislef299/aws-sso/plugins/aws/eks"
	poidc "github.com/louislef299/aws-sso/plugins/aws/oidc"
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
				private = lacct.GetAccountPrivate(requestProfile)
			}

			if token != "" {
				setToken(token)
			} else if t := lacct.GetAccountToken(requestProfile); t != "" {
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

		oidcCfg := &poidc.OIDCLogin{
			Profile:      requestProfile,
			Region:       region,
			NewProfile:   newAuth,
			Private:      private,
			Refresh:      refresh,
			SkipDefaults: skipDefaults,
		}
		if err = dlogin.DLogin(cmd.Context(), "oidc", oidcCfg); err != nil {
			panic(err)
		}

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
				Config:   oidcCfg.Config,
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
					Config:       oidcCfg.Config,
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
