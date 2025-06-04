package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// BlacklistedToken holds the schema definition for the BlacklistedToken entity.
type BlacklistedToken struct {
	ent.Schema
}

// Fields of the BlacklistedToken.
func (BlacklistedToken) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id"),
		field.String("token").Unique(),
		field.Int("user_id"),
		field.Time("blacklisted_at").Default(time.Now()),
		field.Time("expires_at"),
	}
}

// Edges of the BlacklistedToken.
func (BlacklistedToken) Edges() []ent.Edge {
	return nil
}

// Indexes of the BlacklistedToken.
func (BlacklistedToken) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("token"),
		index.Fields("user_id"),
	}
}
