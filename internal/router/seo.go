package router

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/VI-IM/im_backend_go/internal/application"
	"github.com/VI-IM/im_backend_go/internal/config"
	"github.com/VI-IM/im_backend_go/response"
	"github.com/VI-IM/im_backend_go/templates"
)

var (
	// seoTemplates holds parsed HTML templates
	seoTemplates *template.Template
)

// Template data structures for SEO pages
type ProjectSEOData struct {
	Project      *response.Project
	MetaInfo     ProjectMetaInfo
	ProjectImage string
	BaseURL      string
}

type PropertySEOData struct {
	Property      *response.Property
	MetaInfo      PropertyMetaInfo
	PropertyImage string
	BaseURL       string
}

type BlogSEOData struct {
	BlogResponse *response.BlogResponse
	SEOMetaInfo  BlogSEOMetaInfo
	BlogContent  BlogContentInfo
	CanonicalURL string
	CreatedAt    string
	UpdatedAt    string
	BaseURL      string
}

// Meta info structures for templates
type ProjectMetaInfo struct {
	Title         string
	Description   string
	Keywords      string
	ProjectSchema []string
}

type PropertyMetaInfo struct {
	Title       string
	Description string
	Keywords    string
	Canonical   string
}

type BlogSEOMetaInfo struct {
	Title       string
	Description string
	Keywords    string
	BlogSchema  []string
}

type BlogContentInfo struct {
	Title        string
	Description  string
	Image        string
	ImageAlt     string
	ImageCaption string
	Content      string
}

// URL helper functions
func getBaseURL() string {
	cfg := config.GetConfig()
	return cfg.Server.BaseURL
}

func getProjectURL(slug string) string {
	return fmt.Sprintf("%s/%s", getBaseURL(), slug)
}

func getBlogURL(slug string) string {
	return fmt.Sprintf("%s/blogs/%s", getBaseURL(), slug)
}

func getPropertyURL(slug string) string {
	return fmt.Sprintf("%s/propertyforsale/%s", getBaseURL(), slug)
}

func getDefaultImage(imageType string) string {
	return fmt.Sprintf("%s/default-%s-image.jpg", getBaseURL(), imageType)
}

// initSEOTemplates initializes the SEO template system
func initSEOTemplates() error {
	var err error
	seoTemplates, err = templates.ParseSEOTemplates()
	if err != nil {
		return fmt.Errorf("failed to parse SEO templates: %v", err)
	}
	return nil
}

// isBotUserAgent checks if the request is from a bot/crawler
func isBotUserAgent(userAgent string) bool {
	userAgentLower := strings.ToLower(userAgent)
	botSignatures := []string{
		"bot", "crawler", "spider", "googlebot", "bingbot",
		"slurp", "duckduckbot", "baiduspider", "yandexbot",
		"facebookexternalhit", "twitterbot", "linkedinbot",
		"whatsapp", "slack",
	}

	for _, signature := range botSignatures {
		if strings.Contains(userAgentLower, signature) {
			return true
		}
	}
	return false
}

// handleSEORoute handles SEO-specific routes for bots
func handleSEORoute(w http.ResponseWriter, r *http.Request, app application.ApplicationInterface) bool {
	path := r.URL.Path

	// Handle project routes: /<project_slug>
	if projectSlugMatch := strings.TrimPrefix(path, "/"); projectSlugMatch != "" && !strings.Contains(projectSlugMatch, "/") {
		// Check if this is a valid project slug
		if app != nil {
			project, err := app.GetProjectBySlug(projectSlugMatch)
			if err == nil && project != nil {
				generateProjectSEOHTML(w, project)
				return true
			}
		}
	}

	//// Handle blog routes: /blogs/<blog_slug>
	if strings.HasPrefix(path, "/blogs/") {
		blogSlug := strings.TrimPrefix(path, "/blogs/")
		if blogSlug != "" && app != nil {
			blog, err := app.GetBlogBySlug(blogSlug)
			if err == nil && blog != nil {
				generateBlogSEOHTML(w, blog)
				return true
			}
		}
	}

	// Handle property routes: /propertyforsale/<property_slug>
	if strings.HasPrefix(path, "/propertyforsale/") {
		propertySlug := strings.TrimPrefix(path, "/propertyforsale/")
		if propertySlug != "" && app != nil {
			property, err := app.GetPropertyBySlug(r.Context(), propertySlug)
			if err == nil && property != nil {
				generatePropertySEOHTML(w, property)
				return true
			}
		}
	}

	return false
}

