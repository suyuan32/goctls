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

package gitea

import (
	_ "embed"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/gookit/color"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

//go:embed gitea.tpl
var giteaTpl string

var (
	VarStringRepository string
	VarStringOutputDir  string
	VarBoolChina        bool
)

func Gen(_ *cobra.Command, _ []string) (err error) {
	abs, err := filepath.Abs(VarStringOutputDir)
	if err != nil {
		return errors.Wrap(err, "dir not found")
	}

	// validate
	if VarStringRepository == "" || !strings.HasSuffix(VarStringRepository, ".git") {
		return errors.New("wrong repository, please set repository by \"-r\", such as \"-r https://github.com/suyuan32/simple-admin-job.git\"  ")
	}

	color.Green.Println("Generating...")

	if !fileutil.IsExist(filepath.Join(abs, ".gitea/workflows") + "/") {
		err := fileutil.CreateDir(filepath.Join(abs, ".gitea/workflows") + "/")
		if err != nil {
			return errors.Wrap(err, "failed to create the directory")
		}
	}

	tpl, err := template.New("gitea").Parse(giteaTpl)
	if err != nil {
		return errors.Wrap(err, "failed to load gitea template")
	}

	var strOutput strings.Builder

	err = tpl.Execute(&strOutput, map[string]any{
		"china":      VarBoolChina,
		"dir":        extractDir(VarStringRepository),
		"repository": VarStringRepository,
	})
	if err != nil {
		return err
	}

	err = fileutil.WriteStringToFile(filepath.Join(abs, ".gitea/workflows/docker.yml"), strOutput.String(), false)
	if err != nil {
		return errors.Wrap(err, "write file failed")
	}

	return err
}

func extractDir(data string) string {
	var begin, end int
	for i := len(data) - 1; i >= 0; i-- {
		if data[i] == '.' {
			end = i
		}

		if data[i] == '/' {
			begin = i
			break
		}
	}
	return data[begin+1 : end]
}
