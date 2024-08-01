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

package vben

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"

	"github.com/suyuan32/goctls/api/spec"
	"github.com/suyuan32/goctls/util"
)

func genModel(g *GenContext) error {
	var infoData strings.Builder
	for _, v := range g.ApiSpec.Types {
		if v.Name() == fmt.Sprintf("%sInfo", g.ModelName) {
			specData, ok := v.(spec.DefineStruct)
			if !ok {
				return errors.New("cannot get the field")
			}

			for _, val := range specData.Members {
				if val.Name == "" {
					tmpType, _ := val.Type.(spec.DefineStruct)
					if tmpType.Name() == "BaseIDInfo" {
						infoData.WriteString("  id?: number;\n  createdAt?: number;\n  updatedAt?: number;\n")
					} else if tmpType.Name() == "BaseUUIDInfo" {
						infoData.WriteString("  id?: string;\n  createdAt?: number;\n  updatedAt?: number;\n")
						g.UseUUID = true
					}
				} else {
					if val.Name == "Status" {
						g.HasStatus = true
					}

					if val.Name == "State" {
						g.HasState = true
					}

					infoData.WriteString(fmt.Sprintf("  %s?: %s;\n", strcase.ToLowerCamel(val.Name),
						ConvertGoTypeToTsType(val.Type.Name())))
				}
			}

		}
	}

	if infoData.Len() < 5 {
		return errors.New("failed to get the fields of the model, please check the api file and your model name")
	}

	if err := util.With("modelTpl").Parse(modelTpl).SaveTo(map[string]any{
		"modelName": g.ModelName,
		"infoData":  infoData.String(),
	},
		filepath.Join(g.ModelDir, fmt.Sprintf("%sModel.ts", strcase.ToLowerCamel(g.ModelName))), g.Overwrite); err != nil {
		return err
	}
	return nil
}
