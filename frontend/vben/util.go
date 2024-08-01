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

package vben

import "strings"

func ConvertGoTypeToTsType(goType string) string {
	switch goType {
	case "int", "uint", "int8", "uint8", "int16", "uint16", "int32", "uint32", "uint64", "int64", "float", "float32", "float64",
		"*int", "*uint", "*int8", "*uint8", "*int16", "*uint16", "*int32", "*uint32", "*uint64", "*int64", "*float", "*float32", "*float64":
		goType = "number"
	case "[]int", "[]uint", "[]int32", "[]int64", "[]uint32", "[]uint64", "[]float", "[]float32", "[]float64":
		goType = "number[]"
	case "string", "*string":
		goType = "string"
	case "[]string":
		goType = "string[]"
	case "bool", "*bool":
		goType = "boolean"

	}
	return goType
}

func FindBeginEndOfLocaleField(data, target string) (int, int) {

	begin := strings.Index(data, target)

	if begin == -1 {
		return -1, -1
	}

	var end int

	for i := begin; i < len(data); i++ {
		if data[i] == '}' {
			end = i + 2
			break
		}
	}

	return begin, end
}
