package schema

import "entgo.io/ent"

// {{.ModelName}} holds the schema definition for the {{.ModelName}} entity.
type {{.ModelName}} struct {
	ent.Schema
}

// Fields of the {{.ModelName}}.
func ({{.ModelName}}) Fields() []ent.Field {
	return nil
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