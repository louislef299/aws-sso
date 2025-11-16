package eks

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"slices"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/louislef299/aws-sso/internal/envs"
	laws "github.com/louislef299/aws-sso/pkg/aws"
	lconfig "github.com/louislef299/aws-sso/pkg/config"
	"github.com/louislef299/aws-sso/pkg/dlogin"
	lk8s "github.com/louislef299/aws-sso/pkg/kube"
	"github.com/louislef299/aws-sso/pkg/prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const EKS_DISABLE_EKS_LOGIN = "eks.disableEKSLogin"

var (
	ErrClusterDoesNotExist = errors.New("the provided cluster does not exist in this environment")
	ErrClustersDoNotExist  = errors.New("no clusters in this environment")
)

type EKSLogin struct {
	Profile      string
	Cluster      string
	Region       string
	Command      *cobra.Command
	Config       *aws.Config
	SkipDefaults bool
}

func init() {
	dlogin.Register("eks", &EKSLogin{})
}

func (e *EKSLogin) Init(cmd *cobra.Command) error {
	err := dlogin.Activate("eks")
	if err != nil {
		return err
	}

	cmd.Flags().Bool("disableEKSLogin", false, "Disables automatic detection and login for EKS")
	lconfig.AddConfigValue(EKS_DISABLE_EKS_LOGIN, "Disables automatic detection and login for EKS")

	return viper.BindPFlag(EKS_DISABLE_EKS_LOGIN, cmd.Flags().Lookup("disableEKSLogin"))
}

func (a *EKSLogin) Login(ctx context.Context, config any) error {
	cfg, ok := config.(*EKSLogin)
	if !ok {
		return fmt.Errorf("expected EKSLogin, got %T", config)
	}

	if viper.GetBool(EKS_DISABLE_EKS_LOGIN) {
		log.Println("EKS Plugin is disabled, skipping...")
		return nil
	}

	// Configure cluster options before gathering token
	options := []lk8s.ClusterOptionsFunc{
		lk8s.WithProfile(cfg.Profile),
		lk8s.WithRegion(cfg.Region),
	}
	log.Println("set k8s profile to", cfg.Profile)
	imp, err := cfg.Command.Flags().GetString("as")
	if err != nil {
		log.Println(err)
	}
	impG, err := cfg.Command.Flags().GetStringArray("as-group")
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

	cluster, err := getClusterName(ctx, cfg)
	if err == ErrClustersDoNotExist {
		log.Println("there were no clusters found in this environment. skipping Kubernetes EKS and Docker ECR configuration")
		return nil
	} else if err == ErrClusterDoesNotExist {
		log.Printf("cluster %s wasn't found, please select one of the provided clusters:", cluster)
		clusters, err := laws.GetClusters(ctx, cfg.Config)
		if err != nil {
			log.Fatal("couldn't gather clusters: ", err)
		}
		cluster = fuzzyCluster(clusters)
	} else if err != nil {
		log.Printf("could not gather cluster information: %v\n", err)
		return nil
	}
	log.Println("using cluster", cluster)

	return loginEKS(ctx, cfg.Config, cluster, options...)
}

func (a *EKSLogin) Logout(ctx context.Context, config any) error {
	_, ok := config.(*EKSLogin)
	if !ok {
		return fmt.Errorf("expected EKSLogin, got %T", config)
	}

	if viper.GetBool(EKS_DISABLE_EKS_LOGIN) {
		log.Println("EKS Plugin is disabled, skipping...")
		return nil
	}

	return nil
}

func fuzzyCluster(clusters []string) string {
	indexChoice, _ := prompt.Select("Select your cluster", clusters, prompt.FuzzySearchWithPrefixAnchor(clusters))
	log.Printf("Selected cluster %s", clusters[indexChoice])
	return clusters[indexChoice]
}

func getClusterName(ctx context.Context, cfg *EKSLogin) (string, error) {
	clusters, err := laws.GetClusters(ctx, cfg.Config)
	if err != nil {
		return "", err
	}

	if len(clusters) == 0 {
		return "", ErrClustersDoNotExist
	}

	if cfg.SkipDefaults {
		return fuzzyCluster(clusters), nil
	}

	var cluster string
	if cfg.Cluster != "" {
		cluster = cfg.Cluster
	} else if defaultCluster := viper.GetString(envs.CORE_DEFAULT_CLUSTER); defaultCluster != "" {
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
	} else if c := viper.GetString(envs.SESSION_CLUSTER); c != "" {
		cluster = c
	} else {
		cluster = fuzzyCluster(clusters)
	}

	if !slices.Contains(clusters, cluster) {
		return cluster, ErrClusterDoesNotExist
	}
	return cluster, nil
}

func loginEKS(ctx context.Context, cfg *aws.Config, cluster string, optFns ...lk8s.ClusterOptionsFunc) error {
	// configure kubernetes credentials
	clusterInfo, err := laws.GetClusterInfo(ctx, cfg, cluster)
	if err != nil {
		return fmt.Errorf("couldn't gather cluster information: %v", err)
	}

	log.Printf("configuring kubernetes configuration cluster access for %s\n", cluster)
	err = lk8s.ConfigureCluster(ctx, clusterInfo.Cluster, optFns...)
	if err != nil {
		return fmt.Errorf("could not update kubeconfig: %v", err)
	}
	viper.Set(envs.SESSION_CLUSTER, *clusterInfo.Cluster.Name)
	err = viper.WriteConfig()
	if err != nil {
		return fmt.Errorf("could not write cluster information to config: %v", err)
	}
	return nil
}
