package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type Base struct {
	ent.Schema
}

func (Base) Fields() []ent.Field {
	return []ent.Field{
		field.Bool("deleted_at").Optional().Nillable(),
		field.Time("created_at").Default(time.Now()),
		field.Time("updated_at"),
	}
}
