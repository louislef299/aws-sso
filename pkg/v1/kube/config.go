package kube

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/eks"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

// Look at following link:
// https://github.com/kubernetes/client-go/blob/f457a57d6d2564ff06461d22ada3ff5ca6fec9c4/tools/clientcmd/config.go#L166
func ConfigureCluster(info *eks.DescribeClusterOutput, region, profile string) error {
	data, err := base64.StdEncoding.DecodeString(*info.Cluster.CertificateAuthority.Data)
	if err != nil {
		return fmt.Errorf("could not decode certificate: %v", err)
	}

	filepath := GetDefaultConfig()
	kubeConfig, err := readConfig(filepath)
	if err != nil {
		return fmt.Errorf("could not read in kubectl config: %v", err)
	}

	log.Println("setting kube config values for cluster", *info.Cluster.Arn)

	// The name of the cluster must be the AWS cluster ARN, otherwise there will be config errors
	kubeConfig.Clusters[profile] = &api.Cluster{
		LocationOfOrigin:         *info.Cluster.Endpoint,
		Server:                   *info.Cluster.Endpoint,
		CertificateAuthorityData: data,
	}
	kubeConfig.Contexts[profile] = &api.Context{
		Cluster:  profile,
		AuthInfo: profile,
	}
	kubeConfig.AuthInfos[profile] = &api.AuthInfo{
		Exec: &api.ExecConfig{
			APIVersion: "client.authentication.k8s.io/v1beta1",
			Command:    "aws",
			Args: []string{
				"--region",
				region,
				"eks",
				"get-token",
				"--cluster-name",
				*info.Cluster.Name,
			},
			Env: []api.ExecEnvVar{
				{
					Name:  "AWS_PROFILE",
					Value: profile,
				},
			},
			ProvideClusterInfo: true,
		},
	}

	// must set this value to change the kube context
	kubeConfig.CurrentContext = profile

	return clientcmd.WriteToFile(*kubeConfig, filepath)
}

// Returns default configuration filepath for kubectl
func GetDefaultConfig() string {
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
