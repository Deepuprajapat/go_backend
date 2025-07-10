package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
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
		field.String("first_name"),
		field.String("last_name"),
		field.Time("date_of_birth").Optional(),
		field.String("gender").Optional(),
		field.String("phone_number").Optional(),
		field.String("current_address").Optional(),
		field.String("permanent_address").Optional(),
		field.Bool("is_active").Default(true),
		field.Bool("is_deleted").Default(false),
		field.Bool("is_email_verified").Default(false),
		field.Bool("is_verified").Default(false),
		field.Time("last_login_time").Optional(),
		field.Int("parent_id").Optional(),
		field.String("photo_url").Optional(),
		field.Time("created_at").Default(time.Now()),
		field.Time("updated_at").Default(time.Now()).UpdateDefault(time.Now),
		field.Int("created_by").Optional(),
		field.Int("updated_by").Optional(),
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
	}
}
