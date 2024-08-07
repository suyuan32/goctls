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
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// ConvertValidateTagToSwagger converts the validator tag to swagger comments.
func ConvertValidateTagToSwagger(tagData string) ([]string, error) {
	if tagData == "" || !strings.Contains(tagData, "validate") {
		return nil, nil
	}

	validateData := ExtractValidateString(tagData)

	return ConvertTagToComment(validateData)
}

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

// ConvertTagToComment converts validator tag to comments.
func ConvertTagToComment(tagString string) ([]string, error) {
	var result []string
	vals := strings.Split(tagString, ",")
	for _, v := range vals {
		if strings.Contains(v, "required") {
			result = append(result, "// required : true\n")
		}

		if strings.Contains(v, "min") || strings.Contains(v, "max") {
			result = append(result, fmt.Sprintf("// %s\n", strings.Replace(v, "=", " length : ", -1)))
		}

		if strings.Contains(v, "len") {
			tagSplit := strings.Split(v, "=")
			_, tagNum := tagSplit[0], tagSplit[1]
			result = append(result, fmt.Sprintf("// max length : %s\n", tagNum))
			result = append(result, fmt.Sprintf("// min length : %s\n", tagNum))
		}

		if strings.Contains(v, "gt") || strings.Contains(v, "gte") ||
			strings.Contains(v, "lt") || strings.Contains(v, "lte") {
			tagSplit := strings.Split(v, "=")
			tag, tagNum := tagSplit[0], tagSplit[1]
			if strings.Contains(tagNum, ".") {
				bitSize := len(tagNum) - strings.Index(tagNum, ".") - 1
				n, err := strconv.ParseFloat(tagNum, bitSize)
				if err != nil {
					return nil, errors.New("failed to convert the number in validate tag")
				}

				switch tag {
				case "gte":
					result = append(result, fmt.Sprintf("// min : %.*f\n", bitSize, n))
				case "gt":
					result = append(result, fmt.Sprintf("// min : %.*f\n", bitSize, n+1/math.Pow(10, float64(bitSize))))
				case "lte":
					result = append(result, fmt.Sprintf("// max : %.*f\n", bitSize, n))
				case "lt":
					result = append(result, fmt.Sprintf("// max : %.*f\n", bitSize, n-1/math.Pow(10, float64(bitSize))))
				}
			} else {
				n, err := strconv.Atoi(tagNum)
				if err != nil {
					return nil, errors.New("failed to convert the number in validate tag")
				}

				switch tag {
				case "gte":
					result = append(result, fmt.Sprintf("// min : %d\n", n))
				case "gt":
					result = append(result, fmt.Sprintf("// min : %d\n", n))
				case "lte":
					result = append(result, fmt.Sprintf("// max : %d\n", n))
				case "lt":
					result = append(result, fmt.Sprintf("// max : %d\n", n))
				}
			}

		}
	}
	return result, nil
}

// HasCustomValidation returns true if the comment has validations.
func HasCustomValidation(data string) bool {
	lowerCase := strings.ToLower(data)
	if strings.Contains(lowerCase, "max") || strings.Contains(lowerCase, "min") ||
		strings.Contains(lowerCase, "required") {
		return true
	}
	return false
}
