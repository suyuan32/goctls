package gogen

import (
	_ "embed"
	"fmt"
	"path"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/suyuan32/goctls/api/spec"
	"github.com/suyuan32/goctls/config"
	"github.com/suyuan32/goctls/util"
	"github.com/suyuan32/goctls/util/format"
	"github.com/suyuan32/goctls/util/pathx"
)

const defaultLogicPackage = "logic"

//go:embed handler.tpl
var handlerTemplate string

func genHandlers(dir, rootPkg string, cfg *config.Config, api *spec.ApiSpec, g *GenContext) error {
	for _, group := range api.Service.Groups {
		for _, route := range group.Routes {
			if err := genHandler(dir, rootPkg, cfg, group, route, g); err != nil {
				return err
			}
		}
	}

	return nil
}

func genHandler(dir, rootPkg string, cfg *config.Config, group spec.Group, route spec.Route, g *GenContext) error {
	handler := getHandlerName(route)
	handlerPath := getHandlerFolderPath(group, route)
	pkgName := handlerPath[strings.LastIndex(handlerPath, "/")+1:]
	logicName := defaultLogicPackage
	if handlerPath != handlerDir {
		handler = cases.Title(language.English, cases.NoLower).String(handler)
		logicName = pkgName
	}

	var swaggerPath string
	if strings.Contains(route.Path, ":") {
		swaggerPath = ConvertRoutePathToSwagger(route.Path)
	} else {
		swaggerPath = route.Path
	}

	prefix := group.GetAnnotation(spec.RoutePrefixKey)
	prefix = strings.ReplaceAll(prefix, `"`, "")
	prefix = strings.TrimSpace(prefix)
	if len(prefix) > 0 {
		prefix = path.Join("/", prefix)
	}
	swaggerPath = path.Join("/", prefix, swaggerPath)

	// sse
	var useSSE bool
	if group.GetAnnotation("sse") == "true" {
		useSSE = true
	}

	filename, err := format.FileNamingFormat(cfg.NamingFormat, handler)
	if err != nil {
		return err
	}

	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          getHandlerFolderPath(group, route),
		filename:        filename + ".go",
		templateName:    "handlerTemplate",
		category:        category,
		templateFile:    handlerTemplateFile,
		builtinTemplate: handlerTemplate,
		data: map[string]any{
			"PkgName":        pkgName,
			"ImportPackages": genHandlerImports(group, route, rootPkg, useSSE),
			"HandlerName":    handler,
			"RequestType":    util.Title(route.RequestTypeName()),
			"ResponseType":   responseGoTypeName(route, typesPacket),
			"LogicName":      logicName,
			"LogicType":      cases.Title(language.English, cases.NoLower).String(getLogicName(route)),
			"Call":           cases.Title(language.English, cases.NoLower).String(strings.TrimSuffix(handler, "Handler")),
			"HasResp":        len(route.ResponseTypeName()) > 0,
			"HasRequest":     len(route.RequestTypeName()) > 0,
			"HasDoc":         len(route.JoinedDoc()) > 0,
			"Doc":            getDoc(route.JoinedDoc()),
			"TransErr":       g.TransErr,
			"UseValidator":   g.UseValidator,
			"UseSSE":         useSSE,
		},
	})
}

func genHandlerImports(group spec.Group, route spec.Route, parentPkg string, useSSE bool) string {
	imports := []string{
		fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, getLogicFolderPath(group, route))),
		fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, contextDir)),
	}

	if useSSE {
		imports = append(imports, "\"encoding/json\"", "\"github.com/zeromicro/go-zero/core/logc\"",
			"\"github.com/zeromicro/go-zero/core/threading\"")
	}

	if len(route.RequestTypeName()) > 0 {
		imports = append(imports, fmt.Sprintf("\"%s\"\n", pathx.JoinPackages(parentPkg, typesDir)))
	}

	return strings.Join(imports, "\n\t")
}

func getHandlerBaseName(route spec.Route) (string, error) {
	handler := route.Handler
	handler = strings.TrimSpace(handler)
	handler = strings.TrimSuffix(handler, "handler")
	handler = strings.TrimSuffix(handler, "Handler")

	return handler, nil
}

func getHandlerFolderPath(group spec.Group, route spec.Route) string {
	folder := route.GetAnnotation(groupProperty)
	if len(folder) == 0 {
		folder = group.GetAnnotation(groupProperty)
		if len(folder) == 0 {
			return handlerDir
		}
	}

	folder = strings.TrimPrefix(folder, "/")
	folder = strings.TrimSuffix(folder, "/")

	return path.Join(handlerDir, folder)
}

func getHandlerName(route spec.Route) string {
	handler, err := getHandlerBaseName(route)
	if err != nil {
		panic(err)
	}

	return handler + "Handler"
}

func getLogicName(route spec.Route) string {
	handler, err := getHandlerBaseName(route)
	if err != nil {
		panic(err)
	}

	return handler + "Logic"
}
