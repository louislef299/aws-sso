package ecr

import (
	"context"
	"fmt"

	"github.com/louislef299/aws-sso/pkg/dlogin"
)

type DockerLogin struct{}

func init() {
	dlogin.Register("docker", &DockerLogin{})
}

func (a *DockerLogin) Login(ctx context.Context, config any, opts ...dlogin.ConfigOptionsFunc) error {
	c, ok := config.(*dlogin.Config)
	if !ok {
		return fmt.Errorf("expected dlogin.Config, got %T", config)
	}

	fmt.Printf("Cluster: %s\n", c.Cluster)
	return nil
}
