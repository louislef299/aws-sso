//go:build !kube
// +build !kube

package cmd

import (
	"k8s.io/cli-runtime/pkg/genericclioptions"
	utilpointer "k8s.io/utils/ptr"
)

func init() {
	// add impersonation flags here
	impersonateGroup := []string{}
	kc := &genericclioptions.ConfigFlags{
		Impersonate:      utilpointer.To[string](""),
		ImpersonateGroup: &impersonateGroup,
	}
	kc.AddFlags(rootCmd.PersistentFlags())
}
