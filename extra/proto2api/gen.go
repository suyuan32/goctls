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

package proto2api

import (
	"fmt"
	"path/filepath"
	"strings"

	set "github.com/duke-git/lancet/v2/datastructure/set"
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/gookit/color"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/suyuan32/goctls/rpc/parser"
	"github.com/suyuan32/goctls/util/format"
	"github.com/suyuan32/goctls/util/protox"
)

var (
	VarStringProtoPath string
	VarStringModelName string
	VarStringApiPath   string
	VarBoolMultiple    bool
	VarStringGroupName string
	VarStringJsonStyle string
)

type GenContext struct {
	ProtoPath string
	ModelName string
	ApiPath   string
	Multiple  bool
	GroupName string
	JsonStyle string
}

func Gen(_ *cobra.Command, _ []string) error {
	if VarStringGroupName == "" {
		VarStringGroupName = strings.ToLower(VarStringModelName)
	}
	return DoGen(&GenContext{
		ProtoPath: VarStringProtoPath,
		ModelName: VarStringModelName,
		ApiPath:   VarStringApiPath,
		Multiple:  VarBoolMultiple,
		GroupName: VarStringGroupName,
		JsonStyle: VarStringJsonStyle,
	})
}

func DoGen(g *GenContext) error {
	color.Green.Println("Generating...")

	protoFilePath, err := filepath.Abs(g.ProtoPath)
	if err != nil {
		return errors.Wrap(err, "failed to load proto file")
	}

	var apiFilePath, apiData string
	if g.ApiPath != "" {
		apiFilePath, err = filepath.Abs(g.ApiPath)
		if err != nil {
			return errors.Wrap(err, "failed to load api file")
		}

		apiData, err = fileutil.ReadFileToString(apiFilePath)
		if err != nil {
			return errors.Wrap(err, "failed to load api data")
		}
	}

	protoParser := parser.NewDefaultProtoParser()

	protoData, err := protoParser.Parse(protoFilePath, g.Multiple)
	if err != nil {
		return errors.Wrap(err, "failed to parse proto data")
	}

	var typeData, routeData strings.Builder
	typeNameSet := set.New[string]()

	// gen route
	for _, v := range protoData.Service {
		for _, r := range v.RPC {
			if strings.Contains(r.Comment.Message(), g.GroupName) {
				if SkipRpcName(r.Name, g.ModelName) {
					continue
				}
				urlName, err := format.FileNamingFormat("go_zero", r.Name)
				if err != nil {
					return err
				}
				if strings.Contains(apiData, r.Name) {
					continue
				}

				typeNameSet.Add(r.RequestType, r.ReturnsType)

				routeData.WriteString(fmt.Sprintf("\n    // %s\n    @handler %s\n    post /%s/%s (%s) returns (%s)\n",
					r.Name, r.Name, strings.ToLower(g.ModelName), urlName, r.RequestType, ConvertRpcBaseStruct(r.ReturnsType)))
			}

		}
	}

	if routeData.Len() < 2 {
		return errors.New("rpc not found, please check your group name is correct")
	}

	// gen type
	protox.ProtoField = &protox.ProtoFieldData{}
	for _, v := range protoData.Message {
		if typeNameSet.Contain(v.Name) && !strings.Contains(v.Name, "Base") && !strings.Contains(apiData, v.Name) {
			typeData.WriteString(fmt.Sprintf("\n    // %s \n    %s {\n", v.Name, v.Name))
			for _, t := range v.Elements {
				t.Accept(protox.MessageVisitor{})

				if SkipBaseMessage(protox.ProtoField.Name) || (!(protox.ProtoField.Name[0] > 'a' && (protox.ProtoField.Name[0] < 'z') &&
					strings.Contains(apiData, protox.ProtoField.Name))) {
					continue
				}

				optionalStr, repeatStr, pointerStr := "", "", ""
				if protox.ProtoField.Optional {
					optionalStr = ",optional"
					pointerStr = "*"
				}

				if protox.ProtoField.Repeated {
					repeatStr = "[]"
				}

				jsonFieldStr, err := format.FileNamingFormat(g.JsonStyle, protox.ProtoField.Name)
				if err != nil {
					return errors.Wrap(err, "failed to convert json field style")
				}

				fieldStr, err := format.FileNamingFormat("GoZero", protox.ProtoField.Name)
				if err != nil {
					return errors.Wrap(err, "failed to convert json field style")
				}

				typeData.WriteString(fmt.Sprintf("        // %s\n        %s  %s%s%s `json:\"%s%s\"`\n\n",
					fieldStr,
					fieldStr,
					pointerStr,
					repeatStr,
					protox.ProtoField.Type,
					jsonFieldStr,
					optionalStr,
				))
			}
			typeData.WriteString("    }\n")
		}
	}

	if g.ApiPath == "" {
		fmt.Printf("\n\nTYPE DATA  \n\n%s\n\n", typeData.String())
		fmt.Printf("ROUTE DATA  \n\n%s\n\n", routeData.String())
	} else {
		typeIndex := FindTypeContentIndex(apiData)
		if typeIndex != -1 {
			apiData = apiData[:typeIndex+1] + typeData.String() + apiData[typeIndex+1:]
		}

		serviceIndex := FindServiceContentIndex(apiData)
		if serviceIndex != -1 {
			apiData = apiData[:serviceIndex+1] + routeData.String() + apiData[serviceIndex+1:]
		}

		err = fileutil.WriteStringToFile(apiFilePath, apiData, false)
		if err != nil {
			return errors.Wrap(err, "failed to write data to api file")
		}
	}

	color.Green.Println("Generate successfully")

	return err
}
