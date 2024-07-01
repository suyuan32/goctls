package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
)

// {{.ModelName}} holds the schema definition for the {{.ModelName}} entity.
type {{.ModelName}} struct {
	ent.Schema
}

// Fields of the {{.ModelName}}.
func ({{.ModelName}}) Fields() []ent.Field {
	return []ent.Field{}
}

// Edges of the {{.ModelName}}.
func ({{.ModelName}}) Edges() []ent.Edge {
	return nil
}

// Mixin of the {{.ModelName}}.
func ({{.ModelName}}) Mixin() []ent.Mixin {
    return []ent.Mixin{}
}

// Indexes of the {{.ModelName}}.
func ({{.ModelName}}) Indexes() []ent.Index {
    return nil
}

// Annotations of the {{.ModelName}}
func ({{.ModelName}}) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
		entsql.Annotation{Table: "{{.ModelNameLowercase}}"},
	}
}