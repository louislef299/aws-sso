package version

import (
	"fmt"
	"io"

	"text/template"

	"github.com/spf13/cobra"
)

// Main version number being run right now.
var (
	Version    = "development"
	BuildOS    = "not_set"
	BuildArch  = "not_set"
	BuildTime  = "not_set"
	GoVersion  = "not_set"
	CommitHash = "not_set"
	Flavor     = "default"
)

// User Agent name set in requests.
const UserAgent = "aws-sso"

type CommandVersion struct {
	Name       string
	Short      string
	Version    string
	BuildTime  string
	BuildArch  string
	BuildOS    string
	GoVersion  string
	CommitHash string
}

func GetTemplate() string {
	return `{{printf "%s: %s/%s %s/%s built-with/%s\n build-time/%s commit-hash/%s\n" .Short .Name .Version .BuildOS .BuildArch .GoVersion .BuildTime .CommitHash}}`
}

// String prints the version of lnet
func String() string {
	return Version
}

// Returns a semantic version compliant string
func SemverString() string {
	return fmt.Sprintf("v%s", Version)
}

// Prints the version to the provided io.Writer
func PrintVersion(out io.Writer, cmd *cobra.Command) error {
	v := CommandVersion{
		Name:       cmd.Use,
		Short:      cmd.Short,
		Version:    Version,
		BuildTime:  BuildTime,
		BuildArch:  BuildArch,
		BuildOS:    BuildOS,
		GoVersion:  GoVersion,
		CommitHash: CommitHash,
	}
	tmpl, err := template.New("awsLoginPluginVersion").Parse(GetTemplate())
	if err != nil {
		return err
	}
	return tmpl.Execute(out, v)
}
