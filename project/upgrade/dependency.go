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

package upgrade

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/duke-git/lancet/v2/fileutil"

	"github.com/suyuan32/goctls/config"
	"github.com/suyuan32/goctls/rpc/execx"
)

func upgradeDependencies(workDir string) error {
	// drop old replace
	for _, v := range config.OldGoZeroVersion {
		_, err := execx.Run(fmt.Sprintf("go mod edit -dropreplace github.com/zeromicro/go-zero@%s", v), workDir)
		if err != nil {
			return errors.New("failed to drop old replace")
		}
	}

	data, err := fileutil.ReadFileToString(filepath.Join(workDir, "go.mod"))
	if err != nil {
		return err
	}

	err = upgradeOfficialDependencies(data, workDir)
	if err != nil {
		return err
	}

	err = tidy()
	if err != nil {
		return err
	}

	return nil
}

func upgradeOfficialDependencies(data, workDir string) (err error) {
	deps := []struct {
		Repo string
	}{
		{
			Repo: "github.com/suyuan32/simple-admin-common",
		},
		{
			Repo: "github.com/suyuan32/simple-admin-core",
		},
	}

	for _, v := range deps {
		if strings.Contains(data, v.Repo) {
			if strings.Contains(data, "simple-admin-core") {
				_, err = execx.Run(fmt.Sprintf("go mod edit -require=%s@%s", v.Repo,
					config.CoreVersion), workDir)
				if err != nil {
					return err
				}
			} else if strings.Contains(data, "simple-admin-common") {
				_, err = execx.Run(fmt.Sprintf("go mod edit -require=%s@%s", v.Repo,
					config.CommonVersion), workDir)
				if err != nil {
					return err
				}
			}
		}
	}

	return err
}
