package oidc

import (
	"context"
	"fmt"

	"github.com/louislef299/aws-sso/pkg/dlogin"
)

type AWSLogin struct{}

func init() {
	dlogin.Register("aws", &AWSLogin{})
}

func (a *AWSLogin) Login(ctx context.Context, config any, opts ...dlogin.ConfigOptionsFunc) error {
	c, ok := config.(*dlogin.Config)
	if !ok {
		return fmt.Errorf("expected dlogin.Config, got %T", config)
	}

	c.Role = "aws:arn:louis"
	c.Secret = "secret"
	return nil
}
