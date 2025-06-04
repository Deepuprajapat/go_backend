package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Developer struct {
	Base
}

func (Developer) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique(),
		field.String("name"),
		field.String("legal_name"),
		field.String("identifier"),
		field.Int("established_year"),
		field.JSON("media_content", DeveloperMediaContent{}),
		field.Bool("is_verified").Default(false),
	}
}

func (Developer) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("projects", Project.Type),
	}
}

// media and content
type DeveloperMediaContent struct {
	DeveloperAddress string `json:"developer_address"`
	Phone            string `json:"phone"`
	DeveloperLogo    string `json:"developer_logo"`
	AltDeveloperLogo string `json:"alt_developer_logo"`
	About            string `json:"about"`
	Overview         string `json:"overview"`
	Disclaimer       string `json:"disclaimer"`
}
