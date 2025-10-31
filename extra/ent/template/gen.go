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

package template

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/duke-git/lancet/v2/fileutil"

	"github.com/gookit/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"

	"github.com/suyuan32/goctls/rpc/execx"
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

func GenTemplate(_ *cobra.Command, _ []string) error {
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

	tmplDir := filepath.Join(entDir, "template")

	if VarBoolUpdate {
		files, err := pathx.GetFilesPathFromDir(tmplDir, false)
		if err != nil {
			return err
		}

		for _, v := range files {
			fileName := filepath.Base(v)
			tpl := GetTmpl(fileName)
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

		execx.Run("go get -u entgo.io/ent@latest", entDir)
	}

	if VarStringAdd != "" {
		tpl := GetTmpl(VarStringAdd)
		if tpl == "" {
			return errors.New("failed to find the template")
		}

		filePath := filepath.Join(tmplDir, VarStringAdd+".tmpl")

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

func GetTmpl(name string) string {
	switch name {
	case "pagination.tmpl", "pagination":
		return PaginationTmpl
	case "set_not_nil.tmpl", "set_not_nil":
		return NotNilTmpl
	case "set_or_clear.tmpl", "set_or_clear":
		return SetOrClearTmpl
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
		tplInfo.AppendHeader(table.Row{"模板名称", "模板介绍"})
		data = []Info{
			{
				"set_not_nil",
				"Ent 非 nil 模板，用于输入值不为 Nil 时更新",
			},
			{
				"pagination",
				"Ent 分页模板",
			},
			{
				"set_or_clear",
				"Ent 若输入值为 nil 则设置为 null 的模板,使用方法: SetOrClear",
			},
		}
	} else {
		color.Green.Println("The templates supported:\n")
		tplInfo.AppendHeader(table.Row{"Name", "Introduction"})
		data = []Info{
			{
				"set_not_nil",
				"The template for updating the value when it is not nil",
			},
			{
				"pagination",
				"The template for paginating the data",
			},
			{
				"set_or_clear",
				"The template for update the values to null when the input value is nil. Usage: SetOrClear",
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
		color.Green.Println("\n使用方法： goctls extra ent template -a not_empty_update -d ./ent ")
	} else {
		color.Green.Println("\nUsage: goctls extra ent template -a not_empty_update -d ./ent")
	}
}
