package vben

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/suyuan32/goctls/util"
	"path/filepath"
)

func genModalIndex(g *GenContext) error {
	if g.FormType != "modal" {
		return nil
	}

	if err := util.With("modelIndexTpl").Parse(modalIndexTpl).SaveTo(map[string]any{
		"modelName":           g.ModelName,
		"modelNameLowerCamel": strcase.ToLowerCamel(g.ModelName),
		"folderName":          g.FolderName,
		"addButtonTitle":      fmt.Sprintf("{{ t('%s.%s.add%s') }}", g.FolderName, strcase.ToLowerCamel(g.ModelName), g.ModelName),
		"deleteButtonTitle":   "{{ t('common.delete') }}",
		"useUUID":             g.UseUUID,
	},
		filepath.Join(g.ViewDir, "index.vue"), g.Overwrite); err != nil {
		return err
	}
	return nil
}
