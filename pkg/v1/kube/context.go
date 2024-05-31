package kube

import "k8s.io/client-go/tools/clientcmd/api"

func (c *ClusterOptions) GetContext(namespace string) (*api.Context, error) {
	return &api.Context{
		Cluster:   c.Profile,
		AuthInfo:  c.Profile,
		Namespace: namespace,
	}, nil
}

func GetNamespace(c *api.Config) string {
	_, ok := c.Contexts[c.CurrentContext]
	if ok {
		return c.Contexts[c.CurrentContext].Namespace
	}
	return "default"
}
