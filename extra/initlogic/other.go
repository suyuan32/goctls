// Copyright 2023 The Ryan SU Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package initlogic

import (
	_ "embed"
	"fmt"
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/pkg/errors"
	"github.com/suyuan32/goctls/util/format"
	"path/filepath"
	"strings"
	"text/template"

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
