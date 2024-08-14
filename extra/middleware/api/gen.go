package api

import (
	"errors"
	"fmt"
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/gookit/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/suyuan32/goctls/extra/middleware/api/tmpl"
	"github.com/suyuan32/goctls/util/console"
	"github.com/suyuan32/goctls/util/env"
	"github.com/suyuan32/goctls/util/format"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	VarStringName   string
	VarBoolList     bool
	VarStringOutput string
	VarBoolI18n     bool
	VarStringStyle  string
)

func Gen(_ *cobra.Command, _ []string) error {
	if VarBoolList {
		ListAllMiddleware()
		return nil
	}

	if VarStringName == "" {
		return errors.New("please input the middleware name")
	}

	return DoGen(VarStringName, VarStringOutput, VarStringStyle, VarBoolI18n)
}

func DoGen(name, output, style string, useI18n bool) error {
	tmlData := GetTmpl(name)
	if tmlData == "" {
		return errors.New("middleware not found")
	}

	path, err := filepath.Abs(output)
	if err != nil {
		return fmt.Errorf("output path error: %w", err)
	}

	if !strings.Contains(path, "middleware") {
		path = filepath.Join(path, "internal/middleware")
	}

	parseTmpl, err := template.New("middleware").Parse(tmlData)
	if err != nil {
		return err
	}

	var fileData strings.Builder

	err = parseTmpl.Execute(&fileData, map[string]interface{}{
		"useTrans": useI18n,
	})
	if err != nil {
		return err
	}

	fileName := strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(name), "_tenant", ""), "_", "") + "_middleware.go"

	fileName, err = format.FileNamingFormat(style, fileName)
	if err != nil {
		return err
	}

	err = fileutil.WriteStringToFile(filepath.Join(path, fileName), fileData.String(), false)
	if err != nil {
		return err
	}

	console.Success("Generating successfully")

	return nil
}

func GetTmpl(name string) string {
	switch name {
	case "authority":
		return tmpl.AuthorityTmpl
	case "authority_tenant":
		return tmpl.AuthorityTenantTmpl
	case "data_perm":
		return tmpl.DataPermTmpl
	}
	return ""
}

func ListAllMiddleware() {
	type Info struct {
		Name  string
		Intro string
	}

	var data []Info
	tplInfo := table.NewWriter()

	if env.IsChinaEnv() {
		color.Green.Println("支持的中间件:\n")
		tplInfo.AppendHeader(table.Row{"中间件名称", "中间件介绍"})
		data = []Info{
			{
				"authority",
				"Casbin authority 权限中间件",
			},
			{
				"authority_tenant",
				"Casbin autority 多租户权限中间件",
			},
			{
				"data_perm",
				"数据权限中间件",
			},
		}
	} else {
		color.Green.Println("The middleware supported:\n")
		tplInfo.AppendHeader(table.Row{"Name", "Introduction"})
		data = []Info{
			{
				"authority",
				"Casbin authority permission middleware",
			},
			{
				"authority_tenant",
				"Casbin autority multi-tenant permission middleware",
			},
			{
				"data_perm",
				"Data permission middleware",
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
		color.Green.Println("\n使用方法： goctls extra middleware api -a authority")
	} else {
		color.Green.Println("\nUsage: goctls extra middleware api -a authority")
	}
}
