package os

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

// Checks either or not a target file is existing.
// Returns true if the target exists, otherwise false.
func IsFileOrFolderExisting(target string) (bool, error) {
	if _, err := os.Stat(target); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, fmt.Errorf("Could not determine if file or folder %q exists or not. Exiting.", target)
	}
}

func WriteStructToFile(payload interface{}, dest string) error {
	targetDir := path.Dir(dest)
	exists, err := IsFileOrFolderExisting(targetDir)
	if err != nil {
		return err
	}
	if !exists {
		err := os.MkdirAll(targetDir, 0700)
		if err != nil {
			return err
		}
	}
	file, err := json.MarshalIndent(payload, "", " ")
	if err != nil {
		return err
	}
	err = os.WriteFile(dest, file, 0600)
	if err != nil {
		return err
	}
	return nil
}
