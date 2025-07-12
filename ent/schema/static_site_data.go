package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
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
		field.JSON("categories_with_amenities", struct {
			Categories map[string][]struct {
				Icon  string `json:"icon"`
				Value string `json:"value"`
			} `json:"categories"`
		}{}).Optional(),
		field.JSON("property_types", PropertyTypes{}),
		field.JSON("generic_search_data", []GenericSearchData{}),
		field.Bool("is_active").Default(true),
		field.Time("deleted_at").Optional().Nillable(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Indexes of the StaticSiteData.
func (StaticSiteData) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("is_active"),
		index.Fields("categories_with_amenities"),
	}
}

type PropertyTypes struct {
	Commercial  []string `json:"commercial"`
	Residential []string `json:"residential"`
}

type GenericSearchData struct {
	CanonicalURL string            `json:"canonical_url"`
	SearchTerm   string            `json:"title"`
	Filters      map[string]string `json:"filters"`
}
