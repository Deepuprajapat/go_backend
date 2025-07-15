package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Leads struct {
	ent.Schema
}

func (Leads) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique(),
		field.String("email").
			Optional(),
		field.String("name").
			NotEmpty(),
		field.String("phone").
			NotEmpty(),
		field.String("otp").
			Optional(),
		field.Text("message").
			Optional(),
		field.String("source").
			Optional().
			Default("Organic"),
		field.Bool("is_duplicate").
			Optional().
			Default(false),
		field.String("duplicate_reference_id").
			Optional(),
		field.Bool("otp_verified").
			Optional().
			Default(false),
		field.Time("deleted_at").Optional().Nillable(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (Leads) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("property", Property.Type).
			Unique(),
		edge.To("project", Project.Type).
			Unique(),
	}
}
