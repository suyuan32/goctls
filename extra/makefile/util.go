// Copyright 2023 The Ryan SU Authors (https://github.com/suyuan32). All Rights Reserved.
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

package makefile

import (
	"strings"

	"github.com/duke-git/lancet/v2/fileutil"

	"github.com/pkg/errors"
)

// extractInfo extracts the information to context
func extractInfo(g *GenContext) error {
	makefileData, err := fileutil.ReadFileToString(g.TargetPath)
	if err != nil {
		return err
	}

	if strings.Contains(makefileData, "gen-api") {
		if strings.Contains(makefileData, "gen-ent") {
			g.UseEnt = true
			g.IsSingle = true
		} else {
			g.IsApi = true
		}
	}

	if strings.Contains(makefileData, "gen-rpc") {
		g.IsRpc = true
		if strings.Contains(makefileData, "gen-ent") {
			g.UseEnt = true
		}
	}

	dataSplit := strings.Split(makefileData, "\n")

	if g.Style == "" && !strings.Contains(makefileData, "PROJECT_STYLE=") {
		return errors.New("style not set, use -s to set style")
	} else if g.Style == "" && strings.Contains(makefileData, "PROJECT_STYLE=") {
		style := findDefined("PROJECT_STYLE=", dataSplit)
		if style == "" {
			return errors.New("failed to find style definition, please set it manually by -s")
		}
		g.Style = style
	}

	if val := findDefined("PROJECT_I18N", dataSplit); val != "" {
		if val == "true" {
			g.UseI18n = true
		}
	}

	if g.ServiceName == "" {
		g.ServiceName = findDefined("SERVICE", dataSplit)
	}

	g.EntFeature = findDefined("ENT_FEATURE", dataSplit)

	if g.EntFeature == "" {
		g.EntFeature = "sql/execquery,intercept"
	}

	return err
}

func findDefined(target string, data []string) string {
	for _, v := range data {
		if strings.Contains(v, target) {
			dataSplit := strings.Split(v, "=")
			if len(dataSplit) == 2 {
				return strings.TrimSpace(dataSplit[1])
			} else {
				return ""
			}
		}
	}

	return ""
}
