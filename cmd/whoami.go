/*
Copyright © 2023 Louis Lefebvre <louislefebvre1999@gmail.com>
*/
package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/louislef299/aws-sso/internal/envs"
	lregion "github.com/louislef299/aws-sso/internal/region"
	laws "github.com/louislef299/aws-sso/pkg/aws"
	lk8s "github.com/louislef299/aws-sso/pkg/kube"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// whoamiCmd represents the whoami command
var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Check your current AWS STS settings",
	Long: `Returns the current working environment. Uses 
the sts package to get the called identity
information. If in an EKS environment, will
gather cluster information.

Similar to running:
aws sts get-caller-identity`,
	Run: func(cmd *cobra.Command, args []string) {
		if !viper.IsSet(envs.SESSION_PROFILE) || !laws.IsProfileConfigured() {
			log.Println("not currently signed in")
			return
		} else {
			fmt.Println("using profile", laws.CurrentProfile())
		}

		region, err := lregion.GetRegion(lregion.EKS)
		if err != nil {
			log.Fatal("couldn't find region:", err)
		}

		cfg, err := config.LoadDefaultConfig(cmd.Context(), config.WithRegion(region), config.WithSharedConfigProfile(laws.CurrentProfile()))
		if err != nil {
			log.Fatal("couldn't load new config:", err)
		}

		// Set a quick timeout for caller identity
		ctx, cancel := context.WithTimeout(cmd.Context(), commandTimeout)
		defer cancel()

		callerID, err := laws.GetCallerIdentity(ctx, &cfg)
		if err != nil {
			log.Fatal("couldn't gather sts identity information: ", err)
		}
		fmt.Printf("AWS Information:\n\tUser ID: %s\n\tAccount: %s\n\tCaller ARN: %s\n", *callerID.UserId, *callerID.Account, *callerID.Arn)

		if !viper.GetBool(envs.CORE_DISABLE_EKS_LOGIN) {
			cluster := viper.GetString(envs.SESSION_CLUSTER)
			if cluster == "" {
				log.Println("Kubernetes not configured locally")
			}

			err = printClusterInfo(cmd.Context(), &cfg, cluster, os.Stdout)
			if err != nil {
				log.Fatal("couldn't print cluster information:", err)
			}

			kubeapi, filepath, err := lk8s.GetAPIConfig()
			if err != nil {
				log.Fatal("couldn't gather kube config information:", err)
			}
			cc := kubeapi.CurrentContext
			fmt.Printf("Local Kube Config:\n\tUsing Config File: %s\n\tCurrent Context: %s\n",
				filepath,
				cc,
			)
			if i := kubeapi.AuthInfos[cc].Impersonate; i != "" {
				fmt.Printf("\tImpersonating User: %s\n", i)
			}
			if g := kubeapi.AuthInfos[cc].ImpersonateGroups; g != nil {
				fmt.Printf("\tImpersonating Groups: %v\n", g)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}

func printClusterInfo(ctx context.Context, cfg *aws.Config, cluster string, out io.Writer) error {
	resp, err := laws.GetClusterInfo(ctx, cfg, cluster)
	if err != nil {
		return fmt.Errorf("could not gather main cluster information: %v", err)
	}
	_, err = fmt.Fprintf(out, "Kubernetes control plane is running at %s(%v)\n\tCluster Name: %v\n\tRole Arn: %v\n\tPlatform Version: %v\n",
		*resp.Cluster.Endpoint,
		*resp.Cluster.KubernetesNetworkConfig.ServiceIpv4Cidr,
		*resp.Cluster.Name,
		*resp.Cluster.RoleArn,
		*resp.Cluster.PlatformVersion,
	)
	return err
}
