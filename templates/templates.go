package templates

import (
	"embed"
	"html/template"
)

//go:embed seo/*.html
var FS embed.FS

// ParseSEOTemplates parses all SEO templates from embedded filesystem
func ParseSEOTemplates() (*template.Template, error) {
	return template.ParseFS(FS, "seo/*.html")
}