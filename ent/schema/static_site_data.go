package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type StaticSiteData struct {
	ent.Schema
}

func (StaticSiteData) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id"),
		field.JSON("about_us", []byte{}),
		field.JSON("how_we_work", []byte{}),
		field.JSON("testimonials", []byte{}),
		field.JSON("mango_insights", []byte{}),
		field.JSON("our_associations", []byte{}),
		field.Time("updated_at"),
		field.JSON("categories_with_amenities", struct {
			Categories map[string][]struct {
				Icon  string `json:"icon"`
				Value string `json:"value"`
			} `json:"categories"`
		}{}),
		field.JSON("property_types", PropertyTypes{}),
		field.Time("created_at").Default(time.Now),
	}
}

type PropertyTypes struct {
	Commercial  []string `json:"commercial"`
	Residential []string `json:"residential"`
}
