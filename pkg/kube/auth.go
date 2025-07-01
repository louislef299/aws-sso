package kube

import "k8s.io/client-go/tools/clientcmd/api"

func (c *ClusterOptions) GetAuthInfo() (*api.AuthInfo, error) {
	a := &api.AuthInfo{
		Exec: &api.ExecConfig{
			APIVersion: "client.authentication.k8s.io/v1beta1",
			Command:    "aws",
			Args: []string{
				"--region",
				c.Region,
				"eks",
				"get-token",
				"--cluster-name",
				*c.Cluster.Name,
			},
			Env: []api.ExecEnvVar{
				{
					Name:  "AWS_PROFILE",
					Value: c.Profile,
				},
			},
			ProvideClusterInfo: true,
		},
	}

	if c.impersonationEnabled {
		a.Impersonate = c.Impersonate
		a.ImpersonateGroups = c.ImpersonateGroups
	}
	return a, nil
}
