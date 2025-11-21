package plugins

import (
	"github.com/louislef299/knot/pkg/dlogin"
	"github.com/louislef299/knot/plugins/aws/ecr"
	"github.com/louislef299/knot/plugins/aws/eks"
	"github.com/louislef299/knot/plugins/aws/oidc"
)

func GetAvailablePlugins() map[string]dlogin.ILogin {
	return map[string]dlogin.ILogin{
		"oidc": &oidc.OIDCLogin{},
		"eks":  &eks.EKSLogin{},
		"ecr":  &ecr.ECRLogin{},
	}
}
