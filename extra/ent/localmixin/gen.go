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

package localmixin

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"
	"text/template"

	"github.com/duke-git/lancet/v2/fileutil"

	"github.com/suyuan32/goctls/util/ctx"

	"github.com/gookit/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"

	"github.com/suyuan32/goctls/util/console"
	"github.com/suyuan32/goctls/util/env"
	"github.com/suyuan32/goctls/util/pathx"
)

var (
	// VarStringDir is the ent directory
	VarStringDir string

	// VarStringAdd is the template name for adding
	VarStringAdd string

	// VarBoolUpdate describe whether to update the template
	VarBoolUpdate bool

	// VarBoolList describe whether to list all supported templates
	VarBoolList bool

	tplInfo table.Writer
)

func GenLocalMixin(_ *cobra.Command, _ []string) error {
	var entDir string
	var err error

	if VarBoolList {
		ListAllTemplate()
		return nil
	}

	if VarStringDir == "" {
		entDir, err = GetEntDir()
		if err != nil {
			return err
		}
	} else {
		entDir, err = filepath.Abs(VarStringDir)
		if err != nil {
			return err
		}
	}

	tmplDir := filepath.Join(entDir, "schema/localmixin/")

	if !fileutil.IsExist(tmplDir) {
		err = fileutil.CreateDir(tmplDir + "/")
		if err != nil {
			return err
		}
	}

	projectCtx, err := ctx.Prepare(tmplDir)
	if err != nil {
		return err
	}

	if VarBoolUpdate {
		files, err := pathx.GetFilesPathFromDir(tmplDir, false)
		if err != nil {
			return err
		}

		for _, v := range files {
			fileName := filepath.Base(v)
			tpl := GetTmpl(fileName, projectCtx)
			if tpl == "" {
				return errors.New("failed to find the template")
			}

			if fileutil.IsExist(v) {
				err := fileutil.RemoveFile(v)
				if err != nil {
					return errors.Join(err, errors.New("failed to remove the original template"))
				}
			}

			err = fileutil.WriteStringToFile(v, tpl, false)
			if err != nil {
				return err
			}
		}

		//execx.Run("go get -u entgo.io/ent@latest", entDir)
	}

	if VarStringAdd != "" {
		tpl := GetTmpl(VarStringAdd, projectCtx)
		if tpl == "" {
			return errors.New("failed to find the template")
		}

		filePath := filepath.Join(tmplDir, VarStringAdd+".go")

		if pathx.Exists(filePath) {
			return errors.New("the template already exists")
		}

		err := fileutil.WriteStringToFile(filePath, tpl, false)
		if err != nil {
			return err
		}
	}

	console.Success("Generating successfully")

	return nil
}

func GetEntDir() (string, error) {
	entDir, _ := filepath.Abs("./ent")

	if pathx.Exists(entDir) {
		return entDir, nil
	}

	entDir, _ = filepath.Abs("./rpc/ent")

	if pathx.Exists(entDir) {
		return entDir, nil
	}

	entDir, _ = filepath.Abs("./api/ent")

	if pathx.Exists(entDir) {
		return entDir, nil
	}

	return "", errors.New("failed to find the ent directory")
}

func GetTmpl(name string, ctxData *ctx.ProjectContext) string {
	var tplData string
	switch name {
	case "soft_delete", "soft_delete.go":
		tplData = softDeleteTpl
	}

	if tplData != "" {
		mixinTplData := bytes.NewBufferString("")
		mixinTpl, _ := template.New("localmixin").Parse(tplData)
		_ = mixinTpl.Execute(mixinTplData, map[string]any{
			"PackagePath": ctxData.Path,
		})
		return mixinTplData.String()
	}

	return ""
}

func ListAllTemplate() {
	type Info struct {
		Name  string
		Intro string
	}

	var data []Info
	tplInfo = table.NewWriter()

	if env.IsChinaEnv() {
		color.Green.Println("支持的模板:\n")
		tplInfo.AppendHeader(table.Row{"Mixin 模板名称", "Mixin 模板介绍"})
		data = []Info{
			{
				"soft_delete",
				"Ent 的软删除 Mixin 模板",
			},
		}
	} else {
		color.Green.Println("The mixin templates supported:\n")
		tplInfo.AppendHeader(table.Row{"Mixin Name", "Mixin Introduction"})
		data = []Info{
			{
				"soft_delete",
				"Ent soft delete mixin template",
			},
		}
	}

	for _, v := range data {
		tplInfo.AppendRows([]table.Row{
			{
				v.Name,
				v.Intro,
			},
		})
	}

	fmt.Println(tplInfo.Render())

	if env.IsChinaEnv() {
		color.Green.Println("\n使用方法： goctls extra ent mixin -a soft_delete -d ./ent ")
	} else {
		color.Green.Println("\nUsage: goctls extra ent mixin -a soft_delete -d ./ent")
	}
}
