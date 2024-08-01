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
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"

	"github.com/suyuan32/goctls/api/spec"
	"github.com/suyuan32/goctls/util"
	"github.com/suyuan32/goctls/util/pathx"
)

func genLocale(g *GenContext) error {
	var localeEnData, localeZhData strings.Builder
	var enLocaleFileName, zhLocaleFileName string
	enLocaleFileName = filepath.Join(g.LocaleDir, "en", fmt.Sprintf("%s.ts", strings.ToLower(g.FolderName)))
	zhLocaleFileName = filepath.Join(g.LocaleDir, "zh-CN", fmt.Sprintf("%s.ts", strings.ToLower(g.FolderName)))

	for _, v := range g.ApiSpec.Types {
		if v.Name() == fmt.Sprintf("%sInfo", g.ModelName) {
			specData, ok := v.(spec.DefineStruct)
			if !ok {
				return errors.New("cannot get the field")
			}

			localeEnData.WriteString(fmt.Sprintf("  %s: {\n", strcase.ToLowerCamel(g.ModelName)))
			localeZhData.WriteString(fmt.Sprintf("  %s: {\n", strcase.ToLowerCamel(g.ModelName)))

			for _, val := range specData.Members {
				if val.Name != "" {
					localeEnData.WriteString(fmt.Sprintf("    %s: '%s',\n",
						strcase.ToLowerCamel(val.Name), strcase.ToCamel(val.Name)))

					localeZhData.WriteString(fmt.Sprintf("    %s: '%s',\n",
						strcase.ToLowerCamel(val.Name), strcase.ToCamel(val.Name)))
				}
			}

			localeEnData.WriteString(fmt.Sprintf("    add%s: 'Add %s',\n", g.ModelName, g.ModelName))
			localeEnData.WriteString(fmt.Sprintf("    edit%s: 'Edit %s',\n", g.ModelName, g.ModelName))
			localeEnData.WriteString(fmt.Sprintf("    %sList: '%s List',\n", strcase.ToLowerCamel(g.ModelName), g.ModelName))
			localeEnData.WriteString("  },\n")

			localeZhData.WriteString(fmt.Sprintf("    add%s: '添加 %s',\n", g.ModelName, g.ModelName))
			localeZhData.WriteString(fmt.Sprintf("    edit%s: '编辑 %s',\n", g.ModelName, g.ModelName))
			localeZhData.WriteString(fmt.Sprintf("    %sList: '%s 列表',\n", strcase.ToLowerCamel(g.ModelName), g.ModelName))
			localeZhData.WriteString("  },\n")
		}
	}

	if !pathx.FileExists(enLocaleFileName) {
		if err := util.With("localeTpl").Parse(localeTpl).SaveTo(map[string]any{
			"localeData": localeEnData.String(),
		},
			enLocaleFileName, false); err != nil {
			return err
		}
	} else {
		file, err := os.ReadFile(enLocaleFileName)
		if err != nil {
			return err
		}

		data := string(file)

		if !strings.Contains(data, strings.ToLower(g.ModelName)+":") {
			data = data[:len(data)-3] + localeEnData.String() + data[len(data)-3:]
		} else if g.Overwrite {
			begin, end := FindBeginEndOfLocaleField(data, strings.ToLower(g.ModelName))
			data = data[:begin-2] + localeEnData.String() + data[end+1:]
		}

		err = os.WriteFile(enLocaleFileName, []byte(data), os.ModePerm)
		if err != nil {
			return err
		}
	}

	if !pathx.FileExists(zhLocaleFileName) {
		if err := util.With("localeTpl").Parse(localeTpl).SaveTo(map[string]any{
			"localeData": localeZhData.String(),
		},
			zhLocaleFileName, false); err != nil {
			return err
		}
	} else {
		file, err := os.ReadFile(zhLocaleFileName)
		if err != nil {
			return err
		}

		data := string(file)

		if !strings.Contains(data, strings.ToLower(g.ModelName)+":") {
			data = data[:len(data)-3] + localeZhData.String() + data[len(data)-3:]
		} else if g.Overwrite {
			begin, end := FindBeginEndOfLocaleField(data, strings.ToLower(g.ModelName))
			data = data[:begin-2] + localeZhData.String() + data[end+1:]
		}

		err = os.WriteFile(zhLocaleFileName, []byte(data), os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}
