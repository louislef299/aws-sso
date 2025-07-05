package oidc

import (
	"context"

	"github.com/louislef299/aws-sso/pkg/dlogin"
	"github.com/spf13/cobra"
)

type OIDCLogin struct{}

func init() {
	dlogin.Register("oidc", &OIDCLogin{})
}

func (e *OIDCLogin) Init(cmd *cobra.Command) error {
	return nil
}

func (a *OIDCLogin) Login(ctx context.Context, config any) error {
	return nil
}

func (a *OIDCLogin) Logout(ctx context.Context, config any) error {
	return nil
}
