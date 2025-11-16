package plugins

import (
	"github.com/louislef299/aws-sso/pkg/dlogin"
	"github.com/louislef299/aws-sso/plugins/aws/ecr"
	"github.com/louislef299/aws-sso/plugins/aws/eks"
	"github.com/louislef299/aws-sso/plugins/aws/oidc"
)

func GetAvailablePlugins() map[string]dlogin.ILogin {
	return map[string]dlogin.ILogin{
		"oidc": &oidc.OIDCLogin{},
		"eks":  &eks.EKSLogin{},
		"ecr":  &ecr.ECRLogin{},
	}
}
