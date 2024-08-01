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
	"os"

	"github.com/suyuan32/goctls/rpc/execx"
	"github.com/suyuan32/goctls/util/ctx"
)

const (
	goZeroMod = "github.com/zeromicro/go-zero"
	adminTool = "github.com/suyuan32/simple-admin-tools"
)

var errInvalidGoMod = errors.New("it's only working for go module")

func editMod(zeroVersion, toolVersion string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	isGoMod, _ := ctx.IsGoMod(wd)
	if !isGoMod {
		return nil
	}

	mod := fmt.Sprintf("%s@%s", goZeroMod, zeroVersion)

	err = addRequire(mod)
	if err != nil {
		return err
	}

	// add replace
	mod = fmt.Sprintf("%s@%s=%s@%s", goZeroMod, zeroVersion, adminTool, toolVersion)

	err = addReplace(mod)
	if err != nil {
		return err
	}

	return nil
}

func addRequire(mod string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	isGoMod, _ := ctx.IsGoMod(wd)
	if !isGoMod {
		return errInvalidGoMod
	}

	_, err = execx.Run("go mod edit -require "+mod, wd)
	return err
}

func addReplace(mod string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	isGoMod, _ := ctx.IsGoMod(wd)
	if !isGoMod {
		return errInvalidGoMod
	}

	_, err = execx.Run("go mod edit -replace "+mod, wd)
	return err
}

func tidy() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	isGoMod, _ := ctx.IsGoMod(wd)
	if !isGoMod {
		return nil
	}

	_, err = execx.Run("go mod tidy", wd)
	if err != nil {
		return err
	}

	_, err = execx.Run("go mod edit -fmt", wd)
	if err != nil {
		return err
	}

	return err
}
