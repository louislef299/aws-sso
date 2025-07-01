package docker

import (
	"fmt"
	"log"
	"strings"

	"github.com/docker/cli/cli/config"
	"github.com/docker/cli/cli/config/configfile"
	"github.com/docker/cli/cli/config/types"
	"github.com/docker/docker/registry"
)

func Login(username, password, serverAddress string) error {
	// strip excess headers
	serverAddress = registry.ConvertToHostname(serverAddress)

	// Get the auth config.
	dcfg, authConfig, err := configureAWSAuth(username, password, serverAddress)
	if err != nil {
		return err
	}

	// Save the config value.
	if err := dcfg.GetCredentialsStore(authConfig.ServerAddress).Store(authConfig); err != nil {
		return fmt.Errorf("saving credentials failed: %v", err)
	}

	log.Println("Docker login to server", serverAddress, "succeeded")
	return nil
}

// configureAWSAuth returns an types.AuthConfig from the specified user, password and server.
func configureAWSAuth(username, password, serverAddress string) (*configfile.ConfigFile, types.AuthConfig, error) {
	if username == "" || password == "" || serverAddress == "" {
		return nil, types.AuthConfig{}, fmt.Errorf("cannot have empty username, password or server address")
	}

	// load docker configuration values
	dcfg, err := config.Load(config.Dir())
	if err != nil {
		return dcfg, types.AuthConfig{}, fmt.Errorf("loading config file failed: %v", err)
	}
	authConfig, err := dcfg.GetAuthConfig(serverAddress)
	if err != nil {
		return dcfg, authConfig, fmt.Errorf("getting auth config for %s failed: %v", serverAddress, err)
	}

	// A credential helper is being used to populate authentication.
	if dcfg.CredentialHelpers[serverAddress] != "" {
		return dcfg, authConfig, nil
	}

	authConfig.Username = strings.TrimSpace(username)
	authConfig.Password = password
	authConfig.ServerAddress = serverAddress
	authConfig.IdentityToken = ""

	return dcfg, authConfig, nil
}
