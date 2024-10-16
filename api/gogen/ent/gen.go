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

package ent

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/suyuan32/goctls/rpc/execx"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/entc/load"
	"github.com/gookit/color"
	"github.com/iancoleman/strcase"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/suyuan32/goctls/rpc/parser"
	"github.com/suyuan32/goctls/util/ctx"
	"github.com/suyuan32/goctls/util/entx"
	"github.com/suyuan32/goctls/util/format"
	"github.com/suyuan32/goctls/util/pathx"
)

const regularPerm = 0o666

type ApiLogicData struct {
	LogicName string
	LogicCode string
}

type GenEntLogicContext struct {
	Schema           string
	Output           string
	ServiceName      string
	Style            string
	ModelName        string
	SearchKeyNum     int
	GroupName        string
	UseUUID          bool
	JSONStyle        string
	UseI18n          bool
	ImportPrefix     string
	GenApiData       bool
	Overwrite        bool
	RoutePrefix      string
	IdType           string
	HasCreated       bool
	ModelChineseName string
	ModelEnglishName string
}

func (g GenEntLogicContext) Validate() error {
	if g.Schema == "" {
		return errors.New("the schema dir cannot be empty ")
	} else if !strings.HasSuffix(g.Schema, "schema") {
		return errors.New("please input correct schema directory e.g. ./ent/schema ")
	} else if g.ServiceName == "" {
		return errors.New("please set the API service name via --api_service_name")
	}
	return nil
}

// GenEntLogic generates the ent CRUD logic files of the api service.
func GenEntLogic(g *GenEntLogicContext) error {
	return genEntLogic(g)
}

func genEntLogic(g *GenEntLogicContext) error {
	color.Green.Println("Generating...")

	outputDir, err := filepath.Abs(g.Output)
	if err != nil {
		return err
	}

	logicDir := path.Join(outputDir, "internal/logic")

	schemas, err := entc.LoadGraph(g.Schema, &gen.Config{})
	if err != nil {
		return err
	}

	workDir, err := filepath.Abs("./")
	if err != nil {
		return err
	}

	projectCtx, err := ctx.Prepare(workDir)
	if err != nil {
		return err
	}

	successCount := 0

	for _, s := range schemas.Schemas {
		if g.ModelName == s.Name || g.ModelName == "all" {
			color.Blue.Printf("Generating %s...\n", s.Name)
			genCtx := *g
			if g.ModelName == "all" {
				genCtx.GroupName = strings.ToLower(s.Name)
				genCtx.ModelName = s.Name
			}
			// generate logic file
			apiLogicData := GenCRUDData(&genCtx, projectCtx, s)

			for _, v := range apiLogicData {
				logicFilename, err := format.FileNamingFormat(genCtx.Style, v.LogicName)
				if err != nil {
					return err
				}

				// group
				var filename string
				if genCtx.GroupName != "" {
					if err = pathx.MkdirIfNotExist(filepath.Join(logicDir, genCtx.GroupName)); err != nil {
						return err
					}

					filename = filepath.Join(logicDir, genCtx.GroupName, logicFilename+".go")
				} else {
					filename = filepath.Join(logicDir, logicFilename+".go")
				}

				if pathx.FileExists(filename) && !genCtx.Overwrite {
					continue
				}

				err = os.WriteFile(filename, []byte(v.LogicCode), regularPerm)
				if err != nil {
					return err
				}
			}

			// generate api file
			apiData, err := GenApiData(s, genCtx)
			if err != nil {
				return err
			}

			err = pathx.MkdirIfNotExist(filepath.Join(workDir, "desc", strings.ToLower(genCtx.ServiceName)))
			if err != nil {
				return err
			}

			apiFilePath := filepath.Join(workDir, "desc", fmt.Sprintf("%s/%s.api", strings.ToLower(genCtx.ServiceName),
				strcase.ToSnake(genCtx.ModelName)))

			if pathx.FileExists(apiFilePath) && !genCtx.Overwrite {
				return nil
			}

			err = os.WriteFile(apiFilePath, []byte(apiData), regularPerm)
			if err != nil {
				return err
			}

			allApiFile := filepath.Join(workDir, "desc", "all.api")
			allApiData, err := os.ReadFile(allApiFile)
			if err != nil {
				return err
			}
			allApiString := string(allApiData)

			if !strings.Contains(allApiString, fmt.Sprintf("%s.api", strcase.ToSnake(genCtx.ModelName))) {
				allApiString += fmt.Sprintf("\nimport \"%s\"", fmt.Sprintf("./%s/%s.api",
					strings.ToLower(genCtx.ServiceName),
					strcase.ToSnake(genCtx.ModelName)))
			}

			err = os.WriteFile(allApiFile, []byte(allApiString), regularPerm)
			if err != nil {
				return err
			}

			if genCtx.GenApiData {
				prefixStr := ""
				if genCtx.RoutePrefix != "" {
					prefixStr = fmt.Sprintf(" -p %s", genCtx.RoutePrefix)
				}
				_, err := execx.Run(fmt.Sprintf("goctls extra init_code -m %s -t other -n %s%s", genCtx.ModelName, genCtx.ServiceName, prefixStr), g.Output)
				if err != nil {
					color.Red.Printf("the init code of %s already exist, skip... \n", genCtx.ModelName)
				}
			}

			successCount++
		}
	}

	if successCount == 0 {
		color.Red.Printf("Failed. Schema: %s not found. \n", g.ModelName)
	} else {
		color.Green.Println("Generate Ent Logic files for API successfully")
	}

	return nil
}

