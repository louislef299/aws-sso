package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/louislef299/aws-sso/internal/envs"
	laws "github.com/louislef299/aws-sso/pkg/aws"
	"github.com/louislef299/aws-sso/pkg/dlogin"
	pecr "github.com/louislef299/aws-sso/plugins/aws/ecr"
	poidc "github.com/louislef299/aws-sso/plugins/aws/oidc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var force, cleanToken bool

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logs you out of your SSO session",
	Long: `Removes the locally stored SSO tokens from the client-side 
cache, sends an API call to the IAM Identity Center service 
to invalidate the corresponding server-side IAM Identity 
Center sign in session, and removes the token locally.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !laws.IsProfileConfigured() && !force {
			log.Println("local profile not found, nothing to do")
			return
		}

		profile := laws.CurrentProfile()
		// Start up new config
		cfg, err := config.LoadDefaultConfig(cmd.Context(), config.WithRegion(region),
			config.WithSharedConfigProfile(profile))
		if err != nil {
			log.Fatal("couldn't load new config:", err)
		}

		// if session.profile is set, coming from a session
		if laws.IsProfileConfigured() {
			if err = dlogin.DLogout(cmd.Context(), "ecr", &pecr.ECRLogin{
				Username: "AWS",
				Config:   &cfg,
			}); err != nil {
				log.Fatal(err)
			}
		}

		// clean aws configs MUST GO LAST
		if err := dlogin.DLogout(cmd.Context(), "oidc", &poidc.OIDCLogin{
			Config:     &cfg,
			CleanToken: cleanToken,
		}); err != nil {
			log.Fatal("could not logout of AWS:", err)
		}

		if cleanToken {
			log.Println("cleaned out your old sso profiles")
		}

		// reset viper session configs(except for lastVersionCheck)
		sessionTree := viper.Sub(envs.SESSION_HEADER)
		lastVCheck, found := strings.CutPrefix(envs.SESSION_LAST_VCHECK, envs.SESSION_HEADER+".")
		if !found {
			log.Printf("couldn't find the session header prefix %s\n", envs.SESSION_HEADER+".")
		}
		lastVCheck = strings.ToLower(lastVCheck)

		for _, k := range sessionTree.AllKeys() {
			if k == lastVCheck {
				continue
			}
			viper.Set(fmt.Sprintf("%s.%s", envs.SESSION_HEADER, k), "")
		}
		err = viper.WriteConfig()
		if err != nil {
			log.Fatal("couldn't reset config values: ", err)
		}

		log.Println("successfully logged out of", profile)
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
	logoutCmd.Flags().BoolVarP(&cleanToken, "clean", "c", false, "clean out your current SSO cache")
	logoutCmd.Flags().BoolVarP(&force, "force", "f", false, "skip the safety checks and force a logout action")
}
