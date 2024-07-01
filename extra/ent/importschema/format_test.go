package importschema

import "testing"

func Test_getSchemaName(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{data: `
package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/suyuan32/simple-admin-common/orm/ent/mixins"

	mixins2 "github.com/suyuan32/simple-admin-core/rpc/ent/schema/mixins"
)

type User struct {
	ent.Schema
}`},
			want: "User",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSchemaName(tt.args.data); got != tt.want {
				t.Errorf("getSchemaName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_removeOldAnnotation(t *testing.T) {
	type args struct {
		data       string
		schemaName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test_removeOldAnnotation",
			args: args{data: `
package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type SysRole struct {
	ent.Schema
}

func (SysRole) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id"),
		field.Time("created_at").Comment("Create Time | 创建日期"),
	}
}
func (SysRole) Edges() []ent.Edge {
	return nil
}
func (SysRole) Annotations() []schema.Annotation {
	return nil
}

func (SysRole) Indexes() []ent.Index {
	return nil
}
`, schemaName: "SysRole"},
			want: `
package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type SysRole struct {
	ent.Schema
}

func (SysRole) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id"),
		field.Time("created_at").Comment("Create Time | 创建日期"),
	}
}
func (SysRole) Edges() []ent.Edge {
	return nil
}


func (SysRole) Indexes() []ent.Index {
	return nil
}
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeOldAnnotation(tt.args.data, tt.args.schemaName); got != tt.want {
				t.Errorf("removeOldAnnotation() = %v, want %v", got, tt.want)
			}
		})
	}
}
