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
		field.String("id").Unique(),
		field.String("name"),
		field.String("legal_name").Optional(),
		field.String("identifier").Optional(),
		field.Int("established_year"),
		field.JSON("media_content", DeveloperMediaContent{}).Optional(),
		field.Bool("is_verified").Default(false),
		field.Bool("is_active").Default(true),
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
