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

package vben5

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"

	"github.com/suyuan32/goctls/util"
)

func genApi(g *GenContext) error {
	if err := util.With("apiTpl").Parse(apiTpl).SaveTo(map[string]any{
		"modelName":           g.ModelName,
		"modelNameSpace":      strings.Replace(strcase.ToSnake(g.ModelName), "_", " ", -1),
		"modelNameLowerCamel": strcase.ToLowerCamel(g.ModelName),
		"modelNameSnake":      strcase.ToSnake(g.ModelName),
		"prefix":              g.Prefix,
		"useUUID":             g.UseUUID,
		"hasStatus":           g.HasStatus,
	},
		filepath.Join(g.ApiDir, fmt.Sprintf("%s.ts", strcase.ToLowerCamel(g.ModelName))), g.Overwrite); err != nil {
		return err
	}
	return nil
}
