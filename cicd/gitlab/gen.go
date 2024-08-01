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

package gitlab

import (
	_ "embed"
	"path/filepath"

	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/gookit/color"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

//go:embed gitlab.tpl
var gitlabTpl string

var (
	VarStringOutputDir string
)

func Gen(_ *cobra.Command, _ []string) (err error) {
	abs, err := filepath.Abs(VarStringOutputDir)
	if err != nil {
		return errors.Wrap(err, "dir not found")
	}

	color.Green.Println("Generating...")

	err = fileutil.WriteStringToFile(filepath.Join(abs, ".gitlab-ci.yml"), gitlabTpl, false)
	if err != nil {
		return errors.Wrap(err, "write file failed")
	}

	return err
}
