package response

import (
	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/ent/schema"
	"github.com/VI-IM/im_backend_go/internal/domain/enums"
)

type ProjectComparisonResponse struct {
	Projects []*ProjectComparison `json:"projects"`
}

type ProjectComparison struct {
	ProjectID     string                 `json:"project_id"`
	ProjectName   string                 `json:"project_name"`
	Description   string                 `json:"description"`
	Status        enums.ProjectStatus    `json:"status"`
	MinPrice      string                 `json:"min_price"`
	MaxPrice      string                 `json:"max_price"`
	PriceUnit     string                 `json:"price_unit"`
	TimelineInfo  schema.TimelineInfo    `json:"timeline_info"`
	LocationInfo  schema.LocationInfo    `json:"location_info"`
	IsFeatured    bool                   `json:"is_featured"`
	IsPremium     bool                   `json:"is_premium"`
	IsPriority    bool                   `json:"is_priority"`
	WebCards      schema.ProjectWebCards `json:"web_cards"`
	DeveloperName string                 `json:"developer_name,omitempty"`
}

type Project struct {
	ProjectID     string                 `json:"project_id"`
	ProjectName   string                 `json:"project_name"`
	Description   string                 `json:"description"`
	Slug          string                 `json:"slug"`
	Status        enums.ProjectStatus    `json:"status"`
	MinPrice      string                 `json:"min_price"`
	MaxPrice      string                 `json:"max_price"`
	PriceUnit     string                 `json:"price_unit"`
	TimelineInfo  schema.TimelineInfo    `json:"timeline_info"`
	MetaInfo      schema.SEOMeta         `json:"meta_info"`
	WebCards      schema.ProjectWebCards `json:"web_cards"`
	LocationInfo  schema.LocationInfo    `json:"location_info"`
	City          string                 `json:"city"`
	DeveloperInfo DeveloperInfo          `json:"developer_info"`
	IsFeatured    bool                   `json:"is_featured"`
	IsPremium     bool                   `json:"is_premium"`
	IsPriority    bool                   `json:"is_priority"`
}

type DeveloperInfo struct {
	DeveloperID     string `json:"developer_id"`
	DeveloperName   string `json:"developer_name"`
	Phone           string `json:"phone"`
	Logo            string `json:"logo"`
	AltLogo         string `json:"alt_logo"`
	Address         string `json:"address"`
	EstablishedYear int    `json:"established_year"`
	TotalProjects   int    `json:"total_projects"`
}

type AddProjectResponse struct {
	ProjectID string `json:"project_id"`
}

type UpdateProjectResponse struct {
	ProjectID string `json:"project_id"`
}

type ProjectListResponse struct {
	ProjectID     string   `json:"project_id"`
	ProjectName   string   `json:"project_name"`
	ShortAddress  string   `json:"short_address"`
	City          string   `json:"city"`
	Slug     string   `json:"slug"`
	Images        []string `json:"images"`
	Configuration string   `json:"configuration"`
	MinPrice      string   `json:"min_price"`
	Sizes         string   `json:"sizes"`
	IsPremium     bool     `json:"is_premium"`
	VideoURLs     []string `json:"video_urls"`
	FullDetails   *Project `json:"full_details,omitempty"`
}

func GetProjectFromEnt(project *ent.Project) *Project {
	return &Project{
		ProjectID:   project.ID,
		ProjectName: project.Name,
		Description: project.Description,
		Status:      project.Status,
		Slug:        project.Slug,
		MinPrice:    project.MinPrice,
		MaxPrice:    project.MaxPrice,
		TimelineInfo: schema.TimelineInfo{
			ProjectLaunchDate:     project.TimelineInfo.ProjectLaunchDate,
			ProjectPossessionDate: project.TimelineInfo.ProjectPossessionDate,
		},
		MetaInfo: project.MetaInfo,
		WebCards: project.WebCards,
		LocationInfo: schema.LocationInfo{
			ShortAddress:  project.LocationInfo.ShortAddress,
			Longitude:     project.LocationInfo.Longitude,
			Latitude:      project.LocationInfo.Latitude,
			GoogleMapLink: project.LocationInfo.GoogleMapLink, // project.Edges.Location.PhoneNumber
		},
		City: func() string {
			if project.Edges.Location != nil {
				return project.Edges.Location.City
			}
			return ""
		}(),
		DeveloperInfo: DeveloperInfo{
			DeveloperID:   project.Edges.Developer.ID,
			DeveloperName: project.Edges.Developer.Name,
			Phone: func() string {
				if project.Edges.Location != nil {
					return project.Edges.Location.PhoneNumber
				}
				return ""
			}(),
			Logo:            project.Edges.Developer.MediaContent.DeveloperLogo,
			AltLogo:         project.Edges.Developer.MediaContent.AltDeveloperLogo,
			Address:         project.Edges.Developer.MediaContent.DeveloperAddress,
			EstablishedYear: project.Edges.Developer.EstablishedYear,
		},
		IsFeatured: project.IsFeatured,
		IsPremium:  project.IsPremium,
		IsPriority: project.IsPriority,
	}
}

func GetProjectListResponse(project *ent.Project) *ProjectListResponse {
	return &ProjectListResponse{
		ProjectID:   project.ID,
		ProjectName: project.Name,
		City: func() string {
			if project.Edges.Location != nil {
				return project.Edges.Location.City
			}
			return ""
		}(),
		ShortAddress:  project.LocationInfo.ShortAddress,
		IsPremium:     project.IsPremium,
		Images:        project.WebCards.Images,
		Configuration: project.WebCards.Details.Configuration.Value,
		Sizes:         project.WebCards.Details.Sizes.Value,
		VideoURLs:     project.WebCards.VideoPresentation.URLs,
		MinPrice:      project.MinPrice,
		Slug:          project.Slug,
	}
}
