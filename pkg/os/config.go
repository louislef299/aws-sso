package os

import (
	"fmt"
	"os"
	"regexp"
)

const (
	AWS_LOGIN_PATH   = "/.kube/aws-sso"
	AWS_LOGIN_SUFFIX = "aws-sso"

	configProfileRegex = `^*-aws-sso$`
)

func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return home + AWS_LOGIN_PATH, nil
}

func AddProfileSuffix(profile string) string {
	return fmt.Sprintf("%s-%s", profile, AWS_LOGIN_SUFFIX)
}

func IsManagedProfile(profile string) bool {
	r, err := regexp.Compile(configProfileRegex)
	if err != nil {
		panic(err)
	}
	return r.MatchString(profile)
}
