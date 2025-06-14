package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Location struct {
	Base
}

func (Location) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique().Positive(),
		field.String("locality_name").Optional().Nillable(),
		field.String("city").Optional().Nillable(),
		field.String("state").Optional().Nillable(),
		field.String("phone_number").Optional().Nillable(),
		field.String("country").Default("India").NotEmpty(),
		field.String("pincode").Optional().Nillable(),
		field.Bool("is_active").Default(true),
	}
}

func (Location) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("projects", Project.Type),
	}
}
