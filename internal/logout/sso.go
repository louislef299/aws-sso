package logout

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	laws "github.com/louislef299/aws-sso/pkg/v1/aws"
	los "github.com/louislef299/aws-sso/pkg/v1/os"
	"gopkg.in/ini.v1"
)

// Removes all sections that have the AWS_LOGIN_SUFFIX in the name
func Logout(ctx context.Context, cfg *aws.Config, cleanToken bool) error {
	var (
		clientinfo string
		info       *laws.ClientInformation
		err        error
	)

	if cleanToken {
		clientinfo, err = laws.ClientInfoFileDestination()
		if err != nil {
			return err
		}
		info, err = laws.ReadClientInformation(clientinfo)
		if err != nil {
			return fmt.Errorf("couldn't gather client login information: %v", err)
		}
	}

	if err := cleanConfig(); err != nil {
		return err
	}
	if err := cleanCredentials(); err != nil {
		return err
	}

	if cleanToken {
		if err := laws.Logout(ctx, cfg, info.AccessToken); err != nil {
			return fmt.Errorf("failed to logout of account: %v", err)
		}

		return os.Remove(clientinfo)
	}
	return nil
}

// Cleans the provided file with AWS_LOGIN_SUFFIX
func clean(file string) error {
	exists, err := los.IsFileOrFolderExisting(file)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("%s does not exist", file)
	}
	return deleteSections(file, los.AWS_LOGIN_SUFFIX)
}

// Cleans the config file
func cleanConfig() error {
	config, err := getConfigFile()
	if err != nil {
		return fmt.Errorf("issue getting config file: %v", err)
	}
	return clean(config)
}

// Cleans the credentials file
func cleanCredentials() error {
	creds, err := getCredentialsFile()
	if err != nil {
		return fmt.Errorf("issue getting credentials file: %v", err)
	}
	return clean(creds)
}

// Returns the path to the aws config file
func getConfigFile() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return homeDir + "/.aws/config", nil
}

// Returns the path to the aws credentials file
func getCredentialsFile() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return homeDir + "/.aws/credentials", nil
}

// Loads the toml file provided
func loadConfig(file string) (*ini.File, error) {
	cfg, err := ini.Load(file)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// Deletes all toml sections given a prefix
func deleteSections(file string, prefix string) error {
	cfg, err := loadConfig(file)
	if err != nil {
		return err
	}
	sections := cfg.SectionStrings()
	for _, s := range sections {
		if strings.Contains(s, prefix) {
			cfg.DeleteSection(s)
		}
	}
	return cfg.SaveTo(file)
}
