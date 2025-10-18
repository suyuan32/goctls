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

	"github.com/gookit/color"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/entc/load"
	"github.com/iancoleman/strcase"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/suyuan32/goctls/rpc/parser"
	"github.com/suyuan32/goctls/util/ctx"
	"github.com/suyuan32/goctls/util/entx"
	"github.com/suyuan32/goctls/util/format"
	"github.com/suyuan32/goctls/util/pathx"
	"github.com/suyuan32/goctls/util/protox"
)

const regularPerm = 0o666

type RpcLogicData struct {
	LogicName string
	LogicCode string
}

type GenEntLogicContext struct {
	Schema          string
	Output          string
	ServiceName     string
	ProjectName     string
	Style           string
	ProtoFieldStyle string
	ModelName       string
	Multiple        bool
	SearchKeyNum    int
	ModuleName      string
	GroupName       string
	UseUUID         bool
	UseI18n         bool
	ProtoOut        string
	ImportPrefix    string
	Overwrite       bool
	IdType          string // default is empty, if ID belongs types of Uint64 and string, use other base message
	HasCreated      bool   // If true means have created and updated field
	SplitTimeField  bool
}

func (g GenEntLogicContext) Validate() error {
	if g.Schema == "" {
		return errors.New("the schema dir cannot be empty ")
	} else if !strings.HasSuffix(g.Schema, "schema") {
		return errors.New("please input correct schema directory e.g. ./ent/schema ")
	} else if g.ServiceName == "" {
		return errors.New("please set the service name via --service_name")
	} else if g.ModuleName == g.ProjectName {
		return errors.New("do not set the module name if it is the same as project name ")
	}
	return nil
}

// GenEntLogic generates the ent CRUD logic files of the rpc service.
func GenEntLogic(g *GenEntLogicContext) error {
	return genEntLogic(g)
}

func genEntLogic(g *GenEntLogicContext) error {
	color.Green.Println("Generating...")
	outputDir, err := filepath.Abs(g.Output)
	if err != nil {
		return err
	}

	var logicDir string

	if g.Multiple {
		logicDir = path.Join(outputDir, "internal/logic", g.ServiceName)
		err = pathx.MkdirIfNotExist(logicDir)
		if err != nil {
			return err
		}
	} else {
		logicDir = path.Join(outputDir, "internal/logic")
	}

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
			rpcLogicData := GenCRUDData(&genCtx, projectCtx, s)

			for _, v := range rpcLogicData {
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

			// generate proto file
			protoMessage, protoFunctions, err := GenProtoData(s, genCtx)
			if err != nil {
				return err
			}

			var protoFileName string
			if genCtx.ProtoOut == "" {
				protoFileName = filepath.Join(outputDir, genCtx.ProjectName+".proto")
				if !pathx.FileExists(protoFileName) {
					continue
				}
			} else {
				if g.ModelName == "all" {
					genCtx.ProtoOut = strings.ReplaceAll(genCtx.ProtoOut, "all", genCtx.GroupName)
				}
				protoFileName, err = filepath.Abs(genCtx.ProtoOut)
				if err != nil {
					return err
				}
				if !pathx.FileExists(protoFileName) || genCtx.Overwrite {
					err = os.WriteFile(protoFileName, []byte(fmt.Sprintf("syntax = \"proto3\";\n\nservice %s {\n}",
						strcase.ToCamel(genCtx.ServiceName))), os.ModePerm)
					if err != nil {
						return fmt.Errorf("failed to create proto file : %s", err.Error())
					}
				}
			}

			protoFileData, err := os.ReadFile(protoFileName)
			if err != nil {
				return err
			}

			protoDataString := string(protoFileData)

			if strings.Contains(protoDataString, protoMessage) || strings.Contains(protoDataString, protoFunctions) {
				continue
			}

			// generate new proto file
			newProtoData := strings.Builder{}
			serviceBeginIndex, _, serviceEndIndex := protox.FindBeginEndOfService(protoDataString, strcase.ToCamel(genCtx.ServiceName))
			if serviceBeginIndex == -1 {
				continue
			}
			newProtoData.WriteString(protoDataString[:serviceBeginIndex-1])
			newProtoData.WriteString(fmt.Sprintf("\n// %s message\n\n", genCtx.ModelName))
			newProtoData.WriteString(fmt.Sprintf("%s\n", protoMessage))
			newProtoData.WriteString(protoDataString[serviceBeginIndex-1 : serviceEndIndex-1])
			newProtoData.WriteString(fmt.Sprintf("\n\n  // %s management\n", genCtx.ModelName))
			newProtoData.WriteString(fmt.Sprintf("%s\n", protoFunctions))
			newProtoData.WriteString(protoDataString[serviceEndIndex-1:])

			err = os.WriteFile(protoFileName, []byte(newProtoData.String()), regularPerm)
			if err != nil {
				return err
			}
			successCount++
		}
	}

	if successCount == 0 {
		color.Red.Printf("Failed. Schema: %s not found. \n", g.ModelName)
	} else {
		color.Green.Println("Generate Ent Logic files for RPC successfully")
	}
	return nil
}

