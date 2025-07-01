package kube

import (
	"encoding/base64"
	"fmt"

	"k8s.io/client-go/tools/clientcmd/api"
)

func (c *ClusterOptions) GetCluster() (*api.Cluster, error) {
	data, err := base64.StdEncoding.DecodeString(*c.Cluster.CertificateAuthority.Data)
	if err != nil {
		return nil, fmt.Errorf("could not decode certificate: %v", err)
	}

	return &api.Cluster{
		LocationOfOrigin:         *c.Cluster.Endpoint,
		Server:                   *c.Cluster.Endpoint,
		CertificateAuthorityData: data,
	}, nil
}
