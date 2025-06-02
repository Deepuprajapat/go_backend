package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Blogs struct {
	ent.Schema
}

func (Blogs) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique(),
		field.JSON("seo_meta_info", SEOMetaInfo{}),
		field.String("blog_url"),
		field.JSON("blog_content", BlogContent{}),
		field.Bool("is_priority").Default(false),
		field.Time("created_at").Default(time.Now()),
		field.Time("updated_at"),
	}
}

// Edges of the Blogs.
func (Blogs) Edges() []ent.Edge {
	return []ent.Edge{
		// Many blogs can be updated by one user
		edge.From("updated_by", User.Type).
			Ref("updated_blogs").
			Unique(), // Each blog has only one user who last updated it
	}
}

type BlogContent struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	Image        string `json:"image"`
	ImageAlt     string `json:"image_alt"`
	ImageCaption string `json:"image_caption"`
}

type SEOMetaInfo struct {
	BlogSchema  string `json:"blog_schema"`
	Canonical   string `json:"canonical"`
	Title       string `json:"title"`
	Keywords    string `json:"keywords"`
	Description string `json:"description"`
}
