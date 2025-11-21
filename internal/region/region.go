package region

import (
	"errors"
	"os"
	"strings"

	"github.com/louislef299/knot/internal/envs"
	"github.com/spf13/viper"
)

const (
	STS = iota
	EKS
	ECR
)

var (
	AwsRegions = []string{
		"us-east-2",
		"us-east-1",
		"us-west-1",
		"us-west-2",
		"af-south-1",
		"ap-east-1",
		"ap-south-1",
		"ap-northeast-3",
		"ap-northeast-2",
		"ap-southeast-1",
		"ap-southeast-2",
		"ap-northeast-1",
		"ca-central-1",
		"eu-central-1",
		"eu-west-1",
		"eu-west-2",
		"eu-south-1",
		"eu-west-3",
		"eu-north-1",
		"me-south-1",
		"sa-east-1",
		"us-gov-east-1",
		"us-gov-west-1",
		"cn-north-1",
		"cn-northwest-1",
	}

	ErrEmptyResponse  = errors.New("an empty response was returned")
	ErrRegionInvalid  = errors.New("the provided region is invalid")
	ErrRegionNotFound = errors.New("could not find a region in the system")
)

func checkOutput(r string) error {
	if r != "" && validRegion(r) {
		return nil
	}
	if r == "" {
		return ErrEmptyResponse
	}
	if !validRegion(r) {
		return ErrRegionInvalid
	}
	return ErrRegionNotFound
}

func GetRegion(regionType int) (string, error) {
	switch regionType {
	case STS:
		return GetSTSRegion()
	case EKS, ECR:
		return GetResourceRegion()
	default:
		return "", ErrRegionInvalid
	}
}

// Returns the region in precedence of environment
// region, config region and finally default region.
func GetResourceRegion() (string, error) {
	r := viper.GetString(envs.SESSION_REGION)
	if err := checkOutput(r); err == nil {
		return r, nil
	} else if err == ErrRegionInvalid {
		return "", err
	}

	r = viper.GetString(envs.CORE_DEFAULT_REGION)
	if err := checkOutput(r); err == nil {
		return r, nil
	} else if err == ErrRegionInvalid {
		return "", err
	}

	r = os.Getenv("AWS_REGION")
	if err := checkOutput(r); err == nil {
		return r, nil
	} else if err == ErrRegionInvalid {
		return "", err
	}

	r = os.Getenv("AWS_DEFAULT_REGION")
	if err := checkOutput(r); err == nil {
		return r, nil
	} else if err == ErrRegionInvalid {
		return "", err
	}

	return "", ErrRegionNotFound
}

// Returns the region in precedence of environment
// region, config region and finally default region.
func GetSTSRegion() (string, error) {
	r := viper.GetString(envs.CORE_SSO_REGION)
	if err := checkOutput(r); err == nil {
		return r, nil
	} else if err == ErrRegionInvalid {
		return "", err
	}

	r = viper.GetString(envs.CORE_DEFAULT_REGION)
	if err := checkOutput(r); err == nil {
		return r, nil
	} else if err == ErrRegionInvalid {
		return "", err
	}

	r = os.Getenv("AWS_REGION")
	if err := checkOutput(r); err == nil {
		return r, nil
	} else if err == ErrRegionInvalid {
		return "", err
	}

	r = os.Getenv("AWS_DEFAULT_REGION")
	if err := checkOutput(r); err == nil {
		return r, nil
	} else if err == ErrRegionInvalid {
		return "", err
	}

	r = viper.GetString(envs.SESSION_REGION)
	if err := checkOutput(r); err == nil {
		return r, nil
	} else if err == ErrRegionInvalid {
		return "", err
	}

	return "", ErrRegionNotFound
}

func validRegion(region string) bool {
	for _, r := range AwsRegions {
		if strings.Compare(r, strings.TrimSpace(region)) == 0 {
			return true
		}
	}
	return false
}
