package gogen

import (
	_ "embed"

	"github.com/suyuan32/goctls/config"
	"github.com/suyuan32/goctls/util/format"
)

//go:embed authortymiddleware.tpl
var authorityMiddlewareTemplate string

func genCasbin(dir string, cfg *config.Config, g *GenContext) error {
	fileName, err := format.FileNamingFormat(cfg.NamingFormat, "authority_middleware.go")
	if err != nil {
		return err
	}

	err = genFile(fileGenConfig{
		dir:             dir,
		subdir:          middlewareDir,
		filename:        fileName,
		templateName:    "authorityTemplate",
		category:        category,
		templateFile:    authorityTemplateFile,
		builtinTemplate: authorityMiddlewareTemplate,
		data: map[string]any{
			"useTrans": g.TransErr,
		},
	})
	if err != nil {
		return err
	}

	return nil
}
