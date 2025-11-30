package config

import (
	"errors"
	"fmt"
	"io"

	"github.com/louislef299/knot/internal/envs"
	"github.com/spf13/viper"
)

var ErrNoProfilesFound = errors.New("there were no profiles found in config")

func ListProfiles(out io.Writer) error {
	v := viper.Sub(envs.PROFILE_HEADER)
	if v == nil {
		return ErrNoProfilesFound
	}

	for k, v := range v.AllSettings() {
		fmt.Fprintf(out, "%s: %s\n", k, v)
	}
	fmt.Fprint(out, "\n")
	return nil
}
