package aws

import (
	"errors"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	los "github.com/louislef299/aws-sso/pkg/v1/os"
	"gopkg.in/ini.v1"
)

const profileRegex = `^profile .*$`

var ErrFileNotFound = errors.New("the provided file could not be found")

func checkAWSFiles(files []string) error {
	for _, f := range files {
		exists, err := los.IsFileOrFolderExisting(f)
		if err != nil {
			return err
		}

		if !exists {
			dir := path.Dir(f)

			exists, err = los.IsFileOrFolderExisting(dir)
			if err != nil {
				return err
			}

			if dir != "." && !exists {
				err = os.MkdirAll(dir, 0777)
				if err != nil {
					return err
				}
			}

			return os.WriteFile(f, []byte(""), 0644)
		}
	}
	return nil
}

func GetAWSConfigSections(filename string) ([]string, error) {
	exists, err := los.IsFileOrFolderExisting(filename)
	if err != nil {
		return nil, err
	} else if !exists {
		return nil, ErrFileNotFound
	}

	cfg, err := ini.Load(filename)
	if err != nil {
		return nil, err
	}

	profr, err := regexp.Compile(profileRegex)
	if err != nil {
		return nil, err
	}

	sections := cfg.SectionStrings()
	var validSections []string
	for _, s := range sections {
		if profr.MatchString(s) {
			validSections = append(validSections, GetAWSProfileName(s))
		}
	}
	return validSections, nil
}

func GetAWSProfiles() ([]string, error) {
	files := []string{
		config.DefaultSharedConfigFilename(),
		config.DefaultSharedCredentialsFilename(),
	}

	err := checkAWSFiles(files)
	if err != nil {
		return nil, err
	}

	var profiles []string
	for _, f := range files {
		p, err := GetAWSConfigSections(f)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, p...)
	}

	sort.Strings(profiles)
	return profiles, nil
}

func GetAWSProfileName(profile string) string {
	return strings.TrimSpace(strings.TrimLeft(profile, "profile"))
}
