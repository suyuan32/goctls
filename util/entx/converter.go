// Copyright 2023 The Ryan SU Authors (https://github.com/suyuan32). All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package entx

import (
	"regexp"
	"strings"

	"github.com/suyuan32/goctls/rpc/parser"
)

// ConvertEntTypeToProtoType returns prototype from ent type
func ConvertEntTypeToProtoType(typeName string) string {
	switch typeName {
	case "float32":
		typeName = "float"
	case "float64":
		typeName = "double"
	case "float":
		typeName = "double"
	case "int":
		typeName = "int64"
	case "uint":
		typeName = "uint64"
	case "[16]byte":
		typeName = "string"
	case "uint8", "uint16":
		typeName = "uint32"
	case "int8", "int16":
		typeName = "int32"
	}
	return typeName
}

// ConvertProtoTypeToGoType returns go type from proto type
func ConvertProtoTypeToGoType(typeName string) string {
	switch typeName {
	case "float":
		typeName = "float32"
	case "double":
		typeName = "float64"
	}
	return typeName
}

// ConvertSpecificNounToUpper is used to convert snack format to Ent format
func ConvertSpecificNounToUpper(str string) string {
	target := parser.CamelCase(str)

	data := []struct {
		Origin string
		Target string
	}{
		{
			"Uuid",
			"UUID",
		},
		{
			"Api",
			"API",
		},
		{
			"Uri",
			"URI",
		},
		{
			"Url",
			"URL",
		},
		{
			"Ip",
			"IP",
		},
	}

	for _, v := range data {
		target = strings.Replace(target, v.Origin, v.Target, -1)
	}

	target = ConvertIdFieldToUpper(target)

	return target
}

// ConvertIdFieldToUpper is used to convert snack format Id to Ent format
func ConvertIdFieldToUpper(target string) string {
	if IsNotIDField(target) {
		if strings.Contains(target, "Id") {
			target = strings.Replace(target, "Id", "ID", -1)
		}
	} else {
		if strings.HasSuffix(target, "Id") {
			target = target[:len(target)-1] + "D"
		}
	}
	return target
}

// IsNotIDField Judge whether the field is not an ID field
func IsNotIDField(field string) bool {
	compile, err := regexp.Compile("Id[a-z]+/g")
	if err != nil {
		return false
	}

	return compile.MatchString(field)
}

// ConvertEntTypeToGotype returns go type from ent type
func ConvertEntTypeToGotype(prop string) string {
	switch prop {
	case "int":
		return "int64"
	case "uint":
		return "uint64"
	case "uint8", "uint16":
		return "uint32"
	case "int8", "int16":
		return "int32"
	}
	return prop
}

// ConvertEntTypeToGotypeInSingleApi returns go type from ent type in single API service
func ConvertEntTypeToGotypeInSingleApi(prop string) string {
	switch prop {
	case "[16]byte":
		return "string"
	case "time.Time":
		return "int64"
	default:
		return prop
	}
}

// ConvertIDType returns uuid type by uuid flag
func ConvertIDType(useUUID bool, t string) string {
	if useUUID {
		return "string"
	} else {
		switch t {
		case "int32", "int64", "uint32", "uint64", "string":
			return t
		default:
			return "uint64"
		}
	}
}

// ConvertOnlyEntTypeToGoType converts the type that only ent has to go type.
func ConvertOnlyEntTypeToGoType(t string) string {
	switch t {
	case "int8", "int16":
		return "int32"
	case "uint8", "uint16":
		return "uint32"
	default:
		return "uint32"
	}
}

// ConvertIdTypeToBaseMessage returns base message name when id type is not uint64 or string.
func ConvertIdTypeToBaseMessage(t string) string {
	if t == "uint64" || t == "[16]byte" {
		return ""
	} else {
		switch t {
		case "int32":
			return "Int32"
		case "int64":
			return "Int64"
		case "uint32":
			return "Uint32"
		case "string":
			return "String"
		default:
			return ""
		}
	}
}
