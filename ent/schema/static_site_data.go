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
		field.String("id").Unique(),
		field.JSON("about_us", []byte{}).Optional(),
		field.JSON("how_we_work", []byte{}).Optional(),
		field.JSON("testimonials", []byte{}).Optional(),
		field.JSON("mango_insights", []byte{}).Optional(),
		field.JSON("our_associations", []byte{}).Optional(),
		field.Time("updated_at").Optional(),
		field.JSON("categories_with_amenities", struct {
			Categories map[string][]struct {
				Icon  string `json:"icon"`
				Value string `json:"value"`
			} `json:"categories"`
		}{}).Optional(),
		field.JSON("property_types", PropertyTypes{}),
		field.Time("created_at").Default(time.Now),
	}
}

type PropertyTypes struct {
	Commercial  []string `json:"commercial"`
	Residential []string `json:"residential"`
}
