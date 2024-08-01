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

package proto

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gookit/color"
	"github.com/suyuan32/goctls/rpc/execx"

	"github.com/emicklei/proto"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/iancoleman/strcase"

	"github.com/suyuan32/goctls/rpc/parser"
	"github.com/suyuan32/goctls/util/ctx"
	"github.com/suyuan32/goctls/util/entx"
	"github.com/suyuan32/goctls/util/format"
	"github.com/suyuan32/goctls/util/pathx"
	"github.com/suyuan32/goctls/util/protox"
)

const regularPerm = 0o666

// GenLogicByProtoContext describe the data used for logic generation with proto file
type GenLogicByProtoContext struct {
	ProtoDir         string
	OutputDir        string
	APIServiceName   string
	RPCServiceName   string
	RPCPbPackageName string
	Style            string
	ModelName        string
	RpcName          string
	GrpcPackage      string
	UseUUID          bool
	Multiple         bool
	JSONStyle        string
	UseI18n          bool
	ImportPrefix     string
	OptionalService  bool
	GenApiData       bool
	Overwrite        bool
	RoutePrefix      string
	IdType           string
	HasCreated       bool
}

func (g GenLogicByProtoContext) Validate() error {
	if g.APIServiceName == "" {
		return errors.New("please set the API service name via --api_service_name ")
	} else if g.RPCServiceName == "" {
		return errors.New("please set the RPC service name via --rpc_service_name ")
	} else if g.ProtoDir == "" {
		return errors.New("please set the proto dir via --proto ")
	} else if !strings.HasSuffix(g.ProtoDir, "proto") {
		return errors.New("please set the correct proto file ")
	} else if g.ModelName == "" {
		return errors.New("please set the model name via --model ")
	} else if g.RpcName == "" {
		return errors.New("please set the RPC name via --rpc_name ")
	}
	return nil
}

type ApiLogicData struct {
	LogicName string
	LogicCode string
}

