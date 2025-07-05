package eks

import (
	"context"

	lconfig "github.com/louislef299/aws-sso/pkg/config"
	"github.com/louislef299/aws-sso/pkg/dlogin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const EKS_DISABLE_EKS_LOGIN = "eks.disableEKSLogin"

type EKSLogin struct{}

func init() {
	dlogin.Register("eks", &EKSLogin{})
}

func (e *EKSLogin) Init(cmd *cobra.Command) error {
	cmd.Flags().Bool("disableEKSLogin", false, "Disables automatic detection and login for EKS")
	lconfig.AddConfigValue(EKS_DISABLE_EKS_LOGIN, "Disables automatic detection and login for EKS")

	return viper.BindPFlag(EKS_DISABLE_EKS_LOGIN, cmd.Flags().Lookup("disableEKSLogin"))
}

func (a *EKSLogin) Login(ctx context.Context, config any) error {
	return nil
}

func (a *EKSLogin) Logout(ctx context.Context, config any) error {
	return nil
}
