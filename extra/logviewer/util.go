package logviewer

import (
	"errors"
	"os/user"
	"path/filepath"
)

func getWorkspaceConfigDir() (string, error) {
	userDir, err := user.Current()
	if err != nil {
		return "", errors.Join(err, errors.New("failed to get the user directory"))
	}

	configFile := filepath.Join(userDir.HomeDir, ".goctls/log_workspace_config.txt")

	return configFile, nil
}
