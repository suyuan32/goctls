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
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gookit/color"

	"github.com/suyuan32/goctls/rpc/parser"
	"github.com/suyuan32/goctls/util/protox"
)

type ProtoContext struct {
	ProtoDir   string
	OutputPath string
	Multiple   bool
}

func MergeProto(p *ProtoContext) error {
	var protoFiles []string
	err := filepath.WalkDir(p.ProtoDir, func(path string, d fs.DirEntry, err error) error {
		if strings.HasSuffix(path, "proto") {
			protoFiles = append(protoFiles, path)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to merge proto : %v", err)
	}

	if len(protoFiles) == 0 {
		return errors.New("proto files not found")
	}

	var protoFileData []parser.Proto
	protoParser := parser.NewDefaultProtoParser()

	for _, v := range protoFiles {
		data, err := protoParser.Parse(v, p.Multiple)
		if err != nil {
			return fmt.Errorf("fail to parse proto file: %v, path: %s", err, v)
		}
		protoFileData = append(protoFileData, data)
	}

	protoString := genProtoString(parseProto(protoFileData))

	err = os.WriteFile(p.OutputPath, []byte(protoString), os.ModePerm)
	if err != nil {
		return err
	}

	color.Green.Println("Merge proto files successfully")

	return err
}

func parseProto(data []parser.Proto) parser.Proto {
	var result parser.Proto
	importSet := map[string]parser.Import{}
	enumSet := map[string]parser.Enum{}
	messageSet := map[string]parser.Message{}
	serviceSet := map[string]parser.Service{}

	for _, v := range data {
		if v.Name != "" {
			result.Name = v.Name
		}

		if v.PbPackage != "" {
			result.PbPackage = v.PbPackage
		}

		if v.GoPackage != "" {
			result.GoPackage = v.GoPackage
		}

		if v.Package.Package != nil {
			if v.Package.Name != "" {
				result.Package = v.Package
			}
		}

		for _, i := range v.Import {
			importSet[i.Filename] = i
		}

		for _, e := range v.Enum {
			enumSet[e.Name] = e
		}

		for _, m := range v.Message {
			messageSet[m.Name] = m
		}

		for _, s := range v.Service {
			if _, ok := serviceSet[s.Name]; !ok {
				serviceSet[s.Name] = s
			} else {
				tmp := serviceSet[s.Name]
				tmp.RPC = append(tmp.RPC, s.RPC...)
				serviceSet[s.Name] = tmp
			}
		}

	}

	for _, v := range SortImport(importSet) {
		result.Import = append(result.Import, importSet[v])
	}

	for _, e := range SortEnum(enumSet) {
		result.Enum = append(result.Enum, enumSet[e])
	}

	for _, m := range SortMessage(messageSet) {
		result.Message = append(result.Message, messageSet[m])
	}

	for _, s := range SortService(serviceSet) {
		result.Service = append(result.Service, serviceSet[s])
	}

	return result
}

func genProtoString(data parser.Proto) string {
	var protoString strings.Builder
	protox.ProtoField = &protox.ProtoFieldData{}

	protoString.WriteString("syntax = \"proto3\";\n\n")
	protoString.WriteString(fmt.Sprintf("package %s;\n", data.Package.Name))
	protoString.WriteString(fmt.Sprintf("option go_package = \"%s\";\n\n", data.GoPackage))

	for _, i := range data.Import {
		protoString.WriteString(fmt.Sprintf("import \"%s\";\n\n", i.Filename))
	}

	for _, e := range data.Enum {
		if e.Comment != nil {
			protoString.WriteString(protox.GenCommentString(e.Comment.Lines, false))
		}
		if len(e.Elements) != 0 {
			protoString.WriteString(fmt.Sprintf("enum %s {\n", e.Name))
		} else {
			protoString.WriteString(fmt.Sprintf("enum %s {", e.Name))
		}
		for _, ele := range e.Elements {
			ele.Accept(protox.MessageVisitor{})
			protoString.WriteString(fmt.Sprintf("  %s = %d;\n", protox.ProtoField.Name, protox.ProtoField.Sequence))
		}
		protoString.WriteString("}\n\n")
	}

	for _, m := range data.Message {
		if m.Comment != nil {
			protoString.WriteString(protox.GenCommentString(m.Comment.Lines, false))
		}
		if len(m.Elements) != 0 {
			protoString.WriteString(fmt.Sprintf("message %s {\n", m.Name))
		} else {
			protoString.WriteString(fmt.Sprintf("message %s {", m.Name))
		}
		for _, e := range m.Elements {
			e.Accept(protox.MessageVisitor{})
			prefixStr := ""
			if protox.ProtoField.Repeated {
				prefixStr = "repeated "
			}
			if protox.ProtoField.Optional {
				prefixStr = "optional "
			}

			fieldComment := ""
			if protox.ProtoField.Comment != "" {
				fieldComment = fmt.Sprintf("  // %s\n", protox.ProtoField.Comment)
			}

			protoString.WriteString(fmt.Sprintf("%s  %s%s %s = %d;\n", fieldComment, prefixStr, protox.ProtoField.Type, protox.ProtoField.Name, protox.ProtoField.Sequence))
		}
		protoString.WriteString("}\n\n")
	}

	for _, s := range data.Service {
		if s.Comment != nil {
			if s.Comment != nil {
				protoString.WriteString(protox.GenCommentString(s.Comment.Lines, false))
			}
		}
		protoString.WriteString(fmt.Sprintf("service %s {\n", s.Name))
		for _, rpc := range s.RPC {
			if rpc.Comment != nil {
				protoString.WriteString(protox.GenCommentString(rpc.Comment.Lines, true))
			}

			var request, response string
			if !rpc.StreamsRequest {
				request = rpc.RequestType
			} else {
				request = fmt.Sprintf("stream %s", rpc.RequestType)
			}

			if !rpc.StreamsReturns {
				response = rpc.ReturnsType
			} else {
				response = fmt.Sprintf("stream %s", rpc.ReturnsType)
			}

			protoString.WriteString(fmt.Sprintf("  rpc %s(%s) returns (%s);\n", rpc.Name, request, response))
		}
		protoString.WriteString("}\n\n")
	}

	return protoString.String()
}
