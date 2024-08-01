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
