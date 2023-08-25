package schema

import (
	"entgo.io/ent"
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
    return nil
}

// Indexes of the {{.ModelName}}.
func ({{.ModelName}}) Indexes() []ent.Index {
    return nil
}

// Annotations of the {{.ModelName}}
func ({{.ModelName}}) Annotations() []schema.Annotation {
	return nil
}