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

package protox

import (
	"fmt"
	"strings"

	"github.com/emicklei/proto"
)

func FindBeginEndOfService(service, serviceName string) (begin, mid, end int) {
	beginIndex := strings.Index(service, serviceName)
	begin, end = -1, -1
	if beginIndex > 0 {
		for i := beginIndex; i < len(service); i++ {
			if service[i] == '}' {
				end = i
				break
			} else if service[i] == '{' {
				mid = i
			}
		}

		for i := beginIndex; i >= 0; i-- {
			if service[i] == 's' {
				begin = i
				break
			}
		}
	}
	return begin, mid, end
}

var ProtoField *ProtoFieldData

type ProtoFieldData struct {
	Name     string
	Type     string
	Repeated bool
	Optional bool
	Sequence int
	Comment  string
}

type MessageVisitor struct {
	proto.NoopVisitor
}

func (m MessageVisitor) VisitNormalField(i *proto.NormalField) {
	if i.Comment != nil {
		ProtoField.Comment = i.Comment.Message()
	} else {
		ProtoField.Comment = ""
	}
	ProtoField.Name = i.Field.Name
	ProtoField.Type = i.Field.Type
	ProtoField.Repeated = i.Repeated
	ProtoField.Optional = i.Optional
	ProtoField.Sequence = i.Sequence
}

func (m MessageVisitor) VisitMapField(i *proto.MapField) {
	ProtoField.Name = i.Field.Name
	ProtoField.Type = fmt.Sprintf("map<%s,%s>", i.KeyType, i.Field.Type)
	ProtoField.Sequence = i.Sequence
	ProtoField.Repeated = false
	ProtoField.Optional = false
}

func (m MessageVisitor) VisitEnumField(i *proto.EnumField) {
	ProtoField.Name = i.Name
	ProtoField.Type = ""
	ProtoField.Sequence = i.Integer
	ProtoField.Repeated = false
	ProtoField.Optional = false
}

func GenCommentString(comments []string, space bool) string {
	var commentsString strings.Builder
	var spaceString string
	if space {
		spaceString = "  "
	}

	for _, v := range comments {
		commentsString.WriteString(fmt.Sprintf("%s// %s\n", spaceString, v))
	}
	return commentsString.String()
}
