package gogen

import (
	_ "embed"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/suyuan32/goctls/api/spec"
	"github.com/suyuan32/goctls/config"
	"github.com/suyuan32/goctls/util/format"
)

//go:embed middleware.tpl
var middlewareImplementCode string

func genMiddleware(dir string, cfg *config.Config, api *spec.ApiSpec) error {
	middlewares := getMiddleware(api)
	for _, item := range middlewares {
		middlewareFilename := strings.TrimSuffix(strings.ToLower(item), "middleware") + "_middleware"
		filename, err := format.FileNamingFormat(cfg.NamingFormat, middlewareFilename)
		if err != nil {
			return err
		}

		name := strings.TrimSuffix(item, "Middleware") + "Middleware"
		err = genFile(fileGenConfig{
			dir:             dir,
			subdir:          middlewareDir,
			filename:        filename + ".go",
			templateName:    "contextTemplate",
			category:        category,
			templateFile:    middlewareImplementCodeFile,
			builtinTemplate: middlewareImplementCode,
			data: map[string]any{
				"name": cases.Title(language.English, cases.NoLower).String(name),
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}