func GenCRUDData(g *GenEntLogicContext, projectCtx *ctx.ProjectContext, schema *load.Schema) []*RpcLogicData {
	// ent field
	var tmpField = &gen.Field{}
	var data []*RpcLogicData
	hasTime, hasUUID, hasSingle, NoNormalField, hasPointy, hasCreated, hasUpdated := false, false, false, true, false, false, false
	// end string means whether to use \n
	endString := ""
	var packageName string
	if g.GroupName != "" {
		packageName = g.GroupName
	} else {
		packageName = "logic"
	}

	singleSets := []string{}

	setLogic := strings.Builder{}
	for _, v := range schema.Fields {
		camelName := parser.CamelCase(v.Name)
		tmpField.Name = v.Name

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
		} else if entx.IsOnlyEntType(v.Info.Type.String()) {
			singleSets = append(singleSets, fmt.Sprintf("\tif in.%s != nil {\n\t\tquery.SetNotNil%s(pointy.GetPointer(%s(*in.%s)))\n\t}\n",
				tmpField.StructField(),
				camelName,
				v.Info.Type.String(),
				camelName),
			)
			hasSingle = true
		} else {
			if entx.IsTimeProperty(v.Info.Type.String()) {
				hasTime = true
				setLogic.WriteString(fmt.Sprintf("\t\t\tSetNotNil%s(pointy.GetTimeMilliPointer(in.%s)).\n", tmpField.StructField(),
					camelName))
			} else {
				if entx.IsGoTypeNotPrototype(v.Info.Type.String()) {
					if v.Info.Type.String() == "[16]byte" {
						setLogic.WriteString(fmt.Sprintf("\t\t\tSetNotNil%s(uuidx.ParseUUIDStringToPointer(in.%s)).\n", tmpField.StructField(),
							camelName))
						hasUUID = true
					} else {
						singleSets = append(singleSets, fmt.Sprintf("\tif in.%s != nil {\n\t\tquery.SetNotNil%s(pointy.GetPointer(%s(*in.%s)))\n\t}\n",
							camelName,
							tmpField.StructField(),
							v.Info.Type.String(),
							camelName),
						)
						hasSingle = true
					}
				} else {
					setLogic.WriteString(fmt.Sprintf("\t\t\tSetNotNil%s(in.%s).\n", tmpField.StructField(),
						camelName))
				}
			}
		}
	}

	if hasSingle {
		tmp := setLogic.String()
		tmp = strings.TrimSuffix(tmp, ".\n")
		setLogic.Reset()
		setLogic.WriteString(tmp)
		setLogic.WriteString("\n\n")

		for _, v := range singleSets {
			setLogic.WriteString(v)
		}

		setLogic.WriteString("\n\tresult, err := query.Exec(l.ctx)")
	} else {
		setLogic.WriteString("\t\t\tExec(l.ctx)")
	}

	if strings.HasPrefix(setLogic.String(), "\t\t\tSet") {
		NoNormalField = false
	}

	if strings.Contains(setLogic.String(), "pointy") {
		hasPointy = true
	}

	// judge if generate time field
	if hasCreated && hasUpdated {
		g.HasCreated = true
	}

	createLogic := bytes.NewBufferString("")
	createLogicTmpl, _ := template.New("create").Parse(createTpl)
	_ = createLogicTmpl.Execute(createLogic, map[string]any{
		"hasTime":       hasTime,
		"hasUUID":       hasUUID,
		"setLogic":      strings.ReplaceAll(setLogic.String(), "Exec", "Save"),
		"modelName":     schema.Name,
		"projectName":   strings.ToLower(g.ProjectName),
		"projectPath":   projectCtx.Path,
		"packageName":   packageName,
		"useUUID":       g.UseUUID, // UUID primary key
		"useI18n":       g.UseI18n,
		"importPrefix":  g.ImportPrefix,
		"hasSingle":     hasSingle,
		"noNormalField": !NoNormalField,
		"hasPointy":     hasPointy,
		"IdType":        g.IdType,
		"HasCreated":    g.HasCreated,
	})

	data = append(data, &RpcLogicData{
		LogicName: fmt.Sprintf("Create%sLogic", schema.Name),
		LogicCode: createLogic.String(),
	})

	updateLogic := bytes.NewBufferString("")
	updateLogicTmpl, _ := template.New("update").Parse(updateTpl)
	_ = updateLogicTmpl.Execute(updateLogic, map[string]any{
		"hasTime":       hasTime,
		"hasUUID":       hasUUID,
		"setLogic":      strings.Replace(setLogic.String(), "result,", "", 1),
		"modelName":     schema.Name,
		"projectName":   strings.ToLower(g.ProjectName),
		"projectPath":   projectCtx.Path,
		"packageName":   packageName,
		"useUUID":       g.UseUUID, // UUID primary key
		"useI18n":       g.UseI18n,
		"importPrefix":  g.ImportPrefix,
		"hasSingle":     hasSingle,
		"noNormalField": !NoNormalField,
		"hasPointy":     hasPointy,
		"IdType":        g.IdType,
		"HasCreated":    g.HasCreated,
	})

	data = append(data, &RpcLogicData{
		LogicName: fmt.Sprintf("Update%sLogic", schema.Name),
		LogicCode: updateLogic.String(),
	})

	predicateData := strings.Builder{}
	predicateData.WriteString(fmt.Sprintf("\tvar predicates []predicate.%s\n", schema.Name))
	count := 0
	for _, v := range schema.Fields {
		if v.Name == "id" || v.Name == "tenant_id" || count >= g.SearchKeyNum {
			continue
		}

		camelName := parser.CamelCase(v.Name)
		tmpField.Name = v.Name

		if v.Info.Type.String() == "string" && !strings.Contains(strings.ToLower(v.Name), "uuid") {
			predicateData.WriteString(fmt.Sprintf("\tif in.%s != nil {\n\t\tpredicates = append(predicates, %s.%sContains(*in.%s))\n\t}\n",
				camelName, strings.ToLower(schema.Name), tmpField.StructField(), camelName))
			count++
		} else if entx.IsTimeProperty(v.Info.Type.String()) {
			if g.SplitTimeField {
				predicateData.WriteString(fmt.Sprintf("\tif in.%sBegin != nil {\n\t\tpredicates = append(predicates, %s.%sGTE(time.UnixMilli(*in.%sBegin)))\n\t}\n",
					camelName, strings.ToLower(schema.Name), tmpField.StructField(), camelName))
				predicateData.WriteString(fmt.Sprintf("\tif in.%sEnd != nil {\n\t\tpredicates = append(predicates, %s.%sLTE(time.UnixMilli(*in.%sEnd)))\n\t}\n",
					camelName, strings.ToLower(schema.Name), tmpField.StructField(), camelName))
			} else {
				predicateData.WriteString(fmt.Sprintf("\tif in.%s != nil {\n\t\tpredicates = append(predicates, %s.%sGTE(time.UnixMilli(*in.%s)))\n\t}\n",
					camelName, strings.ToLower(schema.Name), tmpField.StructField(), camelName))
			}
			count++
		} else {
			if entx.IsGoTypeNotPrototype(v.Info.Type.String()) {
				if v.Info.Type.String() == "[16]byte" {
					predicateData.WriteString(fmt.Sprintf("\tif in.%s != nil {\n\t\tpredicates = append(predicates, %s.%sEQ(uuidx.ParseUUIDString(*in.%s)))\n\t}\n",
						camelName, strings.ToLower(schema.Name), tmpField.StructField(), camelName))
				} else {
					predicateData.WriteString(fmt.Sprintf("\tif in.%s != nil {\n\t\tpredicates = append(predicates, %s.%sEQ(%s(*in.%s)))\n\t}\n",
						camelName, strings.ToLower(schema.Name), tmpField.StructField(), v.Info.Type.String(), camelName))
				}
			} else if entx.IsOnlyEntType(v.Info.Type.String()) {
				predicateData.WriteString(fmt.Sprintf("\tif in.%s != nil {\n\t\tpredicates = append(predicates, %s.%sEQ(%s(*in.%s)))\n\t}\n",
					camelName, strings.ToLower(schema.Name), tmpField.StructField(), v.Info.Type.String(), camelName))
			} else {
				predicateData.WriteString(fmt.Sprintf("\tif in.%s != nil {\n\t\tpredicates = append(predicates, %s.%sEQ(*in.%s))\n\t}\n",
					camelName, strings.ToLower(schema.Name), tmpField.StructField(), camelName))
			}
			count++
		}
	}
	predicateData.WriteString(fmt.Sprintf("\tresult, err := l.svcCtx.DB.%s.Query().Where(predicates...).Page(l.ctx, in.Page, in.PageSize)",
		schema.Name))

	listData := strings.Builder{}

	for i, v := range schema.Fields {
		tmpField.Name = v.Name
		if entx.IsBaseProperty(v.Name) {
			continue
		} else {
			camelName := parser.CamelCase(v.Name)

			if i < (len(schema.Fields) - 1) {
				endString = "\n"
			} else {
				endString = ""
			}

			if entx.IsUUIDType(v.Info.Type.String()) {
				listData.WriteString(fmt.Sprintf("\t\t\t%s:\tpointy.GetPointer(v.%s.String()),%s", camelName,
					tmpField.StructField(), endString))
			} else if entx.IsOnlyEntType(v.Info.Type.String()) {
				listData.WriteString(fmt.Sprintf("\t\t\t%s:\tpointy.GetPointer(%s(v.%s)),%s", camelName,
					entx.ConvertOnlyEntTypeToGoType(v.Info.Type.String()),
					tmpField.StructField(), endString))
			} else if entx.IsTimeProperty(v.Info.Type.String()) {
				if v.Optional {
					listData.WriteString(fmt.Sprintf("\t\t\t%s:\tpointy.GetUnixMilliPointer(v.%s.UnixMilli()),%s", camelName,
						tmpField.StructField(), endString))
				} else {
					listData.WriteString(fmt.Sprintf("\t\t\t%s:\tpointy.GetPointer(v.%s.UnixMilli()),%s", camelName,
						tmpField.StructField(), endString))
				}
			} else {
				if entx.IsGoTypeNotPrototype(v.Info.Type.String()) {
					listData.WriteString(fmt.Sprintf("\t\t\t%s:\tpointy.GetPointer(%s(v.%s)),%s", camelName,
						entx.ConvertEntTypeToGotype(v.Info.Type.String()), tmpField.StructField(), endString))
				} else {
					listData.WriteString(fmt.Sprintf("\t\t\t%s:\t&v.%s,%s", camelName,
						tmpField.StructField(), endString))
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
		"projectName":        strings.ToLower(g.ProjectName),
		"projectPath":        projectCtx.Path,
		"modelNameLowerCase": strings.ToLower(schema.Name),
		"packageName":        packageName,
		"useUUID":            g.UseUUID,
		"useI18n":            g.UseI18n,
		"importPrefix":       g.ImportPrefix,
		"IdType":             g.IdType,
		"HasCreated":         g.HasCreated,
		"hasTime":            strings.Contains(predicateData.String(), "time."),
		"hasUUID":            strings.Contains(predicateData.String(), "uuidx."),
		"hasPointy":          strings.Contains(listData.String(), "pointy"),
	})

	data = append(data, &RpcLogicData{
		LogicName: fmt.Sprintf("Get%sListLogic", schema.Name),
		LogicCode: getListLogic.String(),
	})

	listDataConv := strings.Replace(listData.String(), "v.", "result.", -1)
	listDataConv = strings.ReplaceAll(listDataConv, "\t\t\t", "\t\t")

	getByIdLogic := bytes.NewBufferString("")
	getByIdLogicTmpl, _ := template.New("getById").Parse(getByIdLogicTpl)
	_ = getByIdLogicTmpl.Execute(getByIdLogic, map[string]any{
		"modelName":          schema.Name,
		"listData":           listDataConv,
		"projectName":        strings.ToLower(g.ProjectName),
		"projectPath":        projectCtx.Path,
		"modelNameLowerCase": strings.ToLower(schema.Name),
		"packageName":        packageName,
		"useUUID":            g.UseUUID,
		"useI18n":            g.UseI18n,
		"importPrefix":       g.ImportPrefix,
		"IdType":             g.IdType,
		"HasCreated":         g.HasCreated,
	})

	data = append(data, &RpcLogicData{
		LogicName: fmt.Sprintf("Get%sByIdLogic", schema.Name),
		LogicCode: getByIdLogic.String(),
	})

	deleteLogic := bytes.NewBufferString("")
	deleteLogicTmpl, _ := template.New("delete").Parse(deleteLogicTpl)
	_ = deleteLogicTmpl.Execute(deleteLogic, map[string]any{
		"modelName":          schema.Name,
		"modelNameLowerCase": strings.ToLower(schema.Name),
		"projectName":        strings.ToLower(g.ProjectName),
		"projectPath":        projectCtx.Path,
		"packageName":        packageName,
		"useUUID":            g.UseUUID,
		"useI18n":            g.UseI18n,
		"importPrefix":       g.ImportPrefix,
		"IdType":             g.IdType,
		"HasCreated":         g.HasCreated,
	})

	data = append(data, &RpcLogicData{
		LogicName: fmt.Sprintf("Delete%sLogic", schema.Name),
		LogicCode: deleteLogic.String(),
	})

	return data
}

func GenProtoData(schema *load.Schema, g GenEntLogicContext) (string, string, error) {
	var protoMessage strings.Builder
	schemacamelName := parser.CamelCase(schema.Name)
	// hasStatus means it has status field
	hasStatus := false
	// end string means whether to use \n
	endString := ""
	// info message
	idString, _ := format.FileNamingFormat(g.ProtoFieldStyle, "id")
	index := 2

	// gen time field
	timeField := ""
	if g.HasCreated {
		createString, _ := format.FileNamingFormat(g.ProtoFieldStyle, "created_at")
		updateString, _ := format.FileNamingFormat(g.ProtoFieldStyle, "updated_at")

		timeField = fmt.Sprintf("  optional int64 %s = 2;\n  optional int64 %s = 3;\n", createString, updateString)
		index = 4
	}

	protoMessage.WriteString(fmt.Sprintf("message %sInfo {\n  optional %s %s = 1;\n%s",
		schemacamelName, entx.ConvertIDType(g.UseUUID, strings.ToLower(g.IdType)), idString, timeField))

	for i, v := range schema.Fields {
		var fieldComment string
		if v.Comment != "" {
			fieldComment = fmt.Sprintf("  // %s\n", v.Comment)
		}

		if entx.IsBaseProperty(v.Name) {
			continue
		} else if v.Name == "status" {
			statusString, _ := format.FileNamingFormat(g.ProtoFieldStyle, v.Name)
			protoMessage.WriteString(fmt.Sprintf("%s  optional uint32 %s = %d;\n", fieldComment, statusString, index))
			hasStatus = true
			index++
		} else {
			if i < (len(schema.Fields) - 1) {
				endString = "\n"
			} else {
				endString = ""
			}

			formatedString, _ := format.FileNamingFormat(g.ProtoFieldStyle, v.Name)
			if entx.IsTimeProperty(v.Info.Type.String()) {
				protoMessage.WriteString(fmt.Sprintf("%s  optional int64  %s = %d;%s", fieldComment, formatedString, index, endString))
			} else {
				protoMessage.WriteString(fmt.Sprintf("%s  optional %s %s = %d;%s", fieldComment, entx.ConvertEntTypeToProtoType(v.Info.Type.String()),
					formatedString, index, endString))
			}

			index++
		}
	}

	protoMessage.WriteString("\n}\n\n")

	// List message
	totalString, _ := format.FileNamingFormat(g.ProtoFieldStyle, "total")
	dataString, _ := format.FileNamingFormat(g.ProtoFieldStyle, "data")
	protoMessage.WriteString(fmt.Sprintf("message %sListResp {\n  uint64 %s = 1;\n  repeated %sInfo %s = 2;\n}\n\n",
		schemacamelName, totalString, schemacamelName, dataString))

	// List Request message
	pageString, _ := format.FileNamingFormat(g.ProtoFieldStyle, "page")
	pageSizeString, _ := format.FileNamingFormat(g.ProtoFieldStyle, "page_size")
	protoMessage.WriteString(fmt.Sprintf("message %sListReq {\n  uint64 %s = 1;\n  uint64 %s = 2;\n",
		schemacamelName, pageString, pageSizeString))
	count := 0
	index = 3

	for _, v := range schema.Fields {
		if v.Name == "id" || v.Name == "tenant_id" || count >= g.SearchKeyNum {
			continue
		}

		if v.Info.Type.String() == "string" && !strings.Contains(strings.ToLower(v.Name), "uuid") {
			formatedString, _ := format.FileNamingFormat(g.ProtoFieldStyle, v.Name)
			protoMessage.WriteString(fmt.Sprintf("  optional %s %s = %d;\n", entx.ConvertEntTypeToProtoType(v.Info.Type.String()),
				formatedString, index))
			index++
			count++
		} else if entx.IsTimeProperty(v.Info.Type.String()) {
			formatedString, _ := format.FileNamingFormat(g.ProtoFieldStyle, v.Name)
			if g.SplitTimeField {
				protoMessage.WriteString(fmt.Sprintf("  optional int64 %s_begin = %d;\n", formatedString, index))
				index++
				count++
				protoMessage.WriteString(fmt.Sprintf("  optional int64 %s_end = %d;\n", formatedString, index))
				index++
				count++
			} else {
				protoMessage.WriteString(fmt.Sprintf("  optional int64 %s = %d;\n", formatedString, index))
				index++
				count++
			}
		} else {
			formatedString, _ := format.FileNamingFormat(g.ProtoFieldStyle, v.Name)
			protoMessage.WriteString(fmt.Sprintf("  optional %s %s = %d;\n", entx.ConvertEntTypeToProtoType(v.Info.Type.String()),
				formatedString, index))
			index++
			count++
		}
	}

	protoMessage.WriteString("}\n")

	// group
	var groupName string
	if g.GroupName != "" {
		groupName = fmt.Sprintf("  // group: %s\n", g.GroupName)
	} else {
		groupName = ""
	}

	protoRpcFunction := bytes.NewBufferString("")
	protoTmpl, err := template.New("proto").Parse(protoTpl)
	err = protoTmpl.Execute(protoRpcFunction, map[string]any{
		"modelName":  schema.Name,
		"groupName":  groupName,
		"useUUID":    g.UseUUID,
		"hasStatus":  hasStatus,
		"IdType":     g.IdType,
		"HasCreated": g.HasCreated,
	})

	if err != nil {
		logx.Error(err)
		return "", "", err
	}

	return protoMessage.String(), protoRpcFunction.String(), nil
}
