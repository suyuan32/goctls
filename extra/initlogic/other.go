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

package initlogic

import (
	_ "embed"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/pkg/errors"
	"github.com/suyuan32/goctls/util/format"

	"github.com/iancoleman/strcase"

	"github.com/suyuan32/goctls/util/console"
)

//go:embed other.tpl
var otherTpl string

//go:embed init.tpl
var initTpl string

func OtherGen(g *CoreGenContext) error {
	var otherString strings.Builder
	otherTemplate, err := template.New("init_other").Parse(otherTpl)
	if err != nil {
		return errors.Wrap(err, "failed to create other init template")
	}

	err = otherTemplate.Execute(&otherString, map[string]any{
		"modelName":      g.ModelName,
		"modelNameSnake": strcase.ToSnake(g.ModelName),
		"modelNameUpper": strings.ToUpper(g.ModelName),
		"serviceName":    g.ServiceName,
		"routePrefix":    g.RoutePrefix,
	})
	if err != nil {
		return err
	}

	if g.Target == "console" {
		console.Info(otherString.String())
	} else {
		absPath, err := filepath.Abs(g.Output)
		if err != nil {
			return errors.Wrap(err, "failed to find the output file")
		}

		apiFileName, err := format.FileNamingFormat(g.Style, "init_api_data.go")
		if err != nil {
			return err
		}

		if g.Output == "." {
			if fileutil.IsExist(filepath.Join(absPath, "internal/logic/base")) {
				path := filepath.Join(absPath, "internal/logic/base/", apiFileName)
				if !fileutil.IsExist(path) {
					err := fileutil.WriteStringToFile(path, initTpl, false)
					if err != nil {
						return errors.Wrap(err, fmt.Sprintf("failed to create API initialization file to path: %s", path))
					}
				}

				originalFileStr, err := fileutil.ReadFileToString(path)
				if err != nil {
					return err
				}

				if strings.Contains(originalFileStr, fmt.Sprintf("/%s/", strcase.ToSnake(g.ModelName))) {
					return errors.New("the init code already exist, if you still want to generate, use \"-o console\" instead")
				}

				index := strings.Index(originalFileStr, "insertApiData")

				if index == -1 {
					return fmt.Errorf("failed to find \"insertApiData\" function in file: %s", path)
				}

				newFileStr := originalFileStr[:index+29] + otherString.String() + originalFileStr[index+29:]

				err = fileutil.WriteStringToFile(path, newFileStr, false)
				if err != nil {
					return err
				}
			}

		}

	}

	return err
}
