package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// OTP holds the schema definition for the OTP entity.
type OTP struct {
	ent.Schema
}

// Fields of the OTP.
func (OTP) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id"),
		field.String("phone"),
		field.String("code"),
		field.Time("expires_at"),
		field.Bool("is_used").Default(false),
		field.Time("created_at"),
		field.Time("updated_at"),
	}
}

// Edges of the OTP.
func (OTP) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("otps").Unique(),
	}
}
