package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type CustomSearchPage struct {
	ent.Schema
}

func (CustomSearchPage) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique(),
		field.String("title"),
		field.String("description"),
		field.String("slug"),
		field.JSON("filters", map[string]interface{}{}).Optional(),
		field.String("search_term"),
		field.JSON("meta_info", MetaInfo{}).Optional(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (CustomSearchPage) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("slug", "filters"),
	}
}

type MetaInfo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Keywords    string `json:"keywords"`
}
