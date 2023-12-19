package aws

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
)

// Returns a list of the clusters in the environment
func GetClusters(ctx context.Context, cfg *aws.Config) ([]string, error) {
	svc := eks.NewFromConfig(*cfg)
	resp, err := svc.ListClusters(ctx, &eks.ListClustersInput{})
	if err != nil {
		return nil, fmt.Errorf("there was an issue getting the clusters: %v", err)
	}
	return resp.Clusters, nil
}

// Get cluster information for provided cluster
func GetClusterInfo(ctx context.Context, cfg *aws.Config, cluster string) (*eks.DescribeClusterOutput, error) {
	return eks.NewFromConfig(*cfg).DescribeCluster(ctx, &eks.DescribeClusterInput{
		Name: aws.String(cluster),
	})
}

func PrintClusterInfo(ctx context.Context, cfg *aws.Config, cluster string, out io.Writer) error {
	resp, err := GetClusterInfo(ctx, cfg, cluster)
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