func GenCRUDData(g *GenEntLogicContext, projectCtx *ctx.ProjectContext, schema *load.Schema) []*ApiLogicData {
	var data []*ApiLogicData
	hasTime, hasUUID, hasCreated, hasUpdated, hasPointy := false, false, false, false, false
	// end string means whether to use \n
	endString := ""
	var packageName string
	if g.GroupName != "" {
		packageName = g.GroupName
	} else {
		packageName = "logic"
	}

	setLogic := strings.Builder{}
	for _, v := range schema.Fields {
		if entx.IsBaseProperty(v.Name) {
			if v.Name == "id" && entx.IsUUIDType(v.Info.Type.String()) {
				g.UseUUID = true
			} else if v.Name == "id" {
				g.IdType = entx.ConvertIdTypeToBaseMessage(v.Info.Type.String())
			}

			if v.Name == "created_at" {
				hasCreated = true
			} else if v.Name == "updated_at" {
				hasUpdated = true
			}
			continue
		} else {
			if entx.IsTimeProperty(v.Info.Type.String()) {
				hasTime = true
				setLogic.WriteString(fmt.Sprintf("\t\t\tSetNotNil%s(pointy.GetTimeMilliPointer(req.%s)).\n", parser.CamelCase(v.Name),
					parser.CamelCase(v.Name)))
			} else if entx.IsUpperProperty(v.Name) {
				if entx.IsUUIDType(v.Info.Type.String()) {
					setLogic.WriteString(fmt.Sprintf("\t\t\tSetNotNil%s(uuidx.ParseUUIDStringToPointer(req.%s)).\n", entx.ConvertSpecificNounToUpper(v.Name),
						parser.CamelCase(v.Name)))
					hasUUID = true
				} else {
					setLogic.WriteString(fmt.Sprintf("\t\t\tSetNotNil%s(req.%s).\n", entx.ConvertSpecificNounToUpper(v.Name),
						parser.CamelCase(v.Name)))
				}
			} else {
				if entx.IsUUIDType(v.Info.Type.String()) {
					setLogic.WriteString(fmt.Sprintf("\t\t\tSetNotNil%s(uuidx.ParseUUIDStringToPointer(req.%s)).\n", parser.CamelCase(v.Name),
						parser.CamelCase(v.Name)))
					hasUUID = true
				} else {
					setLogic.WriteString(fmt.Sprintf("\t\t\tSetNotNil%s(req.%s).\n", parser.CamelCase(v.Name),
						parser.CamelCase(v.Name)))
				}
			}
		}
	}
	setLogic.WriteString("\t\t\tExec(l.ctx)")

	if hasUpdated && hasCreated {
		g.HasCreated = true
		hasPointy = true
	}

	createLogic := bytes.NewBufferString("")
	createLogicTmpl, _ := template.New("create").Parse(createTpl)
	_ = createLogicTmpl.Execute(createLogic, map[string]any{
		"hasTime":      hasTime,
		"hasUUID":      hasUUID,
		"setLogic":     strings.ReplaceAll(setLogic.String(), "Exec", "Save"),
		"modelName":    schema.Name,
		"projectPath":  projectCtx.Path,
		"packageName":  packageName,
		"useUUID":      g.UseUUID, // UUID primary key
		"useI18n":      g.UseI18n,
		"importPrefix": g.ImportPrefix,
		"IdType":       g.IdType,
		"HasCreated":   g.HasCreated,
	})

	data = append(data, &ApiLogicData{
		LogicName: fmt.Sprintf("Create%sLogic", schema.Name),
		LogicCode: createLogic.String(),
	})

	updateLogic := bytes.NewBufferString("")
	updateLogicTmpl, _ := template.New("update").Parse(updateTpl)
	_ = updateLogicTmpl.Execute(updateLogic, map[string]any{
		"hasTime":      hasTime,
		"hasUUID":      hasUUID,
		"setLogic":     setLogic.String(),
		"modelName":    schema.Name,
		"projectPath":  projectCtx.Path,
		"packageName":  packageName,
		"useUUID":      g.UseUUID, // UUID primary key
		"useI18n":      g.UseI18n,
		"importPrefix": g.ImportPrefix,
		"IdType":       g.IdType,
		"HasCreated":   g.HasCreated,
	})

	data = append(data, &ApiLogicData{
		LogicName: fmt.Sprintf("Update%sLogic", schema.Name),
		LogicCode: updateLogic.String(),
	})

	predicateData := strings.Builder{}
	predicateData.WriteString(fmt.Sprintf("\tvar predicates []predicate.%s\n", schema.Name))
	count := 0
	for _, v := range schema.Fields {
		if v.Info.Type.String() == "string" && !strings.Contains(strings.ToLower(v.Name), "uuid") &&
			count < g.SearchKeyNum && !entx.IsBaseProperty(v.Name) {
			camelName := parser.CamelCase(v.Name)
			predicateData.WriteString(fmt.Sprintf("\tif req.%s != nil {\n\t\tpredicates = append(predicates, %s.%sContains(*req.%s))\n\t}\n",
				camelName, strings.ToLower(schema.Name), entx.ConvertSpecificNounToUpper(v.Name), camelName))
			count++
		}
	}
	predicateData.WriteString(fmt.Sprintf("\tdata, err := l.svcCtx.DB.%s.Query().Where(predicates...).Page(l.ctx, req.Page, req.PageSize)",
		schema.Name))

	listData := strings.Builder{}

	for i, v := range schema.Fields {
		if entx.IsBaseProperty(v.Name) {
			continue
		} else {
			nameCamelCase := parser.CamelCase(v.Name)

			if i < (len(schema.Fields) - 1) {
				endString = "\n"
			} else {
				endString = ""
			}

			if entx.IsUUIDType(v.Info.Type.String()) {
				listData.WriteString(fmt.Sprintf("\t\t\t%s:\tpointy.GetPointer(v.%s.String()),%s", nameCamelCase,
					entx.ConvertSpecificNounToUpper(nameCamelCase), endString))
				hasPointy = true
			} else if entx.IsTimeProperty(v.Info.Type.String()) {
				if v.Optional {
					listData.WriteString(fmt.Sprintf("\t\t\t%s:\tpointy.GetUnixMilliPointer(v.%s.UnixMilli()),%s", nameCamelCase,
						entx.ConvertSpecificNounToUpper(nameCamelCase), endString))
				} else {
					listData.WriteString(fmt.Sprintf("\t\t\t%s:\tpointy.GetPointer(v.%s.UnixMilli()),%s", nameCamelCase,
						entx.ConvertSpecificNounToUpper(nameCamelCase), endString))
				}
				hasPointy = true
			} else {
				if entx.IsUpperProperty(v.Name) {
					listData.WriteString(fmt.Sprintf("\t\t\t%s:\t&v.%s,%s", nameCamelCase,
						entx.ConvertSpecificNounToUpper(v.Name), endString))
				} else {
					listData.WriteString(fmt.Sprintf("\t\t\t%s:\t&v.%s,%s", nameCamelCase,
						nameCamelCase, endString))
				}
			}
		}
	}

	getListLogic := bytes.NewBufferString("")
	getListLogicTmpl, _ := template.New("getList").Parse(getListLogicTpl)
	_ = getListLogicTmpl.Execute(getListLogic, map[string]any{
		"predicateData":      predicateData.String(),
		"modelName":          schema.Name,
		"listData":           listData.String(),
		"projectPath":        projectCtx.Path,
		"modelNameLowerCase": strings.ToLower(schema.Name),
		"packageName":        packageName,
		"useUUID":            g.UseUUID,
		"useI18n":            g.UseI18n,
		"importPrefix":       g.ImportPrefix,
		"IdType":             g.IdType,
		"HasCreated":         g.HasCreated,
		"HasPointy":          hasPointy,
	})

	data = append(data, &ApiLogicData{
		LogicName: fmt.Sprintf("Get%sListLogic", schema.Name),
		LogicCode: getListLogic.String(),
	})

	getByIdLogic := bytes.NewBufferString("")
	getByIdLogicTmpl, _ := template.New("getById").Parse(getByIdLogicTpl)
	_ = getByIdLogicTmpl.Execute(getByIdLogic, map[string]any{
		"modelName":          schema.Name,
		"listData":           strings.Replace(listData.String(), "v.", "data.", -1),
		"projectPath":        projectCtx.Path,
		"modelNameLowerCase": strings.ToLower(schema.Name),
		"packageName":        packageName,
		"useUUID":            g.UseUUID,
		"useI18n":            g.UseI18n,
		"importPrefix":       g.ImportPrefix,
		"IdType":             g.IdType,
		"HasCreated":         g.HasCreated,
		"HasPointy":          hasPointy,
	})

	data = append(data, &ApiLogicData{
		LogicName: fmt.Sprintf("Get%sByIdLogic", schema.Name),
		LogicCode: getByIdLogic.String(),
	})

	deleteLogic := bytes.NewBufferString("")
	deleteLogicTmpl, _ := template.New("delete").Parse(deleteLogicTpl)
	_ = deleteLogicTmpl.Execute(deleteLogic, map[string]any{
		"modelName":          schema.Name,
		"modelNameLowerCase": strings.ToLower(schema.Name),
		"projectPath":        projectCtx.Path,
		"packageName":        packageName,
		"useUUID":            g.UseUUID,
		"useI18n":            g.UseI18n,
		"importPrefix":       g.ImportPrefix,
		"IdType":             g.IdType,
		"HasCreated":         g.HasCreated,
	})

	data = append(data, &ApiLogicData{
		LogicName: fmt.Sprintf("Delete%sLogic", schema.Name),
		LogicCode: deleteLogic.String(),
	})

	return data
}