func GenLogicByProto(p *GenLogicByProtoContext) error {
	outputDir, err := filepath.Abs(p.OutputDir)
	if err != nil {
		return err
	}

	// convert api and  rpc service name style
	p.APIServiceName = strcase.ToCamel(p.APIServiceName)
	p.RPCServiceName = strcase.ToCamel(p.RPCServiceName)

	logicDir := path.Join(outputDir, "internal/logic")

	protoParser := parser.NewDefaultProtoParser()
	protoData, err := protoParser.Parse(p.ProtoDir, p.Multiple)
	if err != nil {
		return err
	}

	p.RPCPbPackageName = protoData.PbPackage

	protox.ProtoField = &protox.ProtoFieldData{}

	workDir, err := filepath.Abs("./")
	if err != nil {
		return err
	}

	projectCtx, err := ctx.Prepare(workDir)
	if err != nil {
		return err
	}

	models := strings.Split(p.ModelName, ",")

	for _, v := range models {
		genCtx := *p

		genCtx.ModelName = v

		color.Blue.Printf("Generating %s ...\n", v)

		// generate logic file
		apiLogicData := GenCRUDData(&genCtx, &protoData, projectCtx)

		for _, v := range apiLogicData {
			logicFilename, err := format.FileNamingFormat(genCtx.Style, v.LogicName)
			if err != nil {
				return err
			}

			filename := filepath.Join(logicDir, strings.ToLower(genCtx.ModelName), logicFilename+".go")
			if err = pathx.MkdirIfNotExist(filepath.Join(logicDir, strings.ToLower(genCtx.ModelName))); err != nil {
				return err
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
		apiData, err := GenApiData(&genCtx, &protoData)
		if err != nil {
			return err
		}

		err = pathx.MkdirIfNotExist(filepath.Join(workDir, "desc", strings.ToLower(genCtx.RPCServiceName)))
		if err != nil {
			return err
		}

		apiFilePath := filepath.Join(workDir, "desc", fmt.Sprintf("%s/%s.api", strings.ToLower(genCtx.RPCServiceName),
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
				strings.ToLower(genCtx.RPCServiceName),
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
			_, err := execx.Run(fmt.Sprintf("goctls extra init_code -m %s -t other -n %s%s", genCtx.ModelName, genCtx.RPCServiceName, prefixStr), genCtx.OutputDir)
			if err != nil {
				color.Red.Printf("the init code of %s already exist, skip... \n", genCtx.ModelName)
			}
		}
	}

	color.Green.Println("Generate logic files from proto successfully")

	return nil
}

func GenCRUDData(ctx *GenLogicByProtoContext, p *parser.Proto, projectCtx *ctx.ProjectContext) []*ApiLogicData {
	var data []*ApiLogicData
	var setLogic string

	for _, v := range p.Message {
		if strings.Contains(v.Name, ctx.ModelName) {
			if fmt.Sprintf("%sInfo", ctx.ModelName) == v.Name {
				setLogic = genSetLogic(v.Message, ctx)

				createLogic := bytes.NewBufferString("")
				createLogicTmpl, _ := template.New("create").Parse(createTpl)
				logx.Must(createLogicTmpl.Execute(createLogic, map[string]any{
					"setLogic":           setLogic,
					"modelName":          ctx.ModelName,
					"modelNameLowerCase": strings.ToLower(ctx.ModelName),
					"projectPackage":     projectCtx.Path,
					"rpcPackage":         ctx.GrpcPackage,
					"rpcName":            ctx.RpcName,
					"rpcPbPackageName":   ctx.RPCPbPackageName,
					"useUUID":            ctx.UseUUID,
					"useI18n":            ctx.UseI18n,
					"importPrefix":       ctx.ImportPrefix,
					"optionalService":    ctx.OptionalService,
					"IdType":             ctx.IdType,
					"HasCreated":         ctx.HasCreated,
				}))

				data = append(data, &ApiLogicData{
					LogicName: fmt.Sprintf("Create%sLogic", ctx.ModelName),
					LogicCode: createLogic.String(),
				})

				updateLogic := bytes.NewBufferString("")
				updateLogicTmpl, _ := template.New("update").Parse(updateTpl)
				logx.Must(updateLogicTmpl.Execute(updateLogic, map[string]any{
					"setLogic":           setLogic,
					"modelName":          ctx.ModelName,
					"modelNameLowerCase": strings.ToLower(ctx.ModelName),
					"projectPackage":     projectCtx.Path,
					"rpcPackage":         ctx.GrpcPackage,
					"rpcName":            ctx.RpcName,
					"rpcPbPackageName":   ctx.RPCPbPackageName,
					"useUUID":            ctx.UseUUID,
					"useI18n":            ctx.UseI18n,
					"importPrefix":       ctx.ImportPrefix,
					"optionalService":    ctx.OptionalService,
					"IdType":             ctx.IdType,
					"HasCreated":         ctx.HasCreated,
				}))

				data = append(data, &ApiLogicData{
					LogicName: fmt.Sprintf("Update%sLogic", ctx.ModelName),
					LogicCode: updateLogic.String(),
				})

				// delete logic
				deleteLogic := bytes.NewBufferString("")
				deleteLogicTmpl, _ := template.New("delete").Parse(deleteLogicTpl)
				logx.Must(deleteLogicTmpl.Execute(deleteLogic, map[string]any{
					"setLogic":           setLogic,
					"modelName":          ctx.ModelName,
					"modelNameLowerCase": strings.ToLower(ctx.ModelName),
					"projectPackage":     projectCtx.Path,
					"rpcPackage":         ctx.GrpcPackage,
					"rpcName":            ctx.RpcName,
					"rpcPbPackageName":   ctx.RPCPbPackageName,
					"useUUID":            ctx.UseUUID,
					"useI18n":            ctx.UseI18n,
					"importPrefix":       ctx.ImportPrefix,
					"optionalService":    ctx.OptionalService,
					"IdType":             ctx.IdType,
					"HasCreated":         ctx.HasCreated,
				}))

				data = append(data, &ApiLogicData{
					LogicName: fmt.Sprintf("Delete%sLogic", ctx.ModelName),
					LogicCode: deleteLogic.String(),
				})
			}

			if fmt.Sprintf("%sListReq", ctx.ModelName) == v.Name {
				searchLogic := strings.Builder{}
				for _, field := range v.Elements {
					field.Accept(protox.MessageVisitor{})
					if protox.ProtoField.Name == "page" || protox.ProtoField.Name == "page_size" {
						continue
					}
					searchLogic.WriteString(fmt.Sprintf("\n\t\t\t%s: req.%s,", parser.CamelCase(protox.ProtoField.Name),
						parser.CamelCase(protox.ProtoField.Name)))
				}

				if setLogic == "" {
					for _, m := range p.Message {
						if strings.Contains(m.Name, ctx.ModelName) {
							if fmt.Sprintf("%sInfo", ctx.ModelName) == m.Name {
								setLogic = genSetLogic(m.Message, ctx)
							}
						}
					}
				}

				getListLogic := bytes.NewBufferString("")
				getListLogicTmpl, _ := template.New("getList").Parse(getListLogicTpl)
				logx.Must(getListLogicTmpl.Execute(getListLogic, map[string]any{
					"setLogic":           strings.Replace(setLogic, "req.", "v.", -1),
					"modelName":          ctx.ModelName,
					"modelNameLowerCase": strings.ToLower(ctx.ModelName),
					"projectPackage":     projectCtx.Path,
					"rpcPackage":         ctx.GrpcPackage,
					"rpcName":            ctx.RpcName,
					"rpcPbPackageName":   ctx.RPCPbPackageName,
					"searchKeys":         searchLogic.String(),
					"useUUID":            ctx.UseUUID,
					"useI18n":            ctx.UseI18n,
					"importPrefix":       ctx.ImportPrefix,
					"optionalService":    ctx.OptionalService,
					"IdType":             ctx.IdType,
					"HasCreated":         ctx.HasCreated,
				}))

				data = append(data, &ApiLogicData{
					LogicName: fmt.Sprintf("Get%sListLogic", ctx.ModelName),
					LogicCode: getListLogic.String(),
				})

				getByIdLogic := bytes.NewBufferString("")
				getByIdLogicTmpl, _ := template.New("getById").Parse(getByIdLogicTpl)
				logx.Must(getByIdLogicTmpl.Execute(getByIdLogic, map[string]any{
					"setLogic":           strings.Replace(setLogic, "req.", "data.", -1),
					"modelName":          ctx.ModelName,
					"modelNameLowerCase": strings.ToLower(ctx.ModelName),
					"projectPackage":     projectCtx.Path,
					"rpcPackage":         ctx.GrpcPackage,
					"rpcName":            ctx.RpcName,
					"rpcPbPackageName":   ctx.RPCPbPackageName,
					"useUUID":            ctx.UseUUID,
					"useI18n":            ctx.UseI18n,
					"importPrefix":       ctx.ImportPrefix,
					"optionalService":    ctx.OptionalService,
					"IdType":             ctx.IdType,
					"HasCreated":         ctx.HasCreated,
				}))

				data = append(data, &ApiLogicData{
					LogicName: fmt.Sprintf("Get%sByIdLogic", ctx.ModelName),
					LogicCode: getByIdLogic.String(),
				})
			}

		}
	}

	return data
}

func GenApiData(ctx *GenLogicByProtoContext, p *parser.Proto) (string, error) {
	infoData := strings.Builder{}
	listData := strings.Builder{}
	var data string
	var hasRoutePrefix bool
	if ctx.RoutePrefix != "" {
		hasRoutePrefix = true
	} else {
		hasRoutePrefix = false
	}

	for _, v := range p.Message {
		if strings.Contains(v.Name, ctx.ModelName) {
			if fmt.Sprintf("%sInfo", ctx.ModelName) == v.Name {
				for _, field := range v.Elements {
					field.Accept(protox.MessageVisitor{})
					if entx.IsBaseProperty(protox.ProtoField.Name) {
						continue
					}

					var structData string

					jsonTag, err := format.FileNamingFormat(ctx.JSONStyle, protox.ProtoField.Name)
					if err != nil {
						return "", err
					}

					fieldComment := parser.CamelCase(protox.ProtoField.Name)
					if protox.ProtoField.Comment != "" {
						fieldComment = strings.Trim(protox.ProtoField.Comment, " ")
					}

					structData = fmt.Sprintf("\n\n        // %s\n        %s  *%s `json:\"%s,optional\"`",
						fieldComment,
						parser.CamelCase(protox.ProtoField.Name),
						entx.ConvertProtoTypeToGoType(protox.ProtoField.Type),
						jsonTag)

					infoData.WriteString(structData)
				}
			} else if strings.HasSuffix(v.Name, "ListReq") {
				for _, field := range v.Elements {
					field.Accept(protox.MessageVisitor{})

					var structData string

					jsonTag, err := format.FileNamingFormat(ctx.JSONStyle, protox.ProtoField.Name)
					if err != nil {
						return "", err
					}

					pointerStr := ""
					optionalStr := ""
					if protox.ProtoField.Optional {
						pointerStr = "*"
						optionalStr = "optional"
					}

					structData = fmt.Sprintf("\n\n        // %s\n        %s  %s%s `json:\"%s,%s\"`",
						parser.CamelCase(protox.ProtoField.Name),
						parser.CamelCase(protox.ProtoField.Name),
						pointerStr,
						entx.ConvertProtoTypeToGoType(protox.ProtoField.Type),
						jsonTag,
						optionalStr)

					if protox.ProtoField.Type == "string" {
						listData.WriteString(structData)
					}
				}
			}
		}
	}

	apiTemplateData := bytes.NewBufferString("")
	apiTmpl, _ := template.New("apiTpl").Parse(apiTpl)
	logx.Must(apiTmpl.Execute(apiTemplateData, map[string]any{
		"infoData":           infoData.String(),
		"modelName":          ctx.ModelName,
		"modelNameSpace":     strings.Replace(strcase.ToSnake(ctx.ModelName), "_", " ", -1),
		"modelNameLowerCase": strings.ToLower(ctx.ModelName),
		"modelNameSnake":     strcase.ToSnake(ctx.ModelName),
		"listData":           listData.String(),
		"apiServiceName":     ctx.APIServiceName,
		"useUUID":            ctx.UseUUID,
		"hasRoutePrefix":     hasRoutePrefix,
		"routePrefix":        ctx.RoutePrefix,
		"IdType":             ctx.IdType,
		"IdTypeLower":        strings.ToLower(ctx.IdType),
		"HasCreated":         ctx.HasCreated,
	}))
	data = apiTemplateData.String()

	return data, nil
}

func genSetLogic(v *proto.Message, ctx *GenLogicByProtoContext) string {
	var setLogic strings.Builder
	hasCreated, hasUpdated := false, false
	for _, field := range v.Elements {
		field.Accept(protox.MessageVisitor{})
		if entx.IsBaseProperty(protox.ProtoField.Name) {
			if protox.ProtoField.Name == "id" && protox.ProtoField.Type == "string" {
				ctx.UseUUID = true
			} else if protox.ProtoField.Name == "id" {
				ctx.IdType = entx.ConvertIdTypeToBaseMessage(protox.ProtoField.Type)
			}

			if protox.ProtoField.Name == "created_at" {
				hasCreated = true
			} else if protox.ProtoField.Name == "updated_at" {
				hasUpdated = true
			}
			continue
		}

		setLogic.WriteString(fmt.Sprintf("\n        \t%s: req.%s,", parser.CamelCase(protox.ProtoField.Name),
			parser.CamelCase(protox.ProtoField.Name)))
	}

	if hasUpdated && hasCreated {
		ctx.HasCreated = true
	}
	return setLogic.String()
}
