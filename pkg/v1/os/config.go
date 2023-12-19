package os

import (
	"fmt"
	"os"
)

const (
	AWS_LOGIN_PATH   = "/.kube/aws-sso"
	AWS_LOGIN_PREFIX = "aws-sso"
)

func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return home + AWS_LOGIN_PATH, nil
}

func GetProfile(profile string) string {
	return fmt.Sprintf("%s-%s", profile, AWS_LOGIN_PREFIX)
}
