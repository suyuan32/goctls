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
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/duke-git/lancet/v2/fileutil"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"

	"github.com/suyuan32/goctls/util/console"
)

//go:embed core.tpl
var coreTpl string

type CoreGenContext struct {
	Target      string
	ModelName   string
	Output      string
	Style       string
	ServiceName string
	RoutePrefix string
}

func GenCore(g *CoreGenContext) error {
	var coreString strings.Builder
	coreTemplate, err := template.New("init_core").Parse(coreTpl)
	if err != nil {
		return errors.Wrap(err, "failed to create core init template")
	}

	err = coreTemplate.Execute(&coreString, map[string]any{
		"modelName":      g.ModelName,
		"modelNameSnake": strcase.ToSnake(g.ModelName),
		"modelNameLower": strings.ToLower(g.ModelName),
		"modelNameUpper": strings.ToUpper(g.ModelName),
		"serviceName":    g.ServiceName,
		"routePrefix":    g.RoutePrefix,
	})
	if err != nil {
		return err
	}

	if g.Output != "" {
		absPath, err := filepath.Abs(g.Output)
		if err != nil {
			return errors.Wrap(err, "failed to find the output file")
		}

		if g.Output == "." {
			absPath = filepath.Join(absPath, "internal/logic/base/init_database_api_data.go")
			if !fileutil.IsExist(absPath) {
				return fmt.Errorf("failed to find the target file: %s", absPath)
			}
		}

		apiData, err := os.ReadFile(absPath)

		originalString := string(apiData)

		insertIndex := strings.Index(originalString, "err := l.svcCtx.DB.API.CreateBulk")

		if insertIndex == -1 {
			return errors.New("cannot find the insert place in the output file")
		} else {
			newString := originalString[:insertIndex] + coreString.String() + originalString[insertIndex:]

			err := os.WriteFile(absPath, []byte(newString), 0o666)
			if err != nil {
				return errors.Wrap(err, "failed to write data to output file")
			}
		}
	} else {
		console.Info(coreString.String())
	}

	return err
}
