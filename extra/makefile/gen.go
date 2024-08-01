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
	"bytes"
	_ "embed"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/gookit/color"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/suyuan32/goctls/util/format"
)

//go:embed makefile.tpl
var makefileTpl string

var (
	VarStringServiceName string
	VarStringStyle       string
	VarStringDir         string
	VarStringServiceType string
	VarBoolI18n          bool
	VarBoolEnt           bool
)

type GenContext struct {
	ServiceName string
	Style       string
	IsSingle    bool
	IsApi       bool
	IsRpc       bool
	UseI18n     bool
	UseEnt      bool
	TargetPath  string
	EntFeature  string
}

func Gen(_ *cobra.Command, _ []string) (err error) {
	ctx := GenContext{}

	var absPath string

	if VarStringDir != "" {
		absPath, err = filepath.Abs(VarStringDir)
		if err != nil {
			return errors.Wrap(err, "dir not found")
		}
	} else {
		absPath, err = filepath.Abs(".")
		if err != nil {
			return errors.Wrap(err, "dir not found")
		}
	}

	filePath := filepath.Join(absPath, "Makefile")

	ctx.TargetPath = filePath
	ctx.Style = VarStringStyle
	ctx.ServiceName = VarStringServiceName
	ctx.UseEnt = VarBoolEnt
	ctx.UseI18n = VarBoolI18n

	switch VarStringServiceType {
	case "api":
		ctx.IsApi = true
	case "single":
		ctx.IsSingle = true
	case "rpc":
		ctx.IsRpc = true
	}

	ctx.EntFeature = "sql/execquery,intercept"

	if fileutil.IsExist(ctx.TargetPath) {
		err = extractInfo(&ctx)
		if err != nil {
			return errors.Wrap(err, "failed to extract makefile info")
		}
	}

	err = DoGen(&ctx)

	return err
}

func DoGen(g *GenContext) error {
	color.Green.Println("Generating...")

	serviceNameStyle, err := format.FileNamingFormat(g.Style, g.ServiceName)
	if err != nil {
		return errors.Wrap(err, "failed to format service name")
	}

	makefileData := bytes.NewBufferString("")
	makefileTmpl, _ := template.New("makefile").Parse(makefileTpl)
	_ = makefileTmpl.Execute(makefileData, map[string]any{
		"serviceName":      strcase.ToCamel(g.ServiceName),
		"useEnt":           g.UseEnt,
		"style":            g.Style,
		"useI18n":          g.UseI18n,
		"serviceNameStyle": serviceNameStyle,
		"serviceNameLower": strings.ToLower(g.ServiceName),
		"serviceNameSnake": strcase.ToSnake(g.ServiceName),
		"serviceNameDash":  strings.ReplaceAll(strcase.ToSnake(g.ServiceName), "_", "-"),
		"isApi":            g.IsApi,
		"isSingle":         g.IsSingle,
		"isRpc":            g.IsRpc,
		"entFeature":       g.EntFeature,
	})

	if fileutil.IsExist(g.TargetPath) {
		err = fileutil.RemoveFile(g.TargetPath)
		if err != nil {
			return err
		}
	}

	err = fileutil.WriteStringToFile(g.TargetPath, makefileData.String(), false)
	if err != nil {
		return err
	}

	return nil
}
