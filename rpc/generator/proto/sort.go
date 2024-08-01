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
	"github.com/duke-git/lancet/v2/slice"
	"github.com/suyuan32/goctls/rpc/parser"
)

func SortImport(data map[string]parser.Import) (result []string) {
	result = []string{}
	for k, _ := range data {
		result = append(result, k)
	}

	slice.Sort(result)

	return result
}

func SortEnum(data map[string]parser.Enum) (result []string) {
	result = []string{}
	for k, _ := range data {
		result = append(result, k)
	}

	slice.Sort(result)

	return result
}

func SortMessage(data map[string]parser.Message) (result []string) {
	result = []string{}
	for k, _ := range data {
		result = append(result, k)
	}

	slice.Sort(result)

	return result
}

func SortService(data map[string]parser.Service) (result []string) {
	result = []string{}
	for k, _ := range data {
		result = append(result, k)
	}

	slice.Sort(result)

	return result
}
