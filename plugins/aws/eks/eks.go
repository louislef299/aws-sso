package eks

import (
	"context"
	"fmt"

	"github.com/louislef299/aws-sso/pkg/dlogin"
)

type K8sLogin struct{}

func init() {
	dlogin.Register("k8s", &K8sLogin{})
}

func (a *K8sLogin) Login(ctx context.Context, config any, opts ...dlogin.ConfigOptionsFunc) error {
	c, ok := config.(*dlogin.Config)
	if !ok {
		return fmt.Errorf("expected dlogin.Config, got %T", config)
	}

	fmt.Printf("Role: %s\tSecret: %s\n", c.Role, c.Secret)
	c.Cluster = "louis-dev"
	return nil
}
