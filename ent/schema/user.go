package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("id"),
		field.String("username").Unique(),
		field.String("password"),
		field.String("email").Unique(),
		field.String("name"),
		field.Time("date_of_birth").Optional(),
		field.String("gender").Optional(),
		field.String("phone_number").Optional(),
		field.String("current_address").Optional(),
		field.String("permanent_address").Optional(),
		field.Enum("role").Values("business_partner", "superadmin").Default("business_partner"),
		field.Bool("is_active").Default(true),
		field.Bool("is_email_verified").Default(false),
		field.Bool("is_verified").Default(false),
		field.Time("last_login_time").Optional(),
		field.Int("parent_id").Optional(),
		field.Time("deleted_at").Optional().Nillable(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		// One user can update many blogs
		edge.To("updated_blogs", Blogs.Type),
		// User who created this user
		edge.From("created_by_user", User.Type).
			Ref("created_users").
			Unique(),
		// Users created by this user
		edge.To("created_users", User.Type),
		// User who updated this user
		edge.From("updated_by_user", User.Type).
			Ref("updated_users").
			Unique(),
		// Users updated by this user
		edge.To("updated_users", User.Type),
		// Properties created by this user
		edge.To("created_properties", Property.Type),
	}
}

// Indexes of the User.
func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("username"),
		index.Fields("email"),
		index.Fields("phone_number"),
		index.Fields("parent_id"),
		index.Fields("is_active"),
		index.Fields("is_verified"),
		index.Fields("last_login_time"),
	}
}
