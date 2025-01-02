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

package vben5

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"

	"github.com/suyuan32/goctls/api/spec"
	util2 "github.com/suyuan32/goctls/api/util"
	"github.com/suyuan32/goctls/util"
)

func genData(g *GenContext) error {
	var basicData, searchFormData, formData strings.Builder
	var useBaseInfo bool
	var statusBasicColumnData, statusFormColumnData string
	var stateBasicColumnData, stateFormColumnData string

	// generate basic and search form data
	for _, v := range g.ApiSpec.Types {
		if v.Name() == fmt.Sprintf("%sInfo", g.ModelName) {
			specData, ok := v.(spec.DefineStruct)
			if !ok {
				return errors.New("cannot get the field")
			}

			for _, val := range specData.Members {
				if val.Name == "" {
					tmpType, _ := val.Type.(spec.DefineStruct)
					if tmpType.Name() == "BaseIDInfo" || tmpType.Name() == "BaseUUIDInfo" {
						useBaseInfo = true
					}
				} else if val.Name == "Status" {
					statusRenderData := bytes.NewBufferString("")
					protoTmpl, _ := template.New("proto").Parse(statusRenderTpl)
					_ = protoTmpl.Execute(statusRenderData, map[string]any{
						"modelName": strings.TrimSuffix(specData.RawName, "Info"),
					})
					statusBasicColumnData = statusRenderData.String()

					statusFormColumnData = fmt.Sprintf("\n  {\n    fieldName: '%s',\n    label: $t('%s'),\n    component: 'RadioButtonGroup',\n"+
						"    defaultValue: 1,\n    componentProps: {\n      options: [\n        { label: $t('common.on'), value: 1 },\n    "+
						"    { label: $t('common.off'), value: 2 },\n      ],\n    },\n  },",
						GetJsonTagName(val.Tags()),
						fmt.Sprintf("%s.%s.%s", g.FolderName,
							strcase.ToLowerCamel(strings.TrimSuffix(specData.RawName, "Info")),
							strcase.ToLowerCamel(val.Name)),
					)
				} else if val.Name == "State" {
					stateRenderData := bytes.NewBufferString("")
					protoTmpl, _ := template.New("proto").Parse(stateRenderTpl)
					_ = protoTmpl.Execute(stateRenderData, map[string]any{
						"modelName": strings.TrimSuffix(specData.RawName, "Info"),
					})
					stateBasicColumnData = stateRenderData.String()

					stateFormColumnData = fmt.Sprintf("\n  {\n    fieldName: '%s',\n    label: $t('%s'),\n    component: 'RadioButtonGroup',\n"+
						"    defaultValue: true,\n    componentProps: {\n      options: [\n        { label: $t('common.on'), value: true },\n    "+
						"    { label: $t('common.off'), value: false },\n      ],\n    },\n  },",
						GetJsonTagName(val.Tags()),
						fmt.Sprintf("%s.%s.%s", g.FolderName,
							strcase.ToLowerCamel(strings.TrimSuffix(specData.RawName, "Info")),
							strcase.ToLowerCamel(val.Name)),
					)
				} else {
					basicData.WriteString(fmt.Sprintf("\n  {\n    title: $t('%s'),\n    field: '%s',\n  },",
						fmt.Sprintf("%s.%s.%s", g.FolderName,
							strcase.ToLowerCamel(strings.TrimSuffix(specData.RawName, "Info")),
							strcase.ToLowerCamel(val.Name)), GetJsonTagName(val.Tags())))

					formData.WriteString(fmt.Sprintf("\n  {\n    fieldName: '%s',\n    label: $t('%s'),\n    %s\n%s  },",
						GetJsonTagName(val.Tags()),
						fmt.Sprintf("%s.%s.%s", g.FolderName,
							strcase.ToLowerCamel(strings.TrimSuffix(specData.RawName, "Info")),
							strcase.ToLowerCamel(val.Name)),
						getComponent(val.Type.Name()),
						GetRules(val, false),
					))
				}
			}

			// put here in order to put status in the end
			if g.HasStatus {
				basicData.WriteString(statusBasicColumnData)
				formData.WriteString(statusFormColumnData)
			}

			// put here in order to put state in the end
			if g.HasState {
				basicData.WriteString(stateBasicColumnData)
				formData.WriteString(stateFormColumnData)
			}
		}

		if v.Name() == fmt.Sprintf("%sListReq", g.ModelName) {
			specData, ok := v.(spec.DefineStruct)
			if !ok {
				return errors.New("cannot get field")
			}

			for _, val := range specData.Members {
				if val.Name != "" {
					searchFormData.WriteString(fmt.Sprintf("\n  {\n    fieldName: '%s',\n    label: $t('%s'),\n    %s\n   %s  },",
						strcase.ToLowerCamel(val.Name),
						fmt.Sprintf("%s.%s.%s", g.FolderName,
							strcase.ToLowerCamel(strings.TrimSuffix(specData.RawName, "ListReq")),
							strcase.ToLowerCamel(val.Name)),
						getComponent(val.Type.Name()),
						GetRules(val, true),
					))
				}
			}
		}
	}

	if err := util.With("dataTpl").Parse(dataTpl).SaveTo(map[string]any{
		"modelName":           g.ModelName,
		"modelNameLowerCamel": strcase.ToLowerCamel(g.ModelName),
		"folderName":          g.FolderName,
		"basicData":           basicData.String(),
		"searchFormData":      searchFormData.String(),
		"formData":            formData.String(),
		"useBaseInfo":         useBaseInfo,
		"useUUID":             g.UseUUID,
		"hasStatus":           g.HasStatus || g.HasState,
	},
		filepath.Join(g.ViewDir, "schemas.ts"), g.Overwrite); err != nil {
		return err
	}
	return nil
}

