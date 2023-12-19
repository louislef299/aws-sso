package aws

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	los "github.com/louislef299/aws-sso/pkg/v1/os"
	"gopkg.in/ini.v1"
)

func WriteAWSCredentialsFile(profile string, credentials *sso.GetRoleCredentialsOutput) error {
	credentialsFile := config.DefaultSharedCredentialsFilename()
	exists, err := los.IsFileOrFolderExisting(credentialsFile)
	if err != nil {
		return err
	}
	if !exists {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		err = os.MkdirAll(homeDir+"/.aws", 0777)
		if err != nil {
			return err
		}
		_, err = os.OpenFile(credentialsFile, os.O_CREATE, 0644)
		if err != nil {
			return err
		}
	}
	cfg, err := ini.Load(credentialsFile)
	if err != nil {
		return err
	}

	cfg.Section(profile).Key("aws_access_key_id").SetValue(*credentials.RoleCredentials.AccessKeyId)
	cfg.Section(profile).Key("aws_secret_access_key").SetValue(*credentials.RoleCredentials.SecretAccessKey)
	cfg.Section(profile).Key("aws_session_token").SetValue(*credentials.RoleCredentials.SessionToken)
	return cfg.SaveTo(credentialsFile)
}

func WriteAWSConfigFile(profile, region, output string) error {
	configFile := config.DefaultSharedConfigFilename()
	exists, err := los.IsFileOrFolderExisting(configFile)
	if err != nil {
		return err
	}
	if !exists {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		err = os.MkdirAll(homeDir+"/.aws", 0777)
		if err != nil {
			return err
		}
		_, err = os.OpenFile(configFile, os.O_CREATE, 0644)
		if err != nil {
			return err
		}
	}
	cfg, err := ini.Load(configFile)
	if err != nil {
		return err
	}

	var prof string
	if profile != "default" {
		prof = fmt.Sprintf("%s %s", "profile", profile)
	}

	cfg.Section(prof).Key("region").SetValue(region)
	cfg.Section(prof).Key("output").SetValue(output)
	return cfg.SaveTo(configFile)
}
