package gogen

import (
	_ "embed"
	"fmt"
	"path"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/suyuan32/goctls/api/parser/g4/gen/api"
	"github.com/suyuan32/goctls/api/spec"
	"github.com/suyuan32/goctls/config"
	"github.com/suyuan32/goctls/util/format"
	"github.com/suyuan32/goctls/util/pathx"
	"github.com/suyuan32/goctls/vars"
)

//go:embed logic.tpl
var logicTemplate string

func genLogic(dir, rootPkg string, cfg *config.Config, api *spec.ApiSpec, ctx *GenContext) error {
	for _, g := range api.Service.Groups {
		for _, r := range g.Routes {
			err := genLogicByRoute(dir, rootPkg, cfg, g, r, ctx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func genLogicByRoute(dir, rootPkg string, cfg *config.Config, group spec.Group, route spec.Route, ctx *GenContext) error {
	logic := getLogicName(route)
	goFile, err := format.FileNamingFormat(cfg.NamingFormat, logic)
	if err != nil {
		return err
	}

	// sse
	var useSSE bool
	if group.GetAnnotation("sse") == "true" {
		useSSE = true
	}

	imports := genLogicImports(route, rootPkg, useSSE)
	var responseString string
	var returnString string
	var requestString string
	if len(route.ResponseTypeName()) > 0 {
		resp := responseGoTypeName(route, typesPacket)
		responseString = "(resp " + resp + ", err error)"
		returnString = "return"
	} else {
		responseString = "error"
		returnString = "return nil"
	}
	if len(route.RequestTypeName()) > 0 {
		requestString = "req *" + requestGoTypeName(route, typesPacket)
	}

	if useSSE {
		responseString = "error"
		returnString = "return nil"
		resp := responseGoTypeName(route, typesPacket)
		if len(requestString) == 0 {
			requestString = "client chan<- " + resp
		} else {
			requestString += ", client chan<- " + resp
		}
	}

	subDir := getLogicFolderPath(group, route)
	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          subDir,
		filename:        goFile + ".go",
		templateName:    "logicTemplate",
		category:        category,
		templateFile:    logicTemplateFile,
		builtinTemplate: logicTemplate,
		data: map[string]any{
			"pkgName":      subDir[strings.LastIndex(subDir, "/")+1:],
			"imports":      imports,
			"logic":        cases.Title(language.English, cases.NoLower).String(logic),
			"function":     cases.Title(language.English, cases.NoLower).String(strings.TrimSuffix(logic, "Logic")),
			"responseType": responseString,
			"returnString": returnString,
			"request":      requestString,
			"useSSE":       useSSE,
		},
	})
}

func getLogicFolderPath(group spec.Group, route spec.Route) string {
	folder := route.GetAnnotation(groupProperty)
	if len(folder) == 0 {
		folder = group.GetAnnotation(groupProperty)
		if len(folder) == 0 {
			return logicDir
		}
	}
	folder = strings.TrimPrefix(folder, "/")
	folder = strings.TrimSuffix(folder, "/")
	return path.Join(logicDir, folder)
}

func genLogicImports(route spec.Route, parentPkg string, useSSE bool) string {
	var imports []string
	imports = append(imports, `"context"`+"\n")
	if useSSE {
		imports = append(imports, `"net/http"`+"\n")
	}

	imports = append(imports, fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, contextDir)))
	if shallImportTypesPackage(route) && !useSSE {
		imports = append(imports, fmt.Sprintf("\"%s\"\n", pathx.JoinPackages(parentPkg, typesDir)))
	}
	imports = append(imports, fmt.Sprintf("\"%s/core/logx\"", vars.ProjectOpenSourceURL))

	return strings.Join(imports, "\n\t")
}

func onlyPrimitiveTypes(val string) bool {
	fields := strings.FieldsFunc(val, func(r rune) bool {
		return r == '[' || r == ']' || r == ' '
	})

	for _, field := range fields {
		if field == "map" {
			continue
		}
		// ignore array dimension number, like [5]int
		if _, err := strconv.Atoi(field); err == nil {
			continue
		}
		if !api.IsBasicType(field) {
			return false
		}
	}

	return true
}

func shallImportTypesPackage(route spec.Route) bool {
	if len(route.RequestTypeName()) > 0 {
		return true
	}

	respTypeName := route.ResponseTypeName()
	if len(respTypeName) == 0 {
		return false
	}

	if onlyPrimitiveTypes(respTypeName) {
		return false
	}

	return true
}