func getComponent(dataType string) string {
	switch dataType {
	case "string", "*string":
		return "component: 'Input',"
	case "int", "uint", "int8", "uint8", "int16", "uint16", "int32", "int64", "uint32", "uint64", "float32", "float64",
		"*int", "*uint", "*int8", "*uint8", "*int16", "*uint16", "*int32", "*int64", "*uint32", "*uint64", "*float32", "*float64":
		return "component: 'InputNumber',"
	case "bool", "*bool":
		return "component: 'RadioButtonGroup',\n    defaultValue: true,\n    componentProps: {\n      options: [\n        { label: $t('common.on'), value: true },\n        { label: $t('common.off'), value: false },\n      ],\n    },"
	default:
		return "component: 'Input',"
	}
}

// GetRules returns the rules from tag.
func GetRules(t spec.Member, optional bool) string {
	validatorString := util2.ExtractValidateString(t.Tag)
	if validatorString == "" {
		return ""
	}

	rules, err := ConvertTagToRules(validatorString)
	if err != nil {
		return ""
	}

	optionalStr := ""
	if optional {
		optionalStr = ".optional()"
	}

	switch GetRuleType(t.Type.Name()) {
	case "string":
		return fmt.Sprintf("    rules: z.string()%s%s,\n", rules, optionalStr)
	case "number":
		return fmt.Sprintf("    rules: z.number()%s%s,\n", rules, optionalStr)
	case "float":
		return fmt.Sprintf("    rules: z.number()%s%s,\n", rules, optionalStr)
	default:
		return ""
	}
}

// GetRuleType returns the rule type from go type.
func GetRuleType(t string) string {
	switch t {
	case "string", "*string":
		return "string"
	case "int", "uint", "int8", "uint8", "int16", "uint16", "int32", "int64", "uint32", "uint64",
		"*int", "*uint", "*int8", "*uint8", "*int16", "*uint16", "*int32", "*int64", "*uint32", "*uint64":
		return "number"
	case "float32", "float64", "*float32", "*float64":
		return "float"
	default:
		return "string"
	}
}

// ConvertTagToRules converts validator tag to rules.
func ConvertTagToRules(tagString string) (string, error) {
	vals := strings.Split(tagString, ",")
	resultStr := strings.Builder{}
	for _, v := range vals {
		if strings.Contains(v, "min") || strings.Contains(v, "max") {
			tmp := strings.Split(v, "=")
			if len(tmp) == 2 {
				resultStr.WriteString(fmt.Sprintf(".%s(%s)", strings.ToLower(tmp[0]),
					strings.ToLower(tmp[1])))
			}
		}

		if strings.Contains(v, "len") {
			tmp := strings.Split(v, "=")
			if len(tmp) == 2 {
				resultStr.WriteString(fmt.Sprintf(".length(%s)",
					strings.ToLower(tmp[1])))
			}
		}

		if strings.Contains(v, "gt") || strings.Contains(v, "gte") ||
			strings.Contains(v, "lt") || strings.Contains(v, "lte") {
			tagSplit := strings.Split(v, "=")
			tag, tagNum := tagSplit[0], tagSplit[1]
			if strings.Contains(tagNum, ".") {
				bitSize := len(tagNum) - strings.Index(tagNum, ".") - 1
				n, err := strconv.ParseFloat(tagNum, bitSize)
				if err != nil {
					return "", errors.New("failed to convert the number in validate tag")
				}

				switch tag {
				case "gte":
					resultStr.WriteString(fmt.Sprintf(".min(%.*f)", bitSize, n))
				case "gt":
					resultStr.WriteString(fmt.Sprintf(".min(%.*f)", bitSize, n+1/math.Pow(10, float64(bitSize))))
				case "lte":
					resultStr.WriteString(fmt.Sprintf(".max(%.*f)", bitSize, n))
				case "lt":
					resultStr.WriteString(fmt.Sprintf(".max(%.*f)", bitSize, n-1/math.Pow(10, float64(bitSize))))
				}
			} else {
				n, err := strconv.Atoi(tagNum)
				if err != nil {
					return "", errors.New("failed to convert the number in validate tag")
				}

				switch tag {
				case "gte":
					resultStr.WriteString(fmt.Sprintf(".min(%d)", n))
				case "gt":
					resultStr.WriteString(fmt.Sprintf(".min(%d)", n))
				case "lte":
					resultStr.WriteString(fmt.Sprintf(".max(%d)", n))
				case "lt":
					resultStr.WriteString(fmt.Sprintf(".max(%d)", n))
				}
			}

		}
	}
	return resultStr.String(), nil
}
