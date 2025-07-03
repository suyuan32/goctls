package vben5ele

import (
	"path/filepath"

	"github.com/iancoleman/strcase"

	"github.com/suyuan32/goctls/util"
)

func genForm(g *GenContext) error {
	if g.FormType != "modal" {
		return nil
	}

	if err := util.With("formTpl").Parse(formTpl).SaveTo(map[string]any{
		"modelName":           g.ModelName,
		"modelNameLowerCamel": strcase.ToLowerCamel(g.ModelName),
		"folderName":          g.FolderName,
	},
		filepath.Join(g.ViewDir, "form.vue"), g.Overwrite); err != nil {
		return err
	}
	return nil
}
