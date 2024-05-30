package kube

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

// Look at following link:
// https://github.com/kubernetes/client-go/blob/f457a57d6d2564ff06461d22ada3ff5ca6fec9c4/tools/clientcmd/config.go#L166
func ConfigureCluster(ctx context.Context, c *types.Cluster, optFns ...ClusterOptionsFunc) error {
	return Config(ctx, append(optFns, WithCluster(c))...)
}

func Config(ctx context.Context, optFns ...ClusterOptionsFunc) error {
	option, err := NewClusterOption()
	if err != nil {
		return err
	}

	for _, fn := range optFns {
		err := fn(option)
		if err != nil {
			return err
		}
	}

	kubeConfig, filepath, err := GetAPIConfig()
	if err != nil {
		return err
	}
	log.Println("setting kube config values for cluster", *option.Cluster.Arn)

	// must set this value to change the kube context
	kubeConfig.CurrentContext = option.Profile
	log.Println("current context is", option.Profile)
	cluster, err := option.GetCluster()
	if err != nil {
		return err
	}
	context, err := option.GetContext(GetNamespace(kubeConfig))
	if err != nil {
		return err
	}
	authInfo, err := option.GetAuthInfo()
	if err != nil {
		return err
	}

	// The name of the cluster must be the AWS cluster ARN, otherwise there will be config errors
	kubeConfig.Clusters[option.Profile] = cluster
	kubeConfig.Contexts[option.Profile] = context
	kubeConfig.AuthInfos[option.Profile] = authInfo

	return clientcmd.WriteToFile(*kubeConfig, filepath)
}

// Returns default configuration filepath for kubectl
func getDefaultConfig() string {
	return clientcmd.NewDefaultPathOptions().GetDefaultFilename()
}

// Validates that config file exists, otherwise configures new one
func readConfig(filepath string) (*api.Config, error) {
	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		newConf := api.NewConfig()
		err = clientcmd.WriteToFile(*newConf, filepath)
		if err != nil {
			return nil, err
		}

		return newConf, nil
	} else if err != nil {
		return nil, err
	}
	return clientcmd.LoadFromFile(filepath)
}

func getKubeConfig(filepath string) (*api.Config, error) {
	return readConfig(filepath)
}

func GetAPIConfig() (*api.Config, string, error) {
	filepath := getDefaultConfig()
	api, err := getKubeConfig(filepath)
	if err != nil {
		return nil, "", err
	}
	return api, filepath, nil
}