// generateProjectSEOHTML generates SEO-friendly HTML for project pages using templates
func generateProjectSEOHTML(w http.ResponseWriter, project *response.Project) {
	// Create template data with proper meta info including project_schema
	templateData := ProjectSEOData{
		Project: project,
		MetaInfo: ProjectMetaInfo{
			Title:         project.MetaInfo.Title,
			Description:   project.MetaInfo.Description,
			Keywords:      project.MetaInfo.Keywords,
			ProjectSchema: project.MetaInfo.ProjectSchema, // Include project_schema from database
		},
		ProjectImage: func() string {
			if len(project.WebCards.Images) > 0 {
				return project.WebCards.Images[0]
			}
			return getDefaultImage("project")
		}(),
		BaseURL: getBaseURL(),
	}

	// Provide defaults if SEO meta is empty
	if templateData.MetaInfo.Title == "" {
		templateData.MetaInfo.Title = project.ProjectName
	}
	if templateData.MetaInfo.Description == "" {
		templateData.MetaInfo.Description = fmt.Sprintf("%s - %s", project.ProjectName, project.LocationInfo.ShortAddress)
	}
	if templateData.MetaInfo.Keywords == "" {
		templateData.MetaInfo.Keywords = fmt.Sprintf("%s, %s, real estate, property", project.ProjectName, project.City)
	}

	// Set content type and render template
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	err := seoTemplates.ExecuteTemplate(w, "project_seo.html", templateData)
	if err != nil {
		fmt.Fprintf(w, "Error rendering template: %v", err)
	}
}

// generateBlogSEOHTML generates SEO-friendly HTML for blog pages using templates
func generateBlogSEOHTML(w http.ResponseWriter, blog *response.BlogResponse) {
	// Create canonical URL
	canonicalURL := blog.SEOMetaInfo.Canonical
	if canonicalURL == "" {
		canonicalURL = getBlogURL(blog.Slug)
	}

	// Create template data with proper meta info including blog_schema
	templateData := BlogSEOData{
		BlogResponse: blog,
		SEOMetaInfo: BlogSEOMetaInfo{
			Title:       blog.SEOMetaInfo.Title,
			Description: blog.SEOMetaInfo.Description,
			Keywords:    blog.SEOMetaInfo.Keywords,
			BlogSchema:  blog.SEOMetaInfo.BlogSchema, // Include blog_schema from database
		},
		BlogContent: BlogContentInfo{
			Title:        blog.BlogContent.Title,
			Description:  blog.BlogContent.Description,
			Image:        blog.BlogContent.Image,
			ImageAlt:     blog.BlogContent.ImageAlt,
			ImageCaption: blog.BlogContent.ImageCaption,
			Content:      strings.ReplaceAll(blog.BlogContent.Description, "\n", "<br>"),
		},
		CanonicalURL: canonicalURL,
		CreatedAt:    time.Unix(blog.CreatedAt, 0).Format("2006-01-02T15:04:05Z"),
		UpdatedAt:    time.Unix(blog.UpdatedAt, 0).Format("2006-01-02T15:04:05Z"),
		BaseURL:      getBaseURL(),
	}

	// Set content type and render template
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	err := seoTemplates.ExecuteTemplate(w, "blog_seo.html", templateData)
	if err != nil {
		fmt.Fprintf(w, "Error rendering template: %v", err)
	}
}

// generatePropertySEOHTML generates SEO-friendly HTML for property pages using templates
func generatePropertySEOHTML(w http.ResponseWriter, property *response.Property) {
	// Create template data with proper meta info
	templateData := PropertySEOData{
		Property: property,
		MetaInfo: PropertyMetaInfo{
			Title:       property.MetaInfo.Title,
			Description: property.MetaInfo.Description,
			Keywords:    property.MetaInfo.Keywords,
			Canonical:   property.MetaInfo.Canonical,
		},
		PropertyImage: func() string {
			if len(property.PropertyImages) > 0 {
				return property.PropertyImages[0]
			}
			return getDefaultImage("property")
		}(),
		BaseURL: getBaseURL(),
	}

	// Provide defaults if SEO meta is empty
	if templateData.MetaInfo.Title == "" {
		templateData.MetaInfo.Title = fmt.Sprintf("%s - Property for Sale", property.Name)
	}
	if templateData.MetaInfo.Description == "" {
		templateData.MetaInfo.Description = fmt.Sprintf("Property for sale - %s", property.Name)
	}
	if templateData.MetaInfo.Keywords == "" {
		templateData.MetaInfo.Keywords = fmt.Sprintf("%s, property for sale, real estate", property.Name)
	}

	// Set content type and render template
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	err := seoTemplates.ExecuteTemplate(w, "property_seo.html", templateData)
	if err != nil {
		fmt.Fprintf(w, "Error rendering template: %v", err)
	}
}
