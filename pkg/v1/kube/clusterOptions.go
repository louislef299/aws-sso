package kube

import (
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/louislef299/aws-sso/internal/region"
)

type ClusterOptions struct {
	Cluster *types.Cluster
	Region  string
	Profile string

	impersonationEnabled bool
	Impersonate          string
	ImpersonateGroups    []string
}

type ClusterOptionsFunc func(*ClusterOptions) error

func NewClusterOption() (*ClusterOptions, error) {
	r, err := region.GetRegion(region.EKS)
	if err != nil {
		return nil, err
	}
	return &ClusterOptions{
		Region:  r,
		Profile: "default",
	}, nil
}

func WithCluster(c *types.Cluster) ClusterOptionsFunc {
	return func(o *ClusterOptions) error {
		o.Cluster = c
		return nil
	}
}

func WithImpersonation(user string, groups []string) ClusterOptionsFunc {
	return func(o *ClusterOptions) error {
		o.Impersonate = user
		o.ImpersonateGroups = groups
		o.impersonationEnabled = true
		return nil
	}
}

func WithProfile(p string) ClusterOptionsFunc {
	return func(o *ClusterOptions) error {
		o.Profile = p
		return nil
	}
}

func WithRegion(r string) ClusterOptionsFunc {
	return func(o *ClusterOptions) error {
		o.Region = r
		return nil
	}
}
