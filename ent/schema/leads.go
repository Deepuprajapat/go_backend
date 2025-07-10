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
			NotEmpty(),
		field.String("frequency").
			Optional(),
		field.String("name").
			NotEmpty(),
		field.String("otp").
			Optional(),
		field.String("phone").
			NotEmpty(),
		field.String("project_name").
			Optional(),
		field.String("source").
			Optional(),
		field.Text("message").
			Optional(),
		field.String("user_type").
			Optional(),
		field.Time("deleted_at").Optional().Nillable(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (Leads) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("property", Property.Type).
			Unique().
			Required(),
	}
}
