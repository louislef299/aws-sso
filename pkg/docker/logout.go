package docker

import (
	"fmt"
	"log"

	"github.com/docker/cli/cli/config"
	"github.com/docker/docker/registry"
)

func Logout(registryname string) error {
	registryname = registry.ConvertToHostname(registryname)
	dcfg, err := config.Load(config.Dir())
	if err != nil {
		return fmt.Errorf("loading config file failed: %v", err)
	}

	// check if we're logged in based on the records in the config file
	// which means it couldn't have user/pass cause they may be in the creds store
	if _, loggedIn := dcfg.AuthConfigs[registryname]; loggedIn {
		if err := dcfg.GetCredentialsStore(registryname).Erase(registryname); err != nil {
			return fmt.Errorf("could not erase credentials: %v", err)
		}
		log.Println("erased", registryname)
	} else {
		log.Println("wasn't logged into", registryname)
	}
	return nil
}
