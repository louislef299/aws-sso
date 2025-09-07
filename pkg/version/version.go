package version

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"text/template"

	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
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

	ErrNoValidVersionFound = errors.New("could not find a valid version")
)

const (
	// User Agent name set in requests.
	UserAgent = "aws-sso"

	// Represents the project URL to check for latest release
	releaseURL = "https://api.github.com/repos/louislef299/aws-sso/releases/latest"
)

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

// Represents the required fields for the GitHub API
// docs.github.com/en/rest/releases/releases?apiVersion=2022-11-28#get-the-latest-release
type latestVersion struct {
	Name    string `json:"name"`
	TagName string `json:"tag_name"`
	URL     string `json:"html_url"`
}

func CheckForUpdate() error {
	req, err := http.NewRequest(http.MethodGet, releaseURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	c := http.Client{
		Timeout: time.Second * 2,
	}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var release *latestVersion
	err = json.Unmarshal(body, &release)
	if err != nil {
		return err
	}

	releaseVersion, err := findValidVersion(release.TagName, release.Name)
	if err != nil {
		return err
	}

	r := semver.Compare("v"+Version, releaseVersion)
	if r < 0 {
		log.Printf("A new version of aws-sso is available(%s)!\n%s\n\n", releaseVersion, release.URL)
	} else {
		log.Printf("version looks good!\n\n")
	}
	return nil
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

func findValidVersion(versions ...string) (string, error) {
	for _, v := range versions {
		if semver.IsValid(v) {
			return v, nil
		}
	}
	return "", ErrNoValidVersionFound
}
