package schema

import (
	_ "embed"
	"errors"
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/suyuan32/goctls/util/format"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	VarStringModelName string
)

//go:embed tpl/basic.tpl
var schemaTpl string

func GenSchema(_ *cobra.Command, _ []string) error {
	if VarStringModelName == "" {
		return errors.New("the model name can not be empty")
	}

	var schemaStr strings.Builder

	schemaTpl, err := template.New("schemaTpl").Parse(schemaTpl)
	if err != nil {
		return err
	}

	err = schemaTpl.Execute(&schemaStr, map[string]string{
		"ModelName": VarStringModelName,
	})

	if err != nil {
		return err
	}

	var filePath string
	tmp, err := filepath.Abs(".")
	if err != nil {
		return err
	}

	if strings.HasSuffix(tmp, "schema") {
		filePath = tmp
	} else {
		newPath := filepath.Join(tmp, "ent/schema")
		if fileutil.IsExist(newPath) {
			filePath = newPath
		} else {
			return errors.New("failed to find the ent schema folder")
		}
	}

	filename, err := format.FileNamingFormat("go_zero", VarStringModelName)
	if err != nil {
		return err
	}

	err = fileutil.WriteStringToFile(filepath.Join(filePath, filename+".go"),
		schemaStr.String(), false)
	if err != nil {
		return err
	}

	color.Green.Println("Generate Ent schema successfully")

	return nil
}
