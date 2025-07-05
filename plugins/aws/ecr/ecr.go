package ecr

import (
	"context"
	"fmt"
	"log"

	lcmd "github.com/louislef299/aws-sso/pkg/config"
	"github.com/louislef299/aws-sso/pkg/dlogin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const ECR_DISABLE_ECR_LOGIN = "ecr.disableECRLogin"

type ECRLogin struct{}

func init() {
	dlogin.Register("ecr", &ECRLogin{})
}

func (e *ECRLogin) Init(cmd *cobra.Command) error {
	cmd.Flags().Bool("disableECRLogin", true, "Disables automatic detection and login for ECR")
	lcmd.AddConfigValue(ECR_DISABLE_ECR_LOGIN, "Disables automatic detection and login for ECR")

	return viper.BindPFlag(ECR_DISABLE_ECR_LOGIN, cmd.Flags().Lookup("disableECRLogin"))
}

func (a *ECRLogin) Login(ctx context.Context, config any, opts ...dlogin.ConfigOptionsFunc) error {
	_, ok := config.(*dlogin.Config)
	if !ok {
		return fmt.Errorf("expected dlogin.Config, got %T", config)
	}

	if viper.GetBool(ECR_DISABLE_ECR_LOGIN) {
		log.Println("ECR Plugin is disabled, skipping...")
		return nil
	}

	log.Println("logging into ECR!")
	return nil
}
