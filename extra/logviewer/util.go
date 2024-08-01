// Copyright (C) 2023  Ryan SU (https://github.com/suyuan32)

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
