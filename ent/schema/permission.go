package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Permission holds the schema definition for the Permission entity.
type Permission struct {
	ent.Schema
}

// Fields of the Permission.
func (Permission) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id"),
		field.String("name").Unique(),
		field.String("description").Optional(),
		field.Time("created_at"),
		field.Time("updated_at"),
	}
}

// Edges of the Permission.
func (Permission) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("users", User.Type).Ref("permissions"),
	}
}
