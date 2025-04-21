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

package entx

import (
	"strings"
)

// IsTimeProperty returns true when the string contains time suffix
func IsTimeProperty(prop string) bool {
	if prop == "time.Time" {
		return true
	}
	return false
}

// IsUpperProperty returns true when the string
// contains Ent upper string such as uuid, api and id
func IsUpperProperty(prop string) bool {
	prop = strings.ToLower(prop)

	data := []string{"uuid", "id", "api", "url", "uri", "ip"}

	for _, v := range data {
		if strings.Contains(prop, v) {
			return true
		}
	}

	return false
}

// IsBaseProperty returns true when prop name is
// id, created_at, updated_at, deleted_at
func IsBaseProperty(prop string) bool {
	if prop == "id" || prop == "created_at" || prop == "updated_at" || prop == "deleted_at" || prop == "tenant_id" {
		return true
	}
	return false
}

// IsGoTypeNotPrototype returns true when property type is
// prototype but not go type
func IsGoTypeNotPrototype(prop string) bool {
	if prop == "int" || prop == "uint" || prop == "[16]byte" {
		return true
	}
	return false
}

// IsUUIDType returns true when prop is Ent's UUID type
func IsUUIDType(prop string) bool {
	if prop == "[16]byte" {
		return true
	}
	return false
}

// IsOnlyEntType returns true when the type is only in ent schema. e.g. uint8
func IsOnlyEntType(t string) bool {
	switch t {
	case "int8", "uint8", "int16", "uint16":
		return true
	default:
		return false
	}
}

// IsPageProperty returns true when prop name is
// pageNo, pageSize
func IsPageProperty(prop string) bool {
	if prop == "page" || prop == "pagesize" || prop == "page_size" {
		return true
	}
	return false
}
