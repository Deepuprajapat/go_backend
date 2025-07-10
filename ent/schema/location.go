package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Location struct {
	ent.Schema
}

func (Location) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique(),
		field.String("locality_name").Optional(),
		field.String("city").Optional(),
		field.String("state").Optional(),
		field.String("phone_number").Optional(),
		field.String("country").Default("India"),
		field.String("pincode").Optional(),
		field.Bool("is_active").Default(true),
		field.Time("deleted_at").Optional().Nillable(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (Location) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("projects", Project.Type),
	}
}