func GenApiData(schema *load.Schema, ctx GenEntLogicContext) (string, error) {
	infoData := strings.Builder{}
	listData := strings.Builder{}
	searchKeyNum := ctx.SearchKeyNum
	var data string
	var hasRoutePrefix bool
	if ctx.RoutePrefix != "" {
		hasRoutePrefix = true
	} else {
		hasRoutePrefix = false
	}

	for _, v := range schema.Fields {
		if entx.IsBaseProperty(v.Name) {
			continue
		}

		var structData string

		jsonTag, err := format.FileNamingFormat(ctx.JSONStyle, v.Name)
		if err != nil {
			return "", err
		}

		pointerStr := ""
		optionalStr := ""
		if !entx.IsPageProperty(strings.ToLower(v.Name)) {
			pointerStr = "*"
			optionalStr = "optional"
		}

		fieldComment := parser.CamelCase(v.Name)

		if v.Comment != "" {
			fieldComment = v.Comment
		}

		structData = fmt.Sprintf("\n\n        // %s \n        %s  %s%s `json:\"%s,%s\"`",
			fieldComment,
			parser.CamelCase(v.Name),
			pointerStr,
			entx.ConvertEntTypeToGotypeInSingleApi(v.Info.Type.String()),
			jsonTag, optionalStr)

		infoData.WriteString(structData)

		if v.Info.Type.String() == "string" && searchKeyNum > 0 {
			listData.WriteString(structData)
			searchKeyNum--
		}
	}

	modelChineseName := ctx.ModelName
	if ctx.ModelChineseName != "" {
		modelChineseName = ctx.ModelChineseName
	}

	modelEnglishName := strings.Replace(strcase.ToSnake(ctx.ModelName), "_", " ", -1)
	if ctx.ModelEnglishName != "" {
		modelEnglishName = ctx.ModelEnglishName
	}

	apiTemplateData := bytes.NewBufferString("")
	apiTmpl, _ := template.New("entApiTpl").Parse(apiTpl)
	logx.Must(apiTmpl.Execute(apiTemplateData, map[string]any{
		"infoData":         infoData.String(),
		"modelName":        ctx.ModelName,
		"groupName":        ctx.GroupName,
		"modelNameSnake":   strcase.ToSnake(ctx.ModelName),
		"listData":         listData.String(),
		"apiServiceName":   strcase.ToCamel(ctx.ServiceName),
		"useUUID":          ctx.UseUUID,
		"hasRoutePrefix":   hasRoutePrefix,
		"routePrefix":      ctx.RoutePrefix,
		"IdType":           ctx.IdType,
		"HasCreated":       ctx.HasCreated,
		"IdTypeLower":      strings.ToLower(ctx.IdType),
		"modelChineseName": modelChineseName,
		"modelEnglishName": modelEnglishName,
	}))
	data = apiTemplateData.String()

	return data, nil
}
