package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Location struct {
	Base
}

func (Location) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique(),
		field.String("locality_name"),
		field.String("city"),
		field.String("state"),
		field.String("phone_number"),
		field.String("country").Default("India"),
		field.String("pincode"),
		field.String("area_type"), // Sector, Phase, Block, etc.
		field.Bool("is_active").Default(true),
		field.String("slug"), // URL-friendly version of locality name
		field.Time("created_at").Default(time.Now()),
		field.Time("updated_at"),
	}
}

func (Location) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("projects", Project.Type),
	}
}
