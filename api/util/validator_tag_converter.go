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

package util

import (
	"strings"
)

// ExtractValidateString extracts the validator's string.
func ExtractValidateString(data string) string {
	beginIndex := strings.Index(data, "validate")
	if beginIndex == -1 {
		return ""
	}
	firstQuotationMark := 0
	for i := beginIndex; i < len(data); i++ {
		if data[i] == '"' && firstQuotationMark == 0 {
			firstQuotationMark = i
		} else if data[i] == '"' && firstQuotationMark != 0 {
			return data[firstQuotationMark+1 : i]
		}
	}
	return ""
}
