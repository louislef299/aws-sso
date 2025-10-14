package aws

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/louislef299/aws-sso/internal/envs"
	los "github.com/louislef299/aws-sso/pkg/os"
	"github.com/spf13/viper"
)

const (
	ACCESS_TOKEN_FILE    = "-token.json"
	AWS_TOKEN_PATH       = ".aws/sso/cache/"
	DEFAULT_ACCESS_TOKEN = "access-token.json"
)

type ClientInformation struct {
	AccessTokenExpiresAt    time.Time
	AccessToken             string
	ClientId                string
	ClientSecret            string
	ClientSecretExpiresAt   string
	DeviceCode              string
	VerificationUriComplete string
	StartUrl                string
}

type CredentialsFileTemplate struct {
	AwsAccessKeyId     string `ini:"aws_access_key_id,omitempty"`
	AwsSecretAccessKey string `ini:"aws_secret_access_key,omitempty"`
	AwsSessionToken    string `ini:"aws_session_token,omitempty"`
	CredentialProcess  string `ini:"credential_process,omitempty"`
	Output             string `ini:"output,omitempty"`
	Region             string `ini:"region,omitempty"`
}

// Will attempt to read in client information given a file location
func ReadClientInformation(file string) (*ClientInformation, error) {
	if exists, err := los.IsFileOrFolderExisting(file); err != nil {
		return nil, err
	} else if !exists {
		return nil, os.ErrNotExist
	}

	c := ClientInformation{}
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(content, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

// Checks to see if the ClientInformation AccessTokenExpiresAt is past the current time
func (c *ClientInformation) IsExpired() bool {
	return c.AccessTokenExpiresAt.Before(time.Now())
}

func IsAccessTokenExpired() (bool, error) {
	infoDest, err := ClientInfoFileDestination()
	if err != nil {
		return false, err
	}

	clientInfo, err := ReadClientInformation(infoDest)
	if err != nil {
		return false, err
	}

	return clientInfo.IsExpired(), nil
}

func GetAccessToken() string {
	t := viper.GetString(envs.SESSION_TOKEN)
	if t == DEFAULT_ACCESS_TOKEN || t == "" {
		return DEFAULT_ACCESS_TOKEN
	}
	return fmt.Sprintf("%s%s", t, ACCESS_TOKEN_FILE)
}
