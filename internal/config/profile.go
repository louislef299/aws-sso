package config

import (
	"errors"
	"fmt"
	"io"

	"github.com/louislef299/knot/internal/envs"
	"github.com/louislef299/knot/pkg/provider"
	"github.com/spf13/viper"
)

type Profile struct {
	Name     string                `toml:"name"`
	Provider string                `toml:"provider"`
	Config   provider.ConfigSchema `toml:"config"`
}

var ErrNoProfilesFound = errors.New("there were no profiles found in config")

func ListProfiles(out io.Writer) error {
	v := viper.Sub(envs.PROFILE_HEADER)
	if v == nil {
		return ErrNoProfilesFound
	}

	for _, k := range v.AllKeys() {
		fmt.Fprintf(out, "%s\n", k)
	}
	return nil
}
