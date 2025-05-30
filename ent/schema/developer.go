package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Developer struct {
	ent.Schema
}

func (Developer) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique(),
		field.String("name"),
		field.String("legal_name"),
		field.String("url"),
		field.Int("established_year"),
		field.Int("project_count"),
		field.JSON("contact_info", DeveloperContactInfo{}),
		field.JSON("media_content", DeveloperMediaContent{}),
		field.JSON("seo_meta", DeveloperSEOMeta{}),
		field.Bool("is_active").Default(false),
		field.Bool("is_verified").Default(false),
		field.String("search_context"),
		field.Time("created_at").Default(time.Now()),
		field.Time("updated_at"),
	}
}

func (Developer) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("projects", Project.Type),
	}
}

// contact information
type DeveloperContactInfo struct {
	DeveloperAddress string `json:"developer_address"`
	Phone            string `json:"phone"`
}

// media and content
type DeveloperMediaContent struct {
	DeveloperLogo    string `json:"developer_logo"`
	AltDeveloperLogo string `json:"alt_developer_logo"`
	About            string `json:"about"`
	Overview         string `json:"overview"`
	Disclaimer       string `json:"disclaimer"`
}

type DeveloperSEOMeta struct {
	MetaTitle       string `json:"meta_title"`
	MetaDescription string `json:"meta_description"`
	MetaKeywords    string `json:"meta_keywords"`
	DeveloperUrl    string `json:"developer_url"`
}
