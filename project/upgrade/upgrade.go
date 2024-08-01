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
	"os"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	conf "github.com/suyuan32/goctls/config"
	"github.com/suyuan32/goctls/rpc/execx"
)

var (
	// VarBoolUpgradeMakefile describe whether to upgrade makefile
	VarBoolUpgradeMakefile bool
)

func UpgradeProject(_ *cobra.Command, _ []string) error {
	color.Green.Println("Start upgrading dependencies...")

	err := editMod(conf.DefaultGoZeroVersion, conf.DefaultToolVersion)
	if err != nil {
		return err
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	err = upgradeDependencies(wd)
	if err != nil {
		return err
	}

	if VarBoolUpgradeMakefile {
		color.Green.Println("Start upgrading Makefile ...")
		_, err = execx.Run("goctls extra makefile", wd)
		if err != nil {
			return errors.New("failed to upgrade makefile")
		}
	}

	color.Green.Println("Done.")

	return nil
}
