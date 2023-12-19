//go:build kube
// +build kube

package cmd

import (
	"os"

	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// Add kube flags to cli
func init() {
	// Create the set of flags for your kubectl-aws-sso
	flags := pflag.NewFlagSet("kubectl-aws-sso", pflag.ExitOnError)
	pflag.CommandLine = flags

	// Configure the genericclioptions
	streams := genericclioptions.IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}

	// kubectl configuration flags
	kubeConfigFlags := genericclioptions.NewConfigFlags(false)

	// Join all flags to your root command
	flags.AddFlagSet(rootCmd.PersistentFlags())
	kubeConfigFlags.AddFlags(flags)

	rootCmd.SetOutput(streams.Out)
}
